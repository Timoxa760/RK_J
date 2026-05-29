package categorizer

import "backend_project/services/money-intelligence/ai-processor/internal"

type Categorizer interface {
	Categorize(items []internal.Item) []internal.CategorizedItem
}
