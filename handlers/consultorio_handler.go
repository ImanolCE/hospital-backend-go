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

const consultorioSchema = `{
  "type": "object",
  "required": ["nombre", "tipo", "ubicacion", "id_medico"],
  "properties": {
    "nombre":     { "type": "string" },
    "tipo":       { "type": "string" },
    "ubicacion":  { "type": "string" },
    "id_medico":  { "type": "integer" }
  },
  "additionalProperties": false
}`

// CreateConsultorio inserta un nuevo consultorio en la base de datos
func CreateConsultorio(c *fiber.Ctx) error {
	body := c.Body()
	schemaLoader := gojsonschema.NewStringLoader(consultorioSchema)
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

	dto := new(models.Consultorio)
	if err := c.BodyParser(dto); err != nil {
		return utils.ResponseError(c, 400, "SC02", "Datos inválidos")
	}

	_, err = config.DB.Exec(context.Background(),
		`INSERT INTO consultorios (nombre, tipo, ubicacion, id_medico) VALUES ($1, $2, $3, $4)`,
		dto.Nombre, dto.Tipo, dto.Ubicacion, dto.MedicoID)
	if err != nil {
		return utils.ResponseError(c, 500, "SC03", "Error al crear consultorio")
	}

	return utils.ResponseSuccess(c, 201, "S01", []fiber.Map{
		{"message": "Consultorio creado correctamente"},
	})
}


// GetConsultorios devuelve todos los consultorios
func GetConsultorios(c *fiber.Ctx) error {
	rows, err := config.DB.Query(context.Background(),
		`SELECT id_consultorio, nombre, tipo, ubicacion, id_medico FROM consultorios`)
	if err != nil {
		return utils.ResponseError(c, 500, "F03", "Error al obtener consultorios")
	}
	defer rows.Close()

	var lista []models.Consultorio
	for rows.Next() {
		var d models.Consultorio
		if err := rows.Scan(&d.ID, &d.Nombre, &d.Tipo, &d.Ubicacion, &d.MedicoID); err != nil {
			return utils.ResponseError(c, 500, "F04", "Error al leer datos")
		}
		lista = append(lista, d)
	}

	return utils.ResponseSuccess(c, 200, "S02", lista)
}

// GetConsultorioByID devuelve un consultorio por su ID
func GetConsultorioByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var d models.Consultorio

	err := config.DB.QueryRow(context.Background(),
		`SELECT id_consultorio, nombre, tipo, ubicacion, id_medico FROM consultorios WHERE id_consultorio=$1`, id,
	).Scan(&d.ID, &d.Nombre, &d.Tipo, &d.Ubicacion, &d.MedicoID)
	if err != nil {
		return utils.ResponseError(c, 404, "F05", "Consultorio no encontrado")
	}

	return utils.ResponseSuccess(c, 200, "S03", []fiber.Map{{
		"id":         d.ID,
		"nombre":     d.Nombre,
		"tipo":       d.Tipo,
		"ubicacion":  d.Ubicacion,
		"id_medico":  d.MedicoID,
	}})
}

// UpdateConsultorio modifica un consultorio existente
func UpdateConsultorio(c *fiber.Ctx) error {
	id := c.Params("id")

	body := c.Body()
	schemaLoader := gojsonschema.NewStringLoader(consultorioSchema)
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

	dto := new(models.Consultorio)
	if err := c.BodyParser(dto); err != nil {
		return utils.ResponseError(c, 400, "F06", "Datos inválidos")
	}

	_, err = config.DB.Exec(context.Background(),
		`UPDATE consultorios SET nombre=$1, tipo=$2, ubicacion=$3, id_medico=$4 WHERE id_consultorio=$5`,
		dto.Nombre, dto.Tipo, dto.Ubicacion, dto.MedicoID, id)
	if err != nil {
		return utils.ResponseError(c, 500, "F07", "Error al actualizar consultorio")
	}

	return utils.ResponseSuccess(c, 200, "S04", []fiber.Map{
		{"message": "Consultorio actualizado correctamente"},
	})
}


// DeleteConsultorio elimina un consultorio por su ID
func DeleteConsultorio(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM consultorios WHERE id_consultorio=$1`, id)
	if err != nil {
		return utils.ResponseError(c, 500, "F08", "Error al eliminar consultorio")
	}

	return utils.ResponseSuccess(c, 200, "S05", []fiber.Map{
		{"message": "Consultorio eliminado correctamente"},
	})
}
