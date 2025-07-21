package handlers

import (
    "context"
    "hospital-back/config"
    "hospital-back/models"
    "hospital-back/utils"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/xeipuuv/gojsonschema"
)

const horarioSchema = `{
  "type": "object",
  "required": ["id_consultorio", "turno", "id_medico"],
  "properties": {
    "id_consultorio": { "type": "integer" },
    "turno": { "type": "string" },
    "id_medico": { "type": "integer" }
  },
  "additionalProperties": false
}`

// CreateHorario inserta un nuevo horario en la BD
func CreateHorario(c *fiber.Ctx) error {
    body := c.Body()
    schemaLoader := gojsonschema.NewStringLoader(horarioSchema)
    docLoader := gojsonschema.NewBytesLoader(body)
    result, err := gojsonschema.Validate(schemaLoader, docLoader)
    if err != nil {
        return utils.ResponseError(c, 400, "SC00", "Error validando JSON")
    }
    if !result.Valid() {
        errs := make([]string, 0, len(result.Errors()))
        for _, desc := range result.Errors() {
            errs = append(errs, desc.String())
        }
        return utils.ResponseError(c, 400, "SC01", "Esquema inválido: "+strings.Join(errs, "; "))
    }

    dto := new(models.Horario)
    if err := c.BodyParser(dto); err != nil {
        return utils.ResponseError(c, 400, "F01", "Datos inválidos")
    }

    _, err = config.DB.Exec(context.Background(),
        `INSERT INTO horarios (id_consultorio, turno, id_medico) VALUES ($1, $2, $3)`,
        dto.ConsultorioID, dto.Turno, dto.MedicoID,
    )
    if err != nil {
        return utils.ResponseError(c, 500, "F02", "Error al crear horario")
    }
    return utils.ResponseSuccess(c, 201, "S01", []fiber.Map{
        {"message": "Horario creado correctamente"},
    })
}

// GetHorarios devuelve todos los horarios
func GetHorarios(c *fiber.Ctx) error {
    rows, err := config.DB.Query(context.Background(),
        `SELECT id_horario, id_consultorio, turno, id_medico FROM horarios`,
    )
    if err != nil {
        return utils.ResponseError(c, 500, "F03", "Error al obtener horarios")
    }
    defer rows.Close()

    var lista []models.Horario
    for rows.Next() {
        var h models.Horario
        if err := rows.Scan(&h.ID, &h.ConsultorioID, &h.Turno, &h.MedicoID); err != nil {
            return utils.ResponseError(c, 500, "F04", "Error al leer datos")
        }
        lista = append(lista, h)
    }
    return utils.ResponseSuccess(c, 200, "S02", lista)
}

// GetHorarioByID devuelve un horario por su ID
func GetHorarioByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var h models.Horario

    err := config.DB.QueryRow(context.Background(),
        `SELECT id_horario, id_consultorio, turno, id_medico FROM horarios
         WHERE id_horario=$1`, id,
    ).Scan(&h.ID, &h.ConsultorioID, &h.Turno, &h.MedicoID)
    if err != nil {
        return utils.ResponseError(c, 404, "F05", "Horario no encontrado")
    }

    return utils.ResponseSuccess(c, 200, "S03", []fiber.Map{{
        "id":             h.ID,
        "id_consultorio": h.ConsultorioID,
        "turno":          h.Turno,
        "id_medico":      h.MedicoID,
    }})
}

// UpdateHorario modifica un horario existente
func UpdateHorario(c *fiber.Ctx) error {
    id := c.Params("id")

    body := c.Body()
    schemaLoader := gojsonschema.NewStringLoader(horarioSchema)
    docLoader := gojsonschema.NewBytesLoader(body)
    result, err := gojsonschema.Validate(schemaLoader, docLoader)
    if err != nil {
        return utils.ResponseError(c, 400, "SC10", "Error validando JSON")
    }
    if !result.Valid() {
        errs := make([]string, 0, len(result.Errors()))
        for _, desc := range result.Errors() {
            errs = append(errs, desc.String())
        }
        return utils.ResponseError(c, 400, "SC11", "Esquema inválido: "+strings.Join(errs, "; "))
    }

    dto := new(models.Horario)
    if err := c.BodyParser(dto); err != nil {
        return utils.ResponseError(c, 400, "F06", "Datos inválidos")
    }

    _, err = config.DB.Exec(context.Background(),
        `UPDATE horarios SET id_consultorio=$1, turno=$2, id_medico=$3
         WHERE id_horario=$4`,
        dto.ConsultorioID, dto.Turno, dto.MedicoID, id,
    )
    if err != nil {
        return utils.ResponseError(c, 500, "F07", "Error al actualizar horario")
    }

    return utils.ResponseSuccess(c, 200, "S04", []fiber.Map{
        {"message": "Horario actualizado correctamente"},
    })
}

// DeleteHorario elimina un horario por su ID
func DeleteHorario(c *fiber.Ctx) error {
    id := c.Params("id")
    _, err := config.DB.Exec(context.Background(),
        `DELETE FROM horarios WHERE id_horario=$1`, id,
    )
    if err != nil {
        return utils.ResponseError(c, 500, "F08", "Error al eliminar horario")
    }
    return utils.ResponseSuccess(c, 200, "S05", []fiber.Map{
        {"message": "Horario eliminado correctamente"},
    })
}
