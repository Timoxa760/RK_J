package advisor

import (
	"strings"
	"testing"
)

func TestRepairSplitRussianWords_UserSample(t *testing.T) {
	raw := `Из трат за месяц — только «Продукты» на 3 500 ₽. доход ы вы пропустили при за по лнении, по это му точный свободный остаток я посчитать не могу.
Без полной картины расход ов и доход а я не знаю, сколько у вас реально остаётся каждый месяц.
за сколько месяц ев реально накопить 90 000 ₽ на отпуск в Сочи.`
	got := RepairSplitRussianWords(raw)
	for _, broken := range []string{
		"доход ы", "за по лнении", "по это му", "расход ов", "доход а", "месяц ев",
	} {
		if strings.Contains(got, broken) {
			t.Fatalf("still broken %q in %q", broken, got)
		}
	}
	if !strings.Contains(got, "доходы") || !strings.Contains(got, "заполнении") || !strings.Contains(got, "поэтому") {
		t.Fatalf("expected fixed words in %q", got)
	}
}

func TestRepairSplitRussianWords_PreservesPhrases(t *testing.T) {
	phrases := []string{
		"идти по улице",
		"не могу",
		"что делать",
		"на 3 500 ₽",
	}
	for _, phrase := range phrases {
		got := RepairSplitRussianWords(phrase)
		if got != phrase {
			t.Fatalf("changed valid phrase %q -> %q", phrase, got)
		}
	}
}

func TestRepairSplitRussianWords_HistoricalGluedFixes(t *testing.T) {
	raw := "ука за ть обя зательные транс по рт кредит ы"
	got := RepairSplitRussianWords(raw)
	for _, broken := range []string{"ука за ть", "транс по рт", "кредит ы"} {
		if strings.Contains(got, broken) {
			t.Fatalf("still broken %q in %q", broken, got)
		}
	}
}
