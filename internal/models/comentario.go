package models

import "time"

type Comentario struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TemaID        int64     `json:"tema_id"`    // Asumiendo que el ID del tema es int64
	UsuarioID     uint      `json:"usuario_id"` // Debe coincidir con el tipo del ID en Usuario
	Contenido     string    `json:"contenido"`
	FechaCreacion time.Time `json:"fecha_creacion" gorm:"autoCreateTime"`
}

// TableName establece el nombre de la tabla para GORM
func (Comentario) TableName() string {
	return "Comentario"
}
