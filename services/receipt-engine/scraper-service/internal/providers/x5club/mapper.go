package x5club

import (
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) ToRawReceipts(items []HistoryItem) []scrap.RawReceipt {
	result := make([]scrap.RawReceipt, 0, len(items))
	for _, item := range items {
		parsedDate, _ := time.Parse("2006-01-02", item.Date)
		r := scrap.RawReceipt{
			ID:       item.ID,
			Provider: "x5club",
			Store:    item.StoreName,
			Date:     parsedDate,
			Total:    item.Total,
		}
		for _, it := range item.Items {
			r.Items = append(r.Items, scrap.RawItem{
				Name:     it.Name,
				Price:    it.Price,
				Quantity: it.Quantity,
			})
		}
		result = append(result, r)
	}
	return result
}
