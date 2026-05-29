package magnit

import (
	"time"

	scrap "backend_project/services/receipt-engine/scraper-service/internal"
)

type Mapper struct {
	innDetector *INNDetector
}

func NewMapper(innDetector *INNDetector) *Mapper {
	return &Mapper{innDetector: innDetector}
}

func (m *Mapper) ToRawReceipts(receipts []Receipt) []scrap.RawReceipt {
	result := make([]scrap.RawReceipt, 0, len(receipts))
	for _, r := range receipts {
		parsedDate, _ := time.Parse("2006-01-02", r.Date)
		raw := scrap.RawReceipt{
			ID:       r.ID,
			Provider: "magnit",
			Store:    r.StoreName,
			Date:     parsedDate,
			Total:    r.Total,
		}
		for _, it := range r.Items {
			raw.Items = append(raw.Items, scrap.RawItem{
				Name:     it.Name,
				Price:    it.Price,
				Quantity: it.Quantity,
			})
		}
		result = append(result, raw)
	}
	return result
}

func (m *Mapper) ToRawReceiptsWithINN(receipts []Receipt) []scrap.RawReceipt {
	result := m.ToRawReceipts(receipts)
	for i := range result {
		result[i].ID = m.innDetector.Detect(result[i].Store) + "_" + result[i].ID
	}
	return result
}
