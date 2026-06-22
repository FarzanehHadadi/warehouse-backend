package mapper

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
)

func ToStoreDetailResponse(store *models.Store) dto.StoreSummary {
	if store == nil {

		return dto.StoreSummary{}
	}
	return dto.StoreSummary{
		ID:          store.ID,
		Name:        store.Name,
		ManagerName: store.Manager.Name,
		IsActive:    store.IsActive,
	}
}

func ToStoreSummaries(stores []*models.Store) []dto.StoreSummary {
	if stores == nil {

		return []dto.StoreSummary{}
	}

	str := make([]dto.StoreSummary, len(stores))
	for i, store := range stores {
		str[i] = ToStoreDetailResponse(store)
	}
	return str
}
