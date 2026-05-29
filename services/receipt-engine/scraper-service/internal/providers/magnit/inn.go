package magnit

import "strings"

var storeINNMap = map[string]string{
	"магнит":            "2309085638",
	"магнит косметик":   "2309085639",
	"магнит семейный":   "2309085640",
	"магнит оптик":      "2309085641",
	"магнит аптека":     "2309085642",
	"дикси":             "7702610669",
	"пятёрочка":         "7727563778",
	"перекрёсток":       "7729117676",
	"карусель":          "7743650117",
	"лента":             "7814142098",
}

type INNDetector struct {
	stores map[string]string
}

func NewINNDetector() *INNDetector {
	m := make(map[string]string, len(storeINNMap))
	for k, v := range storeINNMap {
		m[k] = v
	}
	return &INNDetector{stores: m}
}

func (d *INNDetector) Detect(storeName string) string {
	lower := strings.ToLower(strings.TrimSpace(storeName))
	if inn, ok := d.stores[lower]; ok {
		return inn
	}
	for prefix, inn := range d.stores {
		if strings.Contains(lower, prefix) {
			return inn
		}
	}
	return "2309085699"
}
