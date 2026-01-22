package setting

type SettingModel struct {
	Port      string               `json:"port"`      // 服务器端口
	AIHost    string               `json:"aihost"`    // ai服务器地址
	AIPort    int                  `json:"aiport"`    // ai服务器端口默认为6666
	Debug     bool                 `json:"debug"`     // 调试模式，若为真，则AI接口假调用
	Database  DatabaseSettingModel `json:"database"`  // 数据库设置
	RedisPort int                  `json:"redisport"` // redis的端口
	JWT       JWTSettingModel      `json:"jwt"`       // JWT配置
}

type JWTSettingModel struct {
	Secret      string `json:"secret"`      // JWT密钥
	ExpiresIn   int    `json:"expires_in"`  // Token过期时间（秒），默认7200秒
	TokenType   string `json:"token_type"`  // Token类型，默认Bearer
}

type DatabaseSettingModel struct {
	TypeName string `json:"type"`   // 数据库类型,mysql,sqlite
	Source   string `json:"source"` // 数据库地址
}
