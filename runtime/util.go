package runtime

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"math"
)

func Busy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}

func CalculateAllBusy(t1, t2 []cpu.TimesStat) ([]float64, error) {
	// Make sure the CPU measurements have the same length.
	if len(t1) != len(t2) {
		return nil, fmt.Errorf(
			"received two CPU counts: %d != %d",
			len(t1), len(t2),
		)
	}

	ret := make([]float64, len(t1))
	for i, t := range t2 {
		ret[i] = calculateBusy(t1[i], t)
	}
	return ret, nil
}

func calculateBusy(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := Busy(t1)
	t2All, t2Busy := Busy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}

	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}

func Size(ctx *fasthttp.RequestCtx, d int) int {
	n := ctx.QueryArgs().Peek("size")

	size := auxlib.ToInt(string(n))
	if size == 0 {
		return d
	}

	return size
}
