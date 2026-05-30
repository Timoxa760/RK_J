package advisor

import (
	"fmt"
	"strings"
)

func BuildChatActions(snap Snapshot, message, reply string) []ChatAction {
	q := strings.ToLower(strings.TrimSpace(message))
	var actions []ChatAction

	if snap.DataCompleteness["income"] == "skipped" || snap.DataCompleteness["income"] == "unknown" {
		actions = append(actions, ChatAction{
			Type:         "open_profile",
			Label:        "Заполнить доход",
			Path:         "/profile",
			ProfileField: "income",
		})
	}

	if (strings.Contains(q, "кредит") || strings.Contains(q, "ставк")) && len(snap.Credits.Credits) == 0 {
		actions = append(actions, ChatAction{
			Type:  "navigate",
			Label: "Загрузить договор",
			Path:  "/credits",
		})
	}

	if strings.Contains(q, "добав") && (strings.Contains(q, "трат") || strings.Contains(q, "расход") || strings.Contains(q, "покуп")) {
		actions = append(actions, ChatAction{
			Type:  "open_add_expense",
			Label: "Добавить покупку",
		})
	}

	if strings.Contains(q, "урез") || strings.Contains(q, "сократ") || strings.Contains(q, "эконом") {
		if top := TopCategory(snap.Spending); top != nil {
			actions = append(actions, ChatAction{
				Type:  "ask_followup",
				Label: fmt.Sprintf("Сократить «%s»", top.Name),
				Ask:   fmt.Sprintf("Как сократить траты в категории «%s» на 10–20%%?", top.Name),
			})
		}
	}

	if strings.Contains(q, "план") || strings.Contains(q, "шаг") {
		actions = append(actions, ChatAction{
			Type:  "navigate",
			Label: "Открыть план",
			Path:  "/dashboard",
			Hash:  "#plan",
		})
	}

	if len(actions) == 0 {
		plan := BuildPlanResponse(snap, nil)
		if plan.Diagnosis.MainAction.Title != "" {
			actions = append(actions, ChatAction{
				Type:  "navigate",
				Label: "К совету недели",
				Path:  "/dashboard",
				Hash:  "#plan",
			})
		}
	}

	if len(actions) > 3 {
		actions = actions[:3]
	}
	return actions
}
