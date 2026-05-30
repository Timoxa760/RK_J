package onboarding

import "testing"

func TestParseEmergencyBreakdown(t *testing.T) {
	b := parseEmergencyBreakdown("Наличными 50000, на вкладах 15000, в инвестициях 450 рублей")
	if b.Cash != 50000 || b.Deposit != 15000 || b.Investments != 450 {
		t.Fatalf("unexpected %+v", b)
	}
	if b.total() != 65450 {
		t.Fatalf("total=%v", b.total())
	}
}

func TestParse_CushionBreakdown(t *testing.T) {
	res := Parse(StepCushion, "наликом 50 тыщ, вклад 15 косарей, инвестиции 450")
	if !res.Parsed {
		t.Fatal("expected parsed")
	}
	if res.Patch.EmergencyBreakdown == nil {
		t.Fatal("expected breakdown")
	}
	b := res.Patch.EmergencyBreakdown
	if b.Cash != 50000 || b.Deposit != 15000 || b.Investments != 450 {
		t.Fatalf("unexpected breakdown %+v", b)
	}
	if res.Patch.EmergencyFund == nil || *res.Patch.EmergencyFund != 65450 {
		t.Fatalf("fund=%v", res.Patch.EmergencyFund)
	}
}
