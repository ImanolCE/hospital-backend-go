package handlers

import (
    "context"
    "hospital-back/config"
    "hospital-back/models"
    "github.com/gofiber/fiber/v2"
)

// Para inserta un nuevo consultorio en la BD
func CreateConsultorio(c *fiber.Ctx) error {
    dto := new(models.Consultorio)
    if err := c.BodyParser(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        `INSERT INTO consultorios (nombre, tipo, ubicacion, id_medico) VALUES ($1, $2, $3, $4)`,
        dto.Nombre, dto.Tipo, dto.Ubicacion, dto.MedicoID,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al crear consultorio"})
    }
    return c.JSON(fiber.Map{"message": "Consultorio creado correctamente"})
}

// GetConsultorios lista todos los consultorios 
func GetConsultorios(c *fiber.Ctx) error {
    rows, err := config.DB.Query(context.Background(),
        `SELECT id_consultorio, nombre, tipo, ubicacion, id_medico FROM consultorios`,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultorios"})
    }
    defer rows.Close()

    var lista []models.Consultorio
    for rows.Next() {
        var d models.Consultorio
        if err := rows.Scan(
            &d.ID, &d.Nombre, &d.Tipo,
            &d.Ubicacion, &d.MedicoID,
        ); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al leer datos"})
        }
        lista = append(lista, d)
    }
    return c.JSON(lista)
}

// GetConsultorioByID devuelve un consultorio por su ID
func GetConsultorioByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var d models.Consultorio

    err := config.DB.QueryRow(context.Background(),
        `SELECT id_consultorio, nombre, tipo, ubicacion, id_medico FROM consultorios
          WHERE id_consultorio=$1`, id,
    ).Scan(
        &d.ID, &d.Nombre, &d.Tipo,
        &d.Ubicacion, &d.MedicoID,
    )
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Consultorio no encontrado"})
    }
    return c.JSON(d)
}

// UpdateConsultorio modifica un consultorio existente
func UpdateConsultorio(c *fiber.Ctx) error {
    id := c.Params("id")
    dto := new(models.Consultorio)
    if err := c.BodyParser(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        `UPDATE consultorios SET nombre=$1, tipo=$2, ubicacion=$3, id_medico=$4
          WHERE id_consultorio=$5`,
        dto.Nombre, dto.Tipo, dto.Ubicacion, dto.MedicoID, id,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar consultorio"})
    }
    return c.JSON(fiber.Map{"message": "Consultorio actualizado correctamente"})
}

// DeleteConsultorio elimina un consultorio por su ID
func DeleteConsultorio(c *fiber.Ctx) error {
    id := c.Params("id")
    _, err := config.DB.Exec(context.Background(),
        `DELETE FROM consultorios WHERE id_consultorio=$1`, id,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar consultorio"})
    }
    return c.JSON(fiber.Map{"message": "Consultorio eliminado correctamente"})
}
