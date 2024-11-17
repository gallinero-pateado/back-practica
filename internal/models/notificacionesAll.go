package models

import (
	"time"
)

type Notificaciones_All struct {
	Id                 uint      `gorm:"primaryKey;autoIncrement"`
	Titulo             string    `json:"titulo"`
	Mensaje            string    `json:"mensaje"`
	Fecha_hora_mensaje time.Time `json:"default:CURRENT_TIMESTAMP"`
	Estado             string    `json:"estado"`
}

func (Notificaciones_All) TableName() string {
	return "Notificaciones"
}
