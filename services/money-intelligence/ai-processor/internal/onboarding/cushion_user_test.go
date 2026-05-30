package onboarding

import "testing"

func TestParse_Cushion_UserPhrase(t *testing.T) {
	res := Parse(StepCushion, "15 тысяч наличкой, 80 тысяч на вкладах и 10 тысяч на инвестициях")
	if !res.Parsed {
		t.Fatal("expected parsed")
	}
	if res.Patch.EmergencyFund == nil || *res.Patch.EmergencyFund != 105000 {
		t.Fatalf("fund=%v", res.Patch.EmergencyFund)
	}
	b := res.Patch.EmergencyBreakdown
	if b == nil {
		t.Fatal("expected breakdown")
	}
	if b.Cash != 15000 || b.Deposit != 80000 || b.Investments != 10000 {
		t.Fatalf("breakdown=%+v", b)
	}
}
