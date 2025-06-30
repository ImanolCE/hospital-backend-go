package models

// Consulta representa un registro de la tabla consultas
type Consulta struct {
    ID           int     `json:"id_consulta"`    
    PacienteID   int     `json:"id_paciente"`    // la fK de los usuarios(id_usuario)
    MedicoID     int     `json:"id_medico"`      // la fk a usuarios(id_usuario)
    ConsultorioID int    `json:"id_consultorio"` // la fk a consultorios(id_consultorio)
    HorarioID    int     `json:"id_horario"`     // la fk a horarios(id_horario)
    Diagnostico  string  `json:"diagnostico"`    // Tel diagnóstico
    Costo        float64 `json:"costo"`          
    Tipo         string  `json:"tipo"`           // Tipo de consulta, como ("Primera vez","Control")
}
