
package models

// Estructura para User representa un registro de la tabla usuarios
type User struct {
    ID       int    `json:"id_usuario"`           
    Nombre   string `json:"nombre"`       
    Apellido string `json:"apellido"`     
    Correo   string `json:"correo"`      
    Password string `json:"password"`    
    Tipo     string `json:"tipo_usuario"` // Tipo de usuario (paciente, médico, admin, enfermera)
}
