package models

import (
	"time"
)

type NotificacionesFront struct {
	Id                 uint      `gorm:"primaryKey;autoIncrement"`
	Id_mensaje         int       `json:"id_mensaje"`
	Id_receptor        int       `json:"id_receptor"`
	Fecha_hora_mensaje time.Time `json:"default:CURRENT_TIMESTAMP"`
	Estado             string    `json:"estado"`
}

func (NotificacionesFront) TableName() string {
	return "NotificacionesFront"
}
