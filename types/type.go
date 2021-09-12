package types

type (
	Key_t  int
	Size_t uint32
)

type (
	BoardID_t    [IDLEN + 1]byte
	Bid          int32
	BidInStore   int32
	BoardTitle_t [BTLEN + 1]byte
	BM_t         [IDLEN*3 + 3]byte /* BMs' userid, token '/' */
)

type (
	UserID_t [IDLEN + 1]byte
)

type (
	BrdAttr   uint32
	ChessCode byte
	Time4     int32
	PERM      uint32
)
