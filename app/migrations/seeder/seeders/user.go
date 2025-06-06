package seeders

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"FP-DevOps/entity"

	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&entity.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
			return err
		}
	}

	jsonFile, err := os.Open("./migrations/seeder/json/user.json")
	if err != nil {
		return err
	}
	jsonData, _ := io.ReadAll(jsonFile)

	var listUser []entity.User
	json.Unmarshal(jsonData, &listUser)

	// only create if it does not exist
	for _, data := range listUser {
		var user entity.User
		err := db.Where(&entity.User{Username: data.Username}).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		exist := db.Find(&user, "username = ?", data.Username).RowsAffected
		if exist == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
