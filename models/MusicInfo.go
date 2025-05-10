package models

type MusicInfo struct {
	ID       uint `gorm:"primaryKey"`
	Singer   string
	Album    string
	Name     string
	Cover    string
	Location string
}
