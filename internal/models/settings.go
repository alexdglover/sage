package models

import (
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	LaunchBrowserOnStartup bool `gorm:"default:true"`
}

type SettingsRepository struct {
	DB *gorm.DB
}

// Defining a Settings singleton to hold application settings
var settingsSingleton Settings

func (sr *SettingsRepository) GetSettings() (settings *Settings, err error) {
	if settingsSingleton.ID == 0 {
		err = sr.PopulateSettings()
		return &settingsSingleton, err
	}
	return &settingsSingleton, nil
}

func (sr *SettingsRepository) PopulateSettings() (err error) {
	var settings Settings
	result := sr.DB.First(&settings)
	if result.Error != nil {
		return result.Error
	}
	settingsSingleton = settings
	return nil
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (sr *SettingsRepository) Save(settings *Settings) (err error) {
	result := sr.DB.Save(settings)
	return result.Error
}
