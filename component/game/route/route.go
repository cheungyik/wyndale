package route

const (
	EnterGame                  = 3 // 进入游戏（请求-响应）
	SpaceCharactersEnterNotify = 4 // 空间角色进入（通知）
	SpaceCharactersLeave       = 5 // 空间角色离开（请求）
	SpaceCharactersLeaveNotify = 6 // 空间角色离开（通知）
	SpaceEntitySync            = 7 // 空间实体同步（请求-响应）
	SpaceEntitySyncNotify      = 8 // 空间实体同步（通知）
)
