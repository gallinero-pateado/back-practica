package models

import "time"

type Tema struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Titulo        string    `json:"titulo"`
	Contenido     string    `json:"contenido"`
	UsuarioID     uint      `json:"usuario_id"` // Este debería ser el campo donde almacenas el ID del usuario
	FechaCreacion time.Time `json:"fecha_creacion" gorm:"autoCreateTime"`
}

// TableName establece el nombre de la tabla para GORM
func (Tema) TableName() string {
	return "Tema"
}
