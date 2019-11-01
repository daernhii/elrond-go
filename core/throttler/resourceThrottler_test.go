package throttler

import (
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResourceThrottler_CanProcessTargetGoMemLowerThanCurrentShouldRetFalse(t *testing.T) {
	t.Parallel()

	rt, _ := NewResourceThrottler(math.MaxUint64, 1, math.MaxInt64)
	time.Sleep(time.Second)

	assert.False(t, rt.CanProcess())
}

func TestResourceThrottler_CanProcessTargetSysMemLowerThanCurrentShouldRetFalse(t *testing.T) {
	t.Parallel()

	rt, _ := NewResourceThrottler(1, math.MaxUint64, math.MaxInt64)
	time.Sleep(time.Second)

	assert.False(t, rt.CanProcess())
}

func TestResourceThrottler_CanProcessTargetGoRoutinesLowerThanCurrentShouldRetFalse(t *testing.T) {
	t.Parallel()

	rt, _ := NewResourceThrottler(math.MaxUint64, math.MaxInt64, 1)
	time.Sleep(time.Second)

	assert.False(t, rt.CanProcess())
}

func TestResourceThrottler_CanProcessShouldRetTrue(t *testing.T) {
	t.Parallel()

	rt, _ := NewResourceThrottler(math.MaxUint64, math.MaxUint64, math.MaxInt64)
	time.Sleep(time.Second)

	assert.True(t, rt.CanProcess())
}

func Benchmark_TakeMachineStats(b *testing.B) {
	var memStats runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&memStats)
	}
}
