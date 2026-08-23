package verify

import (
	"go-waste-routes/internal/service"
	"sync"
	"testing"
	"time"
)

func TestBug033_ConcurrentTaskGenerationIsIdempotent(t *testing.T) {
	s := service.NewPlanService(nil)
	day := time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC)
	got := make([]string, 10)
	var wg sync.WaitGroup
	for i := range got {
		wg.Add(1)
		go func(i int) { defer wg.Done(); got[i] = s.GenerateDailyTaskSet(33, day)[0] }(i)
	}
	wg.Wait()
	for _, value := range got[1:] {
		if value != got[0] {
			t.Fatalf("same plan/date generated different task sets: %#v", got)
		}
	}
}
