package models

type Canales struct {
	Id         uint         `gorm:"primaryKey;autoIncrement"`
	Nombre     string       `json:"nombre"`
	Subcanales []Subcanales `json:"subcanales"`
}

// TableName establece el nombre de la tabla para GORM
func (Canales) TableName() string {
	return "Canales"
}
