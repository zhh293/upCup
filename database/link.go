package database

import (
	"errors"
	"sync"
)

func LinkAdd(t1 string, t2 string) error {
	if LinkExsit(t1, t2) {
		return errors.New("exist")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if !LinkExsitSingle(t1, t2) {
			var i LinkTable
			i.Telephone1 = t1
			i.Telephone2 = t2
			MainDB.Create(&i)
		}
	}()
	go func() {
		defer wg.Done()
		if !LinkExsitSingle(t2, t1) {
			var i LinkTable
			i.Telephone1 = t2
			i.Telephone2 = t1
			MainDB.Create(&i)
		}
	}()
	wg.Wait()
	return nil
}

func LinkExsitSingle(t1 string, t2 string) bool {
	var i []LinkTable
	MainDB.Model(&LinkTable{}).Where("telephone1 = ?", t1).Find(&i)
	for _, cnt := range i {
		if cnt.Telephone2 == t2 {
			return true
		}
	}
	return false
}

func LinkExsit(t1 string, t2 string) bool {
	var b1, b2 bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		b1 = LinkExsitSingle(t1, t2)
	}()
	go func() {
		defer wg.Done()
		b2 = LinkExsitSingle(t2, t1)
	}()
	wg.Wait()
	return b1 && b2
}

func LinkGetAll(telephone string) []string {
	var i []string
	var cnt []LinkTable
	MainDB.Model(&LinkTable{}).Where("telephone1 = ?", telephone).Find(&cnt)
	for _, i2 := range cnt {
		i = append(i, i2.Telephone2)
	}
	return i
}

type LinkedAccount struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func LinkedAccountList(userID string) ([]LinkedAccount, error) {
	var records []LinkedAccountTable
	result := MainDB.Model(&LinkedAccountTable{}).Where("user_id = ?", userID).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	linkedAccounts := make([]LinkedAccount, 0, len(records))
	for _, record := range records {
		linkedAccounts = append(linkedAccounts, LinkedAccount{
			ID:    record.ID,
			Name:  record.Name,
			Email: record.Email,
		})
	}
	return linkedAccounts, nil
}

func LinkedAccountAdd(userID string, name string, email string) (*LinkedAccount, error) {
	record := LinkedAccountTable{
		UserID: userID,
		Name:   name,
		Email:  email,
	}
	result := MainDB.Create(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return &LinkedAccount{
		ID:    record.ID,
		Name:  record.Name,
		Email: record.Email,
	}, nil
}

func LinkedAccountRemove(userID string, id uint) error {
	var record LinkedAccountTable
	result := MainDB.Where("id = ? AND user_id = ?", id, userID).First(&record)
	if result.Error != nil {
		return result.Error
	}
	deleteResult := MainDB.Delete(&record)
	return deleteResult.Error
}
