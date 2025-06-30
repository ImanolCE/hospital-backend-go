package handlers

import (
    "context"
    "hospital-back/config"
    "hospital-back/models"
    "github.com/gofiber/fiber/v2"
)

// CreateHorario inserta un nuevo horario en la BD
func CreateHorario(c *fiber.Ctx) error {
    dto := new(models.Horario)
    if err := c.BodyParser(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        `INSERT INTO horarios (id_consultorio, turno, id_medico) VALUES ($1, $2, $3)`,
        dto.ConsultorioID, dto.Turno, dto.MedicoID,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al crear horario"})
    }
    return c.JSON(fiber.Map{"message": "Horario creado correctamente"})
}

// GetHorarios devuelve todos los horarios
func GetHorarios(c *fiber.Ctx) error {
    rows, err := config.DB.Query(context.Background(),
        `SELECT id_horario, id_consultorio, turno, id_medico FROM horarios`,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener horarios"})
    }
    defer rows.Close()

    var lista []models.Horario
    for rows.Next() {
        var h models.Horario
        if err := rows.Scan(
            &h.ID, &h.ConsultorioID, &h.Turno, &h.MedicoID,
        ); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al leer datos"})
        }
        lista = append(lista, h)
    }
    return c.JSON(lista)
}

// GetHorarioByID devuelve un horario por su ID
func GetHorarioByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var h models.Horario

    err := config.DB.QueryRow(context.Background(),
        `SELECT id_horario, id_consultorio, turno, id_medico FROM horarios
          WHERE id_horario=$1`, id,
    ).Scan(
        &h.ID, &h.ConsultorioID,
        &h.Turno, &h.MedicoID,
    )
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Horario no encontrado"})
    }
    return c.JSON(h)
}

// UpdateHorario modifica un horario existente
func UpdateHorario(c *fiber.Ctx) error {
    id := c.Params("id")
    dto := new(models.Horario)
    if err := c.BodyParser(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        `UPDATE horarios SET id_consultorio=$1, turno=$2, id_medico=$3
          WHERE id_horario=$4`,
        dto.ConsultorioID, dto.Turno, dto.MedicoID, id,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar horario"})
    }
    return c.JSON(fiber.Map{"message": "Horario actualizado correctamente"})
}

// DeleteHorario elimina un horario por su ID
func DeleteHorario(c *fiber.Ctx) error {
    id := c.Params("id")
    _, err := config.DB.Exec(context.Background(),
        `DELETE FROM horarios WHERE id_horario=$1`, id,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar horario"})
    }
    return c.JSON(fiber.Map{"message": "Horario eliminado correctamente"})
}
