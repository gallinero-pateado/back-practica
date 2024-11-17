package models

type MensajeCanal struct {
	Nombre    string `json:"nombre"`
	Contenido string `json:"contenido"`
	CanalID   uint   `json:"canal_id"`
}
