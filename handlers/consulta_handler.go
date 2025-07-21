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

const consultaSchema = `{
  "type": "object",
  "required": ["id_paciente","id_medico","id_consultorio","id_horario","diagnostico","costo","tipo"],
  "properties": {
    "id_paciente":    { "type": "integer" },
    "id_medico":      { "type": "integer" },
    "id_consultorio": { "type": "integer" },
    "id_horario":     { "type": "integer" },
    "diagnostico":    { "type": "string" },
    "costo":          { "type": "number" },
    "tipo":           { "type": "string" }
  },
  "additionalProperties": false
}`


// CreateConsulta inserta una nueva consulta en la base de datos
func CreateConsulta(c *fiber.Ctx) error {
	body := c.Body()
	schemaLoader := gojsonschema.NewStringLoader(consultaSchema)
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

	consulta := new(models.Consulta)
	if err := c.BodyParser(consulta); err != nil {
		return utils.ResponseError(c, 400, "SC02", "Datos inválidos")
	}

	_, err = config.DB.Exec(context.Background(),
		`INSERT INTO consultas (id_paciente, id_medico, id_consultorio, id_horario, diagnostico, costo, tipo)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		consulta.PacienteID, consulta.MedicoID, consulta.ConsultorioID, consulta.HorarioID, consulta.Diagnostico,
		consulta.Costo, consulta.Tipo)

	if err != nil {
		return utils.ResponseError(c, 500, "SC03", "Error al insertar consulta")
	}

	return utils.ResponseSuccess(c, 201, "SS01", []fiber.Map{{"message": "Consulta creada"}})
}


// GetConsultas devuelve todas las consultas
func GetConsultas(c *fiber.Ctx) error {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id_consulta, id_paciente, id_medico, id_consultorio, id_horario, diagnostico, costo, tipo FROM consultas`)
	if err != nil {
		return utils.ResponseError(c, 500, "F03", "Error al obtener consultas")
	}
	defer rows.Close()

	var lista []models.Consulta
	for rows.Next() {
		var d models.Consulta
		if err := rows.Scan(&d.ID, &d.PacienteID, &d.MedicoID, &d.ConsultorioID, &d.HorarioID,
			&d.Diagnostico, &d.Costo, &d.Tipo); err != nil {
			return utils.ResponseError(c, 500, "F04", "Error al leer datos")
		}
		lista = append(lista, d)
	}

	return utils.ResponseSuccess(c, 200, "S02", lista)
}

// GetConsultaByID devuelve una consulta por su ID
func GetConsultaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var d models.Consulta

	err := config.DB.QueryRow(context.Background(),
		`SELECT id_consulta, id_paciente, id_medico, id_consultorio, id_horario, diagnostico, costo, tipo
		  FROM consultas WHERE id_consulta=$1`, id,
	).Scan(&d.ID, &d.PacienteID, &d.MedicoID, &d.ConsultorioID, &d.HorarioID, &d.Diagnostico, &d.Costo, &d.Tipo)
	if err != nil {
		return utils.ResponseError(c, 404, "F05", "Consulta no encontrada")
	}

	return utils.ResponseSuccess(c, 200, "S03", []fiber.Map{{
		"id":             d.ID,
		"id_paciente":    d.PacienteID,
		"id_medico":      d.MedicoID,
		"id_consultorio": d.ConsultorioID,
		"id_horario":     d.HorarioID,
		"diagnostico":    d.Diagnostico,
		"costo":          d.Costo,
		"tipo":           d.Tipo,
	}})
}

// UpdateConsulta actualiza una consulta existente
func UpdateConsulta(c *fiber.Ctx) error {
	id := c.Params("id")

	body := c.Body()
	schemaLoader := gojsonschema.NewStringLoader(consultaSchema)
	docLoader := gojsonschema.NewBytesLoader(body)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return utils.ResponseError(c, 400, "SU00", "Error validando JSON")
	}
	if !result.Valid() {
		errs := make([]string, 0, len(result.Errors()))
		for _, desc := range result.Errors() {
			errs = append(errs, desc.String())
		}
		return utils.ResponseError(c, 400, "SU01", "Esquema inválido: "+strings.Join(errs, "; "))
	}

	dto := new(models.Consulta)
	if err := c.BodyParser(dto); err != nil {
		return utils.ResponseError(c, 400, "SU02", "Datos inválidos")
	}

	_, err = config.DB.Exec(context.Background(),
		`UPDATE consultas SET id_paciente=$1, id_medico=$2, id_consultorio=$3, id_horario=$4,
		  diagnostico=$5, costo=$6, tipo=$7 WHERE id_consulta=$8`,
		dto.PacienteID, dto.MedicoID, dto.ConsultorioID, dto.HorarioID,
		dto.Diagnostico, dto.Costo, dto.Tipo, id)
	if err != nil {
		return utils.ResponseError(c, 500, "SU03", "Error al actualizar consulta")
	}

	return utils.ResponseSuccess(c, 200, "S04", []fiber.Map{
		{"message": "Consulta actualizada correctamente"},
	})
}

// DeleteConsulta elimina una consulta por su ID
func DeleteConsulta(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM consultas WHERE id_consulta=$1`, id)
	if err != nil {
		return utils.ResponseError(c, 500, "F08", "Error al eliminar consulta")
	}

	return utils.ResponseSuccess(c, 200, "S05", []fiber.Map{
		{"message": "Consulta eliminada correctamente"},
	})
}
