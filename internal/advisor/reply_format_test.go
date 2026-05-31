package advisor

import (
	"strings"
	"testing"
)

func TestFormatAdvisorReply_PreservesNormalWords(t *testing.T) {
	raw := "расходы обязательные указать транспорт кредиты понятно работают без расходов"
	got := FormatAdvisorReply(raw)
	for _, broken := range []string{"расход ы", "обя зательные", "ука за ть", "транс по рт", "кредит ы", "по нятно", "ра бот"} {
		if strings.Contains(got, broken) {
			t.Fatalf("broke normal word %q in %q", broken, got)
		}
	}
}

func TestFormatAdvisorReply_UserSample(t *testing.T) {
	raw := "Цель: Крупная покупка —120 000 ₽Проблема: притекущихданных свободныхденег нет.|||||---||Доход | 15 000₽ |Кредит (фикс.) | −15 000 ₽| Остаток | 0 ₽ |Откладыватьнацель пока нечего— доход полностью уходит накредит.---Главныйвопрос: кредит на15 000₽/мес —это единственный расход, илиесть другие траты"
	got := FormatAdvisorReply(raw)
	if strings.Contains(got, "||||") {
		t.Fatalf("pipes remain: %q", got)
	}
	if strings.Contains(got, "₽Проблема") {
		t.Fatalf("missing space after ruble: %q", got)
	}
	if !strings.Contains(got, "- Доход") && !strings.Contains(got, "- ") {
		t.Fatalf("expected list items: %q", got)
	}
	t.Log(got)
}

func TestFormatAdvisorReply_PipeTable(t *testing.T) {
	raw := "Цель: Крупная покупка —120 000 ₽Проблема: при текущих данных|||||---||Доход | 15 000₽ |Кредит | −15 000 ₽|Откладывать на цель пока нечего"
	got := FormatAdvisorReply(raw)
	if strings.Contains(got, "||||") {
		t.Fatalf("pipes remain: %q", got)
	}
	if !strings.Contains(got, "- ") {
		t.Fatalf("expected list from table: %q", got)
	}
}
