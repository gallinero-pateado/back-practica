package models

type Administrador struct {
	Id               uint   `gorm:"primaryKey;autoIncrement"`
	Firebase_usuario string `gorm:"type:text;uniqueIndex"`
	Correo           string `json:"Correo"`
	Usuario          string `json:"Usuario"`
	Rol              string `json:"ADMINISTRADOR"`
}

// TableName establece el nombre de la tabla para GORM
func (Administrador) TableName() string {
	return "Administrador"
}
