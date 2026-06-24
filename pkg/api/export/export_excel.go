package export

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"warehouse/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// getDisplayValue converts any field into a clean, user-friendly value
func getDisplayValue(v reflect.Value) interface{} {
	if !v.IsValid() {
		return ""
	}

	// Dereference pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return ""
		}

		var items []string
		for i := 0; i < v.Len(); i++ {
			item := getDisplayValue(v.Index(i)) // recursive call
			if str, ok := item.(string); ok && str != "" {
				items = append(items, str)
			} else if item != nil {
				items = append(items, fmt.Sprintf("%v", item))
			}
		}

		if len(items) == 0 {
			return ""
		}
		return strings.Join(items, ", ") // You can change to "\n" for new lines

	case reflect.Struct:
		typeName := v.Type().Name()
		fullType := v.Type().String()

		// Handle models.Date
		if typeName == "Date" || fullType == "models.Date" {
			if tm := v.FieldByName("Time"); tm.IsValid() {
				if t, ok := tm.Interface().(time.Time); ok {
					return t
				}
			}
		}

		// Handle models.Basic
		if typeName == "Basic" || strings.Contains(fullType, "models.Basic") {
			for _, fieldName := range []string{"ID", "UpdatedAt", "CreatedAt"} {
				if f := v.FieldByName(fieldName); f.IsValid() {
					return getDisplayValue(f)
				}
			}
		}

		// Handle other nested structs (SimpleSummary, Department, etc.)
		if nameField := v.FieldByName("Name"); nameField.IsValid() {
			return getDisplayValue(nameField)
		}
		if idField := v.FieldByName("ID"); idField.IsValid() {
			return getDisplayValue(idField)
		}

		// Fallback for unknown struct
		return v.Interface()

	case reflect.String:
		return v.String()
	}

	// Default
	return v.Interface()

}

// ExportToExcel - Generic & User-Friendly
func ExportToExcel[T any](c *gin.Context, data []T, sheetName, filePrefix string, customHeaders map[string]string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to export")
	}

	f := excelize.NewFile()
	defer f.Close()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// Get struct type (handle slice of pointer)
	sample := data[0]
	t := reflect.TypeOf(sample)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("data must be slice of struct or *struct")
	}

	// === Headers ===
	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		header := field.Tag.Get("json")
		if header == "" || header == "-" {
			header = field.Name
		} else if idx := strings.Index(header, ","); idx != -1 {
			header = header[:idx]
		}
		if customHeaders != nil {
			if customName, exists := customHeaders[field.Name]; exists {
				header = customName
			}
		}
		headers[i] = header
	}

	for col, header := range headers {
		f.SetCellValue(sheetName, fmt.Sprintf("%s1", string(rune('A'+col))), header)
	}

	// === Styles ===
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	dateStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 14}) // yyyy-mm-dd
	numberStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 1, Alignment: &excelize.Alignment{Horizontal: "left"}})
	stringStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 49, Alignment: &excelize.Alignment{Horizontal: "center"}})

	f.SetRowHeight(sheetName, 1, 25)
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", string(rune('A'+len(headers)-1))), headerStyle)

	// === Data Rows ===
	row := 2
	for _, item := range data {
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				row++
				continue
			}
			v = v.Elem()
		}
		logger.Log.Info("v", zap.Any("v", v))
		for col := 0; col < t.NumField(); col++ {
			fieldVal := v.Field(col)
			cell := fmt.Sprintf("%s%d", toExcelColumn(col), row)
			displayValue := getDisplayValue(fieldVal)
			switch val := displayValue.(type) {
			case time.Time:
				if !val.IsZero() {
					f.SetCellValue(sheetName, cell, val)
					f.SetCellStyle(sheetName, cell, cell, dateStyle)
				}
			default:
				f.SetCellValue(sheetName, cell, displayValue)
				if k := fieldVal.Kind(); k == reflect.String || k == reflect.Bool || k == reflect.Ptr || k == reflect.Struct {
					f.SetCellStyle(sheetName, cell, cell, stringStyle)
				}
				// Auto number style
				if k := fieldVal.Kind(); k == reflect.Int || k == reflect.Int64 || k == reflect.Float64 || k == reflect.Uint || k == reflect.Uint64 {
					f.SetCellStyle(sheetName, cell, cell, numberStyle)
				}
			}
		}
		row++
	}

	// Auto column width
	for col := 'A'; col < rune('A'+len(headers)); col++ {
		f.SetColWidth(sheetName, string(col), string(col), 22)
	}

	// Send file
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_%s.xlsx", filePrefix, time.Now().Format("2006-01-02_15-04-05"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())

	return nil
}

// toExcelColumn converts 0 -> "A", 1 -> "B", ..., 25 -> "Z", 26 -> "AA", etc.
func toExcelColumn(col int) string {
	var s string
	for col >= 0 {
		s = string(rune('A'+col%26)) + s
		col = col/26 - 1
	}
	return s
}
