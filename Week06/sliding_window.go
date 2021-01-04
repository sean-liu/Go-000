package week06

import (
	"sync"
	"time"
)

// ISlidingCounter uses to do counting and averaging in last N second window
type ISlidingCounter interface {
	Sum() float64
	Max() float64
	Avg() float64
	Increment(i float64)
}

type slidingCounter struct {
	windows  map[int64]*window
	interval int64
	mutex    *sync.RWMutex
}

type window struct {
	Value float64
}

// NewSlidingCounter return a new n-second interval sliding counter.
func NewSlidingCounter(interval int64) ISlidingCounter {
	return &slidingCounter{
		windows:  make(map[int64]*window),
		interval: interval,
		mutex:    &sync.RWMutex{},
	}
}

// getCurrentWindow return current window of time.Now().
func (sc *slidingCounter) getCurrentWindow() *window {
	now := time.Now().Unix()

	if w, found := sc.windows[now]; found {
		return w
	}

	result := &window{}
	sc.windows[now] = result
	return result
}

// removeOldWindows remove 10 second before windows.
func (sc *slidingCounter) removeOldWindows() {
	notExpired := time.Now().Unix() - sc.interval

	for timestamp := range sc.windows {
		if timestamp <= notExpired {
			delete(sc.windows, timestamp)
		}
	}
}

// Increment add i to current window.
func (sc *slidingCounter) Increment(i float64) {
	if i == 0 {
		return
	}
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	w := sc.getCurrentWindow()
	w.Value += i
	sc.removeOldWindows()
}

// Sum sums the values over the windows in the last N seconds.
func (sc *slidingCounter) Sum() float64 {
	now := time.Now().Unix()

	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	var sum float64
	for timestamp, window := range sc.windows {
		if timestamp >= now-sc.interval {
			sum += window.Value
		}
	}

	return sum
}

// Max returns the maximum value seen in the last N seconds.
func (sc *slidingCounter) Max() float64 {

	now := time.Now().Unix()

	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	var max float64

	for timestamp, window := range sc.windows {
		if timestamp >= now-sc.interval {
			if window.Value > max {
				max = window.Value
			}
		}
	}
	return max
}

// Avg returns the average value over the windows in last N seconds.
func (sc *slidingCounter) Avg() float64 {
	return sc.Sum() / float64(sc.interval)
}
