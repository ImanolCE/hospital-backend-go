package models

// Horario que representa un registro de la tabla horarios
type Horario struct {
    ID            int    `json:"id_horario"`    
    ConsultorioID int    `json:"id_consultorio"` 
    Turno         string `json:"turno"`          // el turno como ("Matutino", "Vespertino")
    MedicoID      int    `json:"id_medico"`      // l a fk es aa usuarios(id_usuario) como tipo_medico
}
