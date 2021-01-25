package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RoutineManager provide graceful exit for running rountine
type RoutineManager interface {
	Register(eachFunc func() bool)
	Start() (chan struct{}, context.CancelFunc)
	WaitForSystemSignal(function func())
}

type routineManager struct {
	CheckInterval time.Duration
	funcs         []func() bool
}

// NewDefaultRoutineManager new instance of routine manager
func NewDefaultRoutineManager() RoutineManager {
	return &routineManager{
		CheckInterval: 2 * time.Second,
		funcs:         make([]func() bool, 0),
	}
}

// Register register function to it true : need to exit
func (rm *routineManager) Register(eachFunc func() bool) {
	rm.funcs = append(rm.funcs, eachFunc)
}

// Start starts all functions in it
func (rm *routineManager) Start() (chan struct{}, context.CancelFunc) {
	doneChan := make(chan struct{})
	waitChan := make(chan struct{}, len(rm.funcs))
	ctx, cancel := context.WithCancel(context.Background())
	for _, eachFunc := range rm.funcs {
		go func(ef func() bool, ectx context.Context) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("func fatal err is %s", err)
				}
			}()
			for {
				requireToExit := ef()
				if requireToExit {
					waitChan <- struct{}{}
					return
				}
				select {
				case <-ectx.Done():
					waitChan <- struct{}{}
					return
				case <-time.After(rm.CheckInterval):
					break
				}
			}
		}(eachFunc, ctx)
	}

	go func() {
		for i := 0; i < len(rm.funcs); i++ {
			<-waitChan
			fmt.Printf("%vth wait signal received!\r\n", i+1)
		}
		doneChan <- struct{}{}
	}()

	return doneChan, cancel
}

func (rm *routineManager) WaitForSystemSignal(function func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	function()
}
