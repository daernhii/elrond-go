package throttler

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

var durationBetweenLoadingSystemStats = time.Second

// ResourceThrottler can limit the number of go routines launched
// taking into account the machine capability
type ResourceThrottler struct {
	targetSysMem        uint64
	targetGoMem         uint64
	targetNumGoRoutines int64

	currentSysMem        uint64
	currentGoMem         uint64
	currentNumGoRoutines int64
}

func NewResourceThrottler(
	targetSysMemory uint64,
	targetGoLangMemory uint64,
	targetNumGoRoutines int64,
) (*ResourceThrottler, error) {

	//TODO checks for > 0

	rt := &ResourceThrottler{
		targetSysMem:        targetSysMemory,
		targetGoMem:         targetGoLangMemory,
		targetNumGoRoutines: targetNumGoRoutines,
	}
	go rt.takeStats()

	return rt, nil
}

func (rt *ResourceThrottler) CanProcess() bool {
	isSysMemOccupied := rt.targetSysMem <= atomic.LoadUint64(&rt.currentSysMem)
	if isSysMemOccupied {
		return false
	}

	isGoMemOccupied := rt.targetGoMem <= atomic.LoadUint64(&rt.currentGoMem)
	if isGoMemOccupied {
		return false
	}

	isMaxNumGoRoutinesReached := rt.targetNumGoRoutines <= atomic.LoadInt64(&rt.currentNumGoRoutines)
	if isMaxNumGoRoutinesReached {
		return false
	}

	return true
}

func (rt *ResourceThrottler) StartProcessing() {}

func (rt *ResourceThrottler) EndProcessing() {}

func (rt *ResourceThrottler) takeStats() {
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		atomic.StoreUint64(&rt.currentGoMem, memStats.Alloc)
		atomic.StoreUint64(&rt.currentSysMem, memStats.Sys)
		atomic.StoreInt64(&rt.currentNumGoRoutines, int64(runtime.NumGoroutine()))

		fmt.Printf("sys_mem: %d, go_mem: %d, routines: %d\n", memStats.Sys, memStats.Alloc, runtime.NumGoroutine())

		time.Sleep(durationBetweenLoadingSystemStats)
	}
}

func (rt *ResourceThrottler) IsInterfaceNil() bool {
	return rt == nil
}
