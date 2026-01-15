package database

type UserTable struct {
	Telephone string
	Password  string // 使用bcrypt加密
	UserID    string // userId = "user_" + telephone
	Nickname  string
	Avatar    string
	Email     string
}

func (UserTable) TableName() string {
	return "user"
}

type DataTable struct {
	Telephone string `json:"telephone"`
	Text      string `json:"text"`
	Package   string `json:"package"`
	Type      string `json:"type"`
}

func (DataTable) TableName() string {
	return "data"
}

type LinkTable struct {
	Telephone1 string
	Telephone2 string
}

func (LinkTable) TableName() string {
	return "link"
}
