package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"shm-study/cache"
	"shm-study/types"
	"unsafe"
)

func main() {
	err := cache.NewSHM(types.Key_t(1229), false,true)
	if err != nil {
		log.Errorf("unable to init SHM: e: %v", err)
	}

	number := int32(0)
	cache.Shm.ReadAt(
		unsafe.Offsetof(cache.Shm.Raw.Number),
		unsafe.Sizeof(cache.Shm.Raw.Number),
		unsafe.Pointer(&number),
	)
	fmt.Println(number)
}
