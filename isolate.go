package v8go

// #include "v8go.h"
import "C"

import (
	"runtime"
	"sync"
)

var v8once sync.Once

// Isolate is a JavaScript VM instance with its own heap and
// garbage collector. Most applications will create one isolate
// with many V8 contexts for execution.
type Isolate struct {
	ptr C.IsolatePtr
}

// HeapStatistics represents V8 isolate heap statistics
type HeapStatistics struct {
	TotalHeapSize            uint64
	TotalHeapSizeExecutable  uint64
	TotalPhysicalSize        uint64
	TotalAvailableSize       uint64
	UsedHeapSize             uint64
	HeapSizeLimit            uint64
	MallocedMemory           uint64
	ExternalMemory           uint64
	PeakMallocedMemory       uint64
	NumberOfNativeContexts   uint64
	NumberOfDetachedContexts uint64
}

// NewIsolate creates a new V8 isolate. Only one thread may access
// a given isolate at a time, but different threads may access
// different isolates simultaneously.
func NewIsolate() (*Isolate, error) {
	v8once.Do(func() {
		C.Init()
	})
	iso := &Isolate{C.NewIsolate()}
	runtime.SetFinalizer(iso, (*Isolate).finalizer)
	// TODO: [RC] catch any C++ exceptions and return as error
	return iso, nil
}

// TerminateExecution terminates forcefully the current thread
// of JavaScript execution in the given isolate.
func (i *Isolate) TerminateExecution() {
	C.IsolateTerminateExecution(i.ptr)
}

// GetHeapStatistics returns heap statistics for an isolate.
func (i *Isolate) GetHeapStatistics() HeapStatistics {
	hs := C.IsolationGetHeapStatistics(i.ptr)

	return HeapStatistics{
		TotalHeapSize:            uint64(hs.total_heap_size),
		TotalHeapSizeExecutable:  uint64(hs.total_heap_size_executable),
		TotalPhysicalSize:        uint64(hs.total_physical_size),
		TotalAvailableSize:       uint64(hs.total_available_size),
		UsedHeapSize:             uint64(hs.used_heap_size),
		HeapSizeLimit:            uint64(hs.heap_size_limit),
		MallocedMemory:           uint64(hs.malloced_memory),
		ExternalMemory:           uint64(hs.external_memory),
		PeakMallocedMemory:       uint64(hs.peak_malloced_memory),
		NumberOfNativeContexts:   uint64(hs.number_of_native_contexts),
		NumberOfDetachedContexts: uint64(hs.number_of_detached_contexts),
	}
}

func (i *Isolate) finalizer() {
	C.IsolateDispose(i.ptr)
	i.ptr = nil
	runtime.SetFinalizer(i, nil)
}
