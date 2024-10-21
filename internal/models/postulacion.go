package models

import (
	"time"
)

type Postulacion struct {
	Id                    uint      `gorm:"primaryKey;autoIncrement"`
	Id_usuario            int       `json:"id_usuario"`
	Id_practica           int       `json:"id_practica"`
	Fecha_postulacion     time.Time `json:"default:CURRENT_TIMESTAMP"`
	Mensaje               string    `json:"mensaje"`
	Id_estado_postulacion int       `json:"id_estado_publicacion"`
	Verificación          bool      `json:"verificacion"`
}

// TableName establece el nombre de la tabla para GORM
func (Postulacion) TableName() string {
	return "postulacion"
}
