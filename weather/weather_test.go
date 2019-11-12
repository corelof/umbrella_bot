package weather

import "testing"

func TestTodayForecast(t *testing.T) {
	kind, prob := TodayForecast()
	if len(kind) < 1 || prob < 0 || prob > 100 {
		t.Errorf("unexpected type %s or prob %d\n", kind, prob)
	}
}
