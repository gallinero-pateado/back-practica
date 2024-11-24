package models

import "time"

type Comentario struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TemaID            int64     `json:"tema_id"`
	UsuarioID         uint      `json:"usuario_id"`
	Contenido         string    `json:"contenido"`
	ComentarioPadreID *uint     `json:"comentario_padre_id,omitempty"`
	FechaCreacion     time.Time `json:"fecha_creacion" gorm:"autoCreateTime"`
}

// TableName establece el nombre de la tabla para GORM
func (Comentario) TableName() string {
	return "Comentario"
}
