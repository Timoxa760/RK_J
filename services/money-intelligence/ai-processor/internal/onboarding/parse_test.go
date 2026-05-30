package onboarding

import "testing"

func TestParse_IncomeSlang(t *testing.T) {
	res := Parse(StepIncome, "180 тысяч зарплата и 20 косарей с аренды")
	if !res.Parsed {
		t.Fatal("expected parsed")
	}
	if res.Patch.ActiveIncome == nil || *res.Patch.ActiveIncome != 180000 {
		t.Fatalf("active=%v", res.Patch.ActiveIncome)
	}
	if res.Patch.PassiveIncome == nil || *res.Patch.PassiveIncome != 20000 {
		t.Fatalf("passive=%v", res.Patch.PassiveIncome)
	}
}

func TestParse_CushionTysh(t *testing.T) {
	res := Parse(StepCushion, "300 тыщ на вкладе")
	if !res.Parsed || res.Patch.EmergencyFund == nil || *res.Patch.EmergencyFund != 300000 {
		t.Fatalf("unexpected %+v", res)
	}
}

func TestParse_ExpenseVoice(t *testing.T) {
	res := Parse(StepExpenses, "аренда 45 тыщ")
	if !res.Parsed || len(res.Patch.FixedExpenses) != 1 || res.Patch.FixedExpenses[0].Amount != 45000 {
		t.Fatalf("unexpected %+v", res)
	}
}

func TestParse_Empty(t *testing.T) {
	res := Parse(StepIncome, "пока ничего")
	if res.Parsed {
		t.Fatal("expected not parsed")
	}
}
