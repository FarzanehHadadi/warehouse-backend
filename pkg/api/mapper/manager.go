// mapper/manager_mapper.go
package mapper

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
)

// ToManagerDetailResponse converts Manager model to ManagerDetailResponse DTO
func ToManagerDetailResponse(manager *models.Manager) dto.ManagerSummary {
	if manager == nil {
		return dto.ManagerSummary{}
	}

	return dto.ManagerSummary{
		ID:          manager.ID,
		Name:        manager.Name,
		Phone:       manager.Phone,
		Departments: ToDepartmentSummaries(manager.Departments),
	}
}

// ToManagerSummaries converts slice of Manager models to ManagerSummary DTOs
func ToManagerSummaries(managers []*models.Manager) []dto.ManagerSummary {
	if managers == nil {
		return []dto.ManagerSummary{}
	}

	summaries := make([]dto.ManagerSummary, len(managers))
	for i, manager := range managers {
		summaries[i] = ToManagerDetailResponse(manager)
	}
	return summaries
}

// ToDepartmentSummaries converts slice of Department models to DepartmentSummary DTOs
func ToDepartmentSummaries(departments []models.Department) []dto.DepartmentSummary {
	if departments == nil {
		return []dto.DepartmentSummary{}
	}

	summaries := make([]dto.DepartmentSummary, len(departments))
	for i, dept := range departments {
		summaries[i] = dto.DepartmentSummary{
			ID:   dept.ID,
			Name: dept.Name,
		}
	}
	return summaries
}
