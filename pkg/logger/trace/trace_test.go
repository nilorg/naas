package trace

import "testing"

// func TestNextSpanIDTimer(t *testing.B) {
// 	t.ResetTimer()
// 	t.StartTimer()

// }

func TestNextSpanID(t *testing.T) {
	spanID := NextSpanID("")
	t.Logf("next_spanID:%s", spanID)
	spanID = NextSpanID(spanID)
	t.Logf("next_spanID:%s", spanID)
	spanID = NextSpanID(spanID)
	t.Logf("next_spanID:%s", spanID)
	spanID = NextSpanID(spanID)
	t.Logf("next_spanID:%s", spanID)
}

// func TestStartSpanID(t *testing.T) {
// 	spanID := StartSpanID("")
// 	t.Logf("start_spanID:%s", spanID)
// 	spanID = StartSpanID(spanID)
// 	t.Logf("start_spanID:%s", spanID)
// 	spanID = StartSpanID(spanID)
// 	t.Logf("start_spanID:%s", spanID)
// 	spanID = StartSpanID(spanID)
// 	t.Logf("start_spanID:%s", spanID)
// }

// func TestAllSpanID(t *testing.T) {
// 	spanID := StartSpanID("")
// 	t.Logf("start_spanID:%s", spanID)
// 	spanID = StartSpanID(spanID)
// 	t.Logf("start_spanID:%s", spanID)
// 	spanID = NextSpanID(spanID)
// 	t.Logf("next_spanID:%s", spanID)
// 	spanID = NextSpanID(spanID)
// 	t.Logf("next_spanID:%s", spanID)
// }
