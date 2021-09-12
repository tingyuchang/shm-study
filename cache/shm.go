package cache

import (
	"encoding/binary"
	"errors"
	"reflect"
	"runtime/debug"
	"shm-study/shm"
	"shm-study/types"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

type SHM struct {
	ShmID   int
	IsNew   bool
	ShmAddr unsafe.Pointer

	Raw SHMRaw // dummy variable to get the offset and size of the shm-fields.
}

// NewSHM is to init SHM with Version and Size checked.
func NewSHM(key types.Key_t, isUseHugeTlb bool, isCreate bool) error {
	if Shm != nil {
		return nil
	}

	shmID := int(0)
	var shmAddr unsafe.Pointer
	isNew := false
	var err error

	SHMSIZE := types.Size_t(SHM_RAW_SZ)
	if SHMALIGNEDSIZE != 0 {
		SHMSIZE = types.Size_t((int(SHM_RAW_SZ)/(SHMALIGNEDSIZE) + 1) * SHMALIGNEDSIZE)
	}

	log.Infof("cache.NewSHM: SHMSIZE: %v SHM_RAW_SZ: %v SHMALIGNEDSIZE: %v", SHMSIZE, SHM_RAW_SZ, SHMALIGNEDSIZE)
	if !IsTest {
		debugShm()
	}

	size := SHMSIZE

	if isCreate {
		shmID, shmAddr, isNew, err = shm.CreateShm(key, size, isUseHugeTlb)
		if err != nil {
			return err
		}
	} else {
		shmID, shmAddr, err = shm.OpenShm(key, size, isUseHugeTlb)
		if err != nil {
			return err
		}
	}

	Shm = &SHM{
		ShmID:   shmID,
		IsNew:   isNew,
		ShmAddr: shmAddr,
	}

	if isNew {
		in_version := SHM_VERSION
		in_size := int32(SHM_RAW_SZ)
		in_number := int32(0)
		in_loaded := int32(0)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Version),
			unsafe.Sizeof(Shm.Raw.Version),
			unsafe.Pointer(&in_version),
		)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Size),
			unsafe.Sizeof(Shm.Raw.Size),
			unsafe.Pointer(&in_size),
		)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Number),
			unsafe.Sizeof(Shm.Raw.Number),
			unsafe.Pointer(&in_number),
		)
		Shm.WriteAt(
			unsafe.Offsetof(Shm.Raw.Loaded),
			unsafe.Sizeof(Shm.Raw.Loaded),
			unsafe.Pointer(&in_loaded),
		)
	}

	// version and size should be fixed.
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Version),
		unsafe.Sizeof(Shm.Raw.Version),
		unsafe.Pointer(&Shm.Raw.Version),
	)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Size),
		unsafe.Sizeof(Shm.Raw.Size),
		unsafe.Pointer(&Shm.Raw.Size),
	)

	Shm.SetBCACHEPTR(
		unsafe.Offsetof(Shm.Raw.BCache),
	)

	// verify version
	if Shm.Raw.Version != SHM_VERSION {
		log.Errorf("cache.NewSHM: version not match: key: %v Shm.Raw.Version: %v SHM_VERSION: %v isCreate: %v isNew: %v", key, Shm.Raw.Version, SHM_VERSION, isCreate, isNew)
		debug.PrintStack()
		_ = CloseSHM()
		return errors.New("ErrShmVersion")
	}
	if Shm.Raw.Size != int32(SHM_RAW_SZ) {
		log.Warnf("cache.NewSHM: size not match (version matched): key: %v Shm.Raw.Size: %v SHM_RAW_SZ: %v size: %v isCreate: %v isNew: %v", key, Shm.Raw.Size, SHM_RAW_SZ, size, isCreate, isNew)

		_ = CloseSHM()
		return errors.New("ErrShmSize")
	}

	if isCreate && !isNew {
		log.Warnf("cache.NewSHM: is expected to create, but not: key: %v", key)
	}

	log.Infof("cache.NewSHM: shm created: key: %v shmID: %v shmAddr: %v size: %v isNew: %v", key, Shm.ShmID, Shm.ShmAddr, Shm.Raw.Size, isNew)

	return nil
}

func debugShm() {
	val := reflect.ValueOf(EMPTY_SHM_RAW)
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		field := val.Type().Field(i)
		fieldName := field.Name
		offset := field.Offset

		value := val.Field(i).Interface()
		theSize := binary.Size(value)

		log.Debugf("cache.Shm.Raw: (%v/%v) %v: offset: %v size: %v", i, numField, fieldName, offset, theSize)
	}
}

//ReadAt
//
//Require precalculated offset and size and outptr to efficiently get the data.
//See tests for exact usage.
//[!!!] If we are reading from the array, make sure that have unit-size * n in the size.
//
//Params
//  offsetOfSHMRawComponent: offset from SHMRaw
//  size: size of the variable, usually can be referred from SHMRaw
//        [!!!]If we are reading from the array, make sure that have unit-size * n in the size.
//  outptr: the ptr of the object to read.
func (s *SHM) ReadAt(offsetOfSHMRawComponent uintptr, size uintptr, outptr unsafe.Pointer) {
	shm.ReadAt(s.ShmAddr, int(offsetOfSHMRawComponent), size, outptr)
}

//WriteAt
//
//Require recalculated offset and size and outptr to efficiently get the data.
//See tests for exact usage.
//[!!!]If we are reading from the array, make sure that have unit-size * n in the size.
//
//Params
//  offsetOfSHMRawComponent: offset from SHMRaw
//  size: size of the variable
//        [!!!]If we are reading from the array, make sure that have unit-size * n in the size.
//  inptr: the ptr of the object to write.
func (s *SHM) WriteAt(offsetOfSHMRawComponent uintptr, size uintptr, inptr unsafe.Pointer) {
	shm.WriteAt(s.ShmAddr, int(offsetOfSHMRawComponent), size, inptr)
}

//SetBCACHEPTR
//
//!!!Required in NewSHM (and should be set only once in NewSHM)
func (s *SHM) SetBCACHEPTR(offsetOfSHMRawComponent uintptr) {
	shm.SetBCACHEPTR(s.ShmAddr, int(offsetOfSHMRawComponent))
}

//Close
//
//XXX [WARNING] know what you are doing before using Close!.
//This is to be able to close the shared mem for the completeness of the mem-usage.
//However, in production, we create shm without the need of closing the shm.
//
//We simply use ipcrm to delete the shm if necessary.
//
//Currently used only in test.
//XXX not doing close shm to prevent opening too many shms in tests.
func CloseSHM() error {
	if !IsTest {
		return errors.New("ErrInvalidOp")
	}

	if Shm == nil {
		// Already Closed
		log.Errorf("cache.CloseSHM: already closed")
		return errors.New("ErrShmNotInit")
	}

	err := Shm.Close()
	if err != nil {
		log.Errorf("cache.CloseSHM: unable to close: e: %v", err)
		return err
	}

	Shm = nil

	log.Infof("cache.CloseSHM: done")

	return nil
}

//Close
//
//XXX [WARNING] know what you are doing before using Close!.
//This is to be able to close the shared mem for the completeness of the mem-usage.
//However, in production, we create shm without the need of closing the shm.
//
//We simply use ipcrm to delete the shm if necessary.
//
//Currently used only in test.
//XXX not doing close shm to prevent opening too many shms in tests.
func (s *SHM) Close() error {
	if !IsTest {
		return errors.New("ErrInvalidOp")
	}

	return shm.CloseShm(s.ShmID, s.ShmAddr)
}

func (s *SHM) Reset() {
	if !IsTest {
		return
	}
	s.WriteAt(
		unsafe.Offsetof(s.Raw.Userid),
		SHM_RAW_SZ-uintptr(types.INT32_SZ*2),
		unsafe.Pointer(&EMPTY_SHM_RAW.Userid),
	)
}
