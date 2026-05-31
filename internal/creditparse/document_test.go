package creditparse

import "testing"

func TestIsGeneralConditionsTemplate_AlfaSample(t *testing.T) {
	text := `Общие условия Договора потребительского кредита,
предусматривающего выдачу Кредита наличными
ОБЩИЕ УСЛОВИЯ ДОГОВОРА ПОТРЕБИТЕЛЬСКОГО КРЕДИТА,
ПРЕДУСМАТРИВАЮЩЕГО ВЫДАЧУ КРЕДИТА НАЛИЧНЫМИ
Заемщик уплачивает Банку проценты по ставке, указанной в Индивидуальных условиях выдачи Кредита наличными.`
	if !IsGeneralConditionsTemplate(text, "general_conditions_nal_28052020.pdf") {
		t.Fatal("expected general conditions template")
	}
	_, ok := ParseFromText(text)
	if ok {
		t.Fatal("general conditions must not parse as individual contract")
	}
}

func TestIsGeneralConditionsTemplate_IndividualContract(t *testing.T) {
	text := `ИНДИВИДУАЛЬНЫЕ УСЛОВИЯ ДОГОВОРА ПОТРЕБИТЕЛЬСКОГО КРЕДИТА
Сумма кредита 500 000 рублей
Срок кредита 60 месяцев
Процентная ставка 15.9% годовых`
	if IsGeneralConditionsTemplate(text, "individual.pdf") {
		t.Fatal("individual contract must not be classified as general template")
	}
	got, ok := ParseFromText(text)
	if !ok {
		t.Fatal("expected parse ok")
	}
	if got.TermMonths != 60 || got.Rate != 15.9 {
		t.Fatalf("unexpected fields: %+v", got)
	}
}
