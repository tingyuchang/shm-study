package cache

import "unsafe"

var (
	Shm           *SHM
	EMPTY_SHM_RAW = SHMRaw{}
)

const SHM_RAW_SZ = unsafe.Sizeof(EMPTY_SHM_RAW)

var (
	IsTest         = false
	SHMALIGNEDSIZE = 1048576 * 1 // 4MB for x86_64, 1MB for development
)

const (
	SHM_VERSION         = 4842
	PRE_ALLOCATED_USERS = 1000
)
