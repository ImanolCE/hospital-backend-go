
package models

// Consultorio que representa a un registro de la tala de consultorios
type Consultorio struct {
    ID        int    `json:"id_consultorio"` 
    Nombre    string `json:"nombre"`         // Nombre del consultorio
    Tipo      string `json:"tipo"`           // Tipo como: (“general”, “pediatría”)
    Ubicacion string `json:"ubicacion"`      // Ubicación física
    MedicoID  int    `json:"id_medico"`      
}

