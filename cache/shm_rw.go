package cache

import "shm-study/types"

type SHMRaw struct {
	Version int32 // SHM_VERSION   for verification
	Size    int32
	Number  int32 /* # of users total */
	Loaded  int32 /* .PASSWD has been loaded? */
	BCache  [100]types.BoardHeaderRaw
	Userid  [50]types.UserID_t
}
