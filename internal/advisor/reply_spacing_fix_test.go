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

func TestRepairSplitRussianWords_EmergencyFundSample(t *testing.T) {
	raw := `Сейчас вся подушка — 350 000 ₽ — лежит без укания, где именно (в профиле все три поля: наличные, де по зит, инвестиции — равны нулю).
Депозит дают около 18–21 % го до вых. Для подушки накопительный счёт пред по чтительнее — главная задача подушки это доступность в кр из ис.
Откройте накопительный счёт и укажите сумму в поле «де по зит».`
	got := RepairSplitRussianWords(raw)
	for _, broken := range []string{
		"де по зит", "го до вых", "пред по чтительнее", "кр из ис", "укания",
	} {
		if strings.Contains(got, broken) {
			t.Fatalf("still broken %q in %q", broken, got)
		}
	}
	for _, fixed := range []string{"депозит", "годовых", "предпочтительнее", "кризис", "уточнения"} {
		if !strings.Contains(got, fixed) {
			t.Fatalf("expected %q in %q", fixed, got)
		}
	}
}

func TestRepairSplitRussianWords_MonthGlued(t *testing.T) {
	got := RepairSplitRussianWords("в месяцона покрывает")
	if strings.Contains(got, "месяцона") {
		t.Fatalf("got %q", got)
	}
	if !strings.Contains(got, "месяц она") {
		t.Fatalf("got %q", got)
	}
}

func TestRepairSplitRussianWords_UvasUze(t *testing.T) {
	got := RepairSplitRussianWords("увасуже есть")
	if strings.Contains(got, "увасуже") {
		t.Fatalf("got %q", got)
	}
}

func TestRepairSplitRussianWords_ZamenyNet(t *testing.T) {
	re := reExplicitSplitFixes[18]
	rep := explicitSplitReplacements[18]
	if !re.MatchString("заменынet") {
		t.Fatalf("no match for %q in %q", "заменынet", re.String())
	}
	direct := re.ReplaceAllString("заменынet.", rep)
	if !strings.Contains(direct, "замены нет") {
		t.Fatalf("direct %q rep=%q", direct, rep)
	}
	got := RepairSplitRussianWords("заменынet.")
	if !strings.Contains(got, "замены нет") {
		t.Fatalf("got %q", got)
	}
}

func TestRepairSplitRussianWords_PassiveIncomeSample(t *testing.T) {
	raw := `Идея понятна — ис по льзовать часть подушки. единственныйфи на нсовый буфер. ` +
		`При доходе 15 000 ₽ в месяцона покрывает. заменынet. увасуже есть 300 000 ₽. не до стающие 50 000 ₽.`
	got := RepairSplitRussianWords(raw)
	for _, broken := range []string{
		"ис по льз", "единственныйфи", "на нсовый", "месяцона", "заменынet", "увасуже", "не до стающ",
	} {
		if strings.Contains(got, broken) {
			t.Fatalf("still broken %q in %q", broken, got)
		}
	}
	for _, fixed := range []string{
		"использовать", "единственный финансовый", "месяц она", "замены нет", "у вас уже", "недостающие",
	} {
		if !strings.Contains(got, fixed) {
			t.Fatalf("expected %q in %q", fixed, got)
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
