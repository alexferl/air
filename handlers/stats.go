package handlers

import (
	"net/http"
	"runtime"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/labstack/echo/v4"
)

type Stats struct {
	Vips    Vips    `json:"vips"`
	Runtime Runtime `json:"runtime"`
}

type Runtime struct {
	Alloc      uint64 `json:"alloc"`
	HeapInuse  uint64 `json:"heap_inuse"`
	Sys        uint64 `json:"sys"`
	TotalAlloc uint64 `json:"total_alloc"`
}

type Vips struct {
	Mem     VipsMem     `json:"mem"`
	Runtime VipsRuntime `json:"runtime"`
}

type VipsMem struct {
	Mem     int64 `json:"mem"`
	MemHigh int64 `json:"mem_high"`
	Files   int64 `json:"files"`
	Allocs  int64 `json:"allocs"`
}

type VipsRuntime struct {
	OperationCounts map[string]int64 `json:"operations_counts"`
}

// Stats returns runtime and vips stats
func (h *Handler) Stats(c echo.Context) error {
	var mem runtime.MemStats
	var vms vips.MemoryStats
	var rs vips.RuntimeStats
	vips.ReadVipsMemStats(&vms)
	vips.ReadRuntimeStats(&rs)
	runtime.ReadMemStats(&mem)

	s := Stats{
		Runtime: Runtime{
			Alloc:      mem.Alloc,
			HeapInuse:  mem.HeapInuse,
			Sys:        mem.Sys,
			TotalAlloc: mem.TotalAlloc,
		},
		Vips: Vips{
			Mem: VipsMem{
				Mem:     vms.Mem,
				MemHigh: vms.MemHigh,
				Files:   vms.Files,
				Allocs:  vms.Allocs,
			},
			Runtime: VipsRuntime{
				OperationCounts: rs.OperationCounts,
			},
		}}

	return c.JSON(http.StatusOK, s)
}
