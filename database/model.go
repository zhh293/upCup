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

type AudioTable struct {
	AudioID   string `json:"audioId" gorm:"column:audio_id"`
	Telephone string `json:"telephone" gorm:"column:telephone"`
	FileName  string `json:"fileName" gorm:"column:file_name"`
	FileSize  int64  `json:"fileSize" gorm:"column:file_size"`
	UploadTime int64 `json:"uploadTime" gorm:"column:upload_time"`
	Duration  int    `json:"duration" gorm:"column:duration"`
	Format    string `json:"format" gorm:"column:format"`
	AudioURL  string `json:"audioUrl" gorm:"column:audio_url"`
}

func (AudioTable) TableName() string {
	return "audio"
}

type LinkTable struct {
	Telephone1 string
	Telephone2 string
}

func (LinkTable) TableName() string {
	return "link"
}

type LinkedAccountTable struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	UserID string `gorm:"column:user_id"`
	Name   string `gorm:"column:name"`
	Email  string `gorm:"column:email"`
}

func (LinkedAccountTable) TableName() string {
	return "linked_account"
}
