package types

type BoardHeaderRaw struct {
	// Require updating SHM_VERSION if BOARD_HEADER_RAW_SZ is changed.
	Brdname            BoardID_t /* bid */
	Title              BoardTitle_t
	BM                 BM_t /* BMs' userid, token '/' */
	Pad1               [3]byte
	BrdAttr            BrdAttr   /* board的屬性 */
	ChessCountry       ChessCode /* 棋國 */
	VoteLimitPosts_    uint8     /* (已停用) 連署 : 文章篇數下限 */
	VoteLimitLogins    uint8     /* 連署 : 登入次數下限 */
	Pad2_1             [1]uint8  /* (已停用) 連署 : 註冊時間限制 */
	BUpdate            Time4     /* note update time */
	PostLimitPosts_    uint8     /* (已停用) 發表文章 : 文章篇數下限 */
	PostLimitLogins    uint8     /* 發表文章 : 登入次數下限 */
	Pad2_2             [1]uint8  /* (已停用) 發表文章 : 註冊時間限制 */
	BVote              uint8     /* 正舉辦 Vote 數 */
	VTime              Time4     /* Vote close time */
	Level              PERM      /* 可以看此板的權限 */
	PermReload         Time4     /* 最後設定看板的時間 */
	Gid                Bid       /* 看板所屬的類別 ID */
	Next               [2]Bid    /* 在同一個gid下一個看板 動態產生*/
	FirstChild         [2]Bid    /* 屬於這個看板的第一個子看板 */
	Parent             Bid       /* 這個看板的 parent 看板 bid */
	ChildCount         int32     /* 有多少個child */
	NUser              int32     /* 多少人在這板 */
	PostExpire         int32     /* postexpire */
	EndGamble          Time4
	PostType           [33]byte
	PostTypeF          byte
	FastRecommendPause uint8 /* 快速連推間隔 */
	VoteLimitBadpost   uint8 /* 連署 : 劣文上限 */
	PostLimitBadpost   uint8 /* 發表文章 : 劣文上限 */
	Pad3               [3]byte
	SRexpire           Time4 /* SR Records expire time */
	Pad4               [40]byte
}
