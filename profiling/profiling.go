package profiling

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

type Profiler struct {
	outPath     string
	samplesFile string
	cleanup     bool
	cpuprof     *os.File
	memprof     *os.File
}

func (b *Profiler) StartProfiling() error {
	var err error

	// Start CPU profiling.
	b.cpuprof, err = os.Create(filepath.Join(b.outPath, "cpu.prof"))
	if err != nil {
		return fmt.Errorf("bench: could not create cpu profile: %v", err)
	}
	if err := pprof.StartCPUProfile(b.cpuprof); err != nil {
		return fmt.Errorf("bench: could not start CPU profile: %v", err)
	}

	// Start memory profiling.
	b.memprof, err = os.Create(filepath.Join(b.outPath, "mem.prof"))
	if err != nil {
		return fmt.Errorf("bench: could not create memory profile: %v", err)
	}
	runtime.MemProfileRate = 64 * 1024

	return nil
}

func (b *Profiler) StopProfiling() error {
	if b.cpuprof != nil {
		pprof.StopCPUProfile()
		b.cpuprof.Close()
		b.cpuprof = nil
	}
	if b.memprof != nil {
		if err := pprof.Lookup("heap").WriteTo(b.memprof, 0); err != nil {
			return fmt.Errorf("error writing mem profile: %v", err)
		}
		b.memprof.Close()
		b.memprof = nil
	}
	return nil
}