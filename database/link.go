package database

import (
	"errors"

	"gorm.io/gorm"
)

// LinkExsitSingle 检查单向关联是否存在 (t1 -> t2)
func LinkExsitSingle(t1 string, t2 string) bool {
	var count int64
	MainDB.Model(&LinkTable{}).Where("telephone1 = ? AND telephone2 = ?", t1, t2).Count(&count)
	return count > 0
}

// LinkExsit 检查双向关联是否都存在
func LinkExsit(t1 string, t2 string) bool {
	return LinkExsitSingle(t1, t2) && LinkExsitSingle(t2, t1)
}

// LinkGetAll 获取某人关联的所有号码 (底层权限)
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
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	TargetTelephone string `json:"targetTelephone"`
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
			ID:              record.ID,
			Name:            record.Name,
			TargetTelephone: record.TargetTelephone,
		})
	}
	return linkedAccounts, nil
}

func LinkedAccountAdd(userID string, userTelephone string, name string, targetTelephone string) (*LinkedAccount, error) {
	// 0. 检查目标用户是否存在
	if !UserExist(targetTelephone) {
		return nil, errors.New("target user not found")
	}

	// 使用事务保证原子性
	err := MainDB.Transaction(func(tx *gorm.DB) error {
		// 1. 检查是否已经存在于展示列表 (LinkedAccountTable)
		var count int64
		if err := tx.Model(&LinkedAccountTable{}).Where("user_id = ? AND target_telephone = ?", userID, targetTelephone).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("account already linked")
		}

		// 2. 建立底层权限关联 (LinkTable) - 双向
		// 检查并插入 t1 -> t2
		var c1 int64
		tx.Model(&LinkTable{}).Where("telephone1 = ? AND telephone2 = ?", userTelephone, targetTelephone).Count(&c1)
		if c1 == 0 {
			if err := tx.Create(&LinkTable{Telephone1: userTelephone, Telephone2: targetTelephone}).Error; err != nil {
				return err
			}
		}

		// 检查并插入 t2 -> t1
		var c2 int64
		tx.Model(&LinkTable{}).Where("telephone1 = ? AND telephone2 = ?", targetTelephone, userTelephone).Count(&c2)
		if c2 == 0 {
			if err := tx.Create(&LinkTable{Telephone1: targetTelephone, Telephone2: userTelephone}).Error; err != nil {
				return err
			}
		}

		// 3. 添加到展示列表
		record := LinkedAccountTable{
			UserID:          userID,
			Name:            name,
			TargetTelephone: targetTelephone,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 事务成功后查询并返回刚创建的记录（为了获取 ID）
	var record LinkedAccountTable
	MainDB.Where("user_id = ? AND target_telephone = ?", userID, targetTelephone).First(&record)

	return &LinkedAccount{
		ID:              record.ID,
		Name:            record.Name,
		TargetTelephone: record.TargetTelephone,
	}, nil
}

func LinkedAccountRemove(userID string, userTelephone string, id uint) error {
	return MainDB.Transaction(func(tx *gorm.DB) error {
		var record LinkedAccountTable
		// 查找记录
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
			return err
		}

		// 1. 删除底层权限关联 (LinkTable) - 双向
		// 删除 t1 -> t2
		if err := tx.Where("telephone1 = ? AND telephone2 = ?", userTelephone, record.TargetTelephone).Delete(&LinkTable{}).Error; err != nil {
			return err
		}
		// 删除 t2 -> t1
		if err := tx.Where("telephone1 = ? AND telephone2 = ?", record.TargetTelephone, userTelephone).Delete(&LinkTable{}).Error; err != nil {
			return err
		}

		// 2. 删除展示列表记录
		if err := tx.Delete(&record).Error; err != nil {
			return err
		}

		return nil
	})
}
