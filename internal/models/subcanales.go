package models

type Subcanales struct {
	Id       uint           `gorm:"primaryKey;autoIncrement"`
	Id_canal uint           `json:"id_canal"`
	Nombre   string         `json:"nombre"`
	Mensajes []MensajeCanal `json:"mensajes"`
	Usuarios []Usuario      `json:"usuarios"`
}

func (Subcanales) TableName() string {
	return "Canales"
}
