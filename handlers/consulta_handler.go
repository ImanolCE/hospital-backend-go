package handlers

import (
    "context"
    "hospital-back/config"
    "hospital-back/models"
    "github.com/gofiber/fiber/v2"
)

// CreateConsulta inserta una nueva consulta en la BD
func CreateConsulta(c *fiber.Ctx) error {
    dto := new(models.Consulta)
    if err := c.BodyParser(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        `INSERT INTO consultas (id_paciente, id_medico, id_consultorio, id_horario, diagnostico, costo, tipo)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        dto.PacienteID, dto.MedicoID, dto.ConsultorioID, dto.HorarioID, dto.Diagnostico, dto.Costo, dto.Tipo,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al crear consulta"})
    }
    return c.JSON(fiber.Map{"message": "Consulta creada correctamente"})
}

// GetConsultas devuelve todas las consultas
func GetConsultas(c *fiber.Ctx) error {
    rows, err := config.DB.Query(context.Background(),
        `SELECT id_consulta, id_paciente, id_medico, id_consultorio, id_horario, diagnostico, costo, tipo
           FROM consultas`,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultas"})
    }
    defer rows.Close()

    var lista []models.Consulta
    for rows.Next() {
        var d models.Consulta
        if err := rows.Scan(
            &d.ID, &d.PacienteID, &d.MedicoID, &d.ConsultorioID, &d.HorarioID,
            &d.Diagnostico, &d.Costo, &d.Tipo,
        ); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al leer datos"})
        }
        lista = append(lista, d)
    }
    return c.JSON(lista)
}

// GetConsultaByID devuelve una consulta por su ID
func GetConsultaByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var d models.Consulta

    err := config.DB.QueryRow(context.Background(),
        `SELECT id_consulta, id_paciente, id_medico, id_consultorio, id_horario, diagnostico, costo, tipo
           FROM consultas
          WHERE id_consulta=$1`, id,
    ).Scan(
        &d.ID, &d.PacienteID, &d.MedicoID, &d.ConsultorioID, &d.HorarioID,
        &d.Diagnostico, &d.Costo, &d.Tipo,
    )
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Consulta no encontrada"})
    }
    return c.JSON(d)
}

// UpdateConsulta modifica una consulta existente (rutas protegidas)
func UpdateConsulta(c *fiber.Ctx) error {
    id := c.Params("id")
    dto := new(models.Consulta)
    if err := c.BodyParser(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        `UPDATE consultas
            SET id_paciente=$1,id_medico=$2, id_consultorio=$3, id_horario=$4, diagnostico=$5, costo=$6, tipo=$7
          WHERE id_consulta=$8`,
        dto.PacienteID, dto.MedicoID, dto.ConsultorioID, dto.HorarioID, dto.Diagnostico, dto.Costo, dto.Tipo, id,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar consulta"})
    }
    return c.JSON(fiber.Map{"message": "Consulta actualizada correctamente"})
}

// DeleteConsulta elimina una consulta por su ID
func DeleteConsulta(c *fiber.Ctx) error {
    id := c.Params("id")
    _, err := config.DB.Exec(context.Background(),
        `DELETE FROM consultas WHERE id_consulta=$1`, id,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar consulta"})
    }
    return c.JSON(fiber.Map{"message": "Consulta eliminada correctamente"})
}
