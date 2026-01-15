package database

type UserTable struct {
	Telephone string
	Password  string // 应当用md5形式存储
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
