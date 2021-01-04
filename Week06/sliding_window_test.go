package week06

import (
	"testing"
	"time"
)

func TestISlidingCounter(t *testing.T) {

	counter := NewSlidingCounter(10)
	for _, request := range []float64{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		counter.Increment(request)
		time.Sleep(1 * time.Second)
	}

	if counter.Sum() != 45 {
		t.FailNow()
	}

	if counter.Avg() != 4.5 {
		t.FailNow()
	}

	if counter.Max() != 9 {
		t.FailNow()
	}

	t.Logf("sum is %v, max is %v, avg is %v", counter.Sum(), counter.Max(), counter.Avg())
}
