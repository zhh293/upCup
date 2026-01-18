package database

func DataAdd(p string, telephone string, text string, t string) {
	var i DataTable
	i.Package = p
	i.Telephone = telephone
	i.Type = t
	i.Text = text
	MainDB.Create(&i)
}

func DataGet(telephone string) []DataTable {
	var i []DataTable
	MainDB.Model(&DataTable{}).Where("telephone = ?", telephone).Find(&i)
	for a, b := 0, len(i)-1; a <= b; a, b = a+1, b-1 {
		i[a], i[b] = i[b], i[a]
	}
	return i
}

func DataCounts(telephone string) int {
	var i int64
	MainDB.Model(&DataTable{}).Where("telephone = ?", telephone).Count(&i)
	return int(i)
}

func SaveAudio(audio *AudioTable) error {
	return MainDB.Create(audio).Error
}

func GetAudioByID(audioID string) (*AudioTable, error) {
	var audio AudioTable
	result := MainDB.Model(&AudioTable{}).Where("audio_id = ?", audioID).First(&audio)
	if result.Error != nil {
		return nil, result.Error
	}
	return &audio, nil
}
