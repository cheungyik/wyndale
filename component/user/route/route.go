package route

const (
	Register             = 1  // 注册（请求-响应）
	Login                = 2  // 登录（请求-响应）
	EnterGame            = 3  // 进入游戏（请求-响应）
	SpaceEntitySync      = 4  // 空间实体同步（请求-响应）
	SpaceCharactersEnter = 5  // 空间角色进入（通知）
	SpaceCharactersLeave = 6  // 空间角色离开（通知）
	EntitySync           = 7  // 实体同步（通知）
	EntityEnter          = 8  // 实体进入（通知）
	CreateCharacter      = 9  // 创建角色（请求-响应）
	DeleteCharacter      = 10 // 删除角色（请求-响应）
	CharacterList        = 11 // 角色列表（请求-响应）
)
