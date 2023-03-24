package models

type Application struct {
	Id        uint   `gorm:"primaryKey"`
	ServerKey string `gorm:"not null"`
	ClientKey string `gorm:"not null"`
	AppName   string `gorm:"not null"`
	Callback  string `gorm:"not null"`
}
