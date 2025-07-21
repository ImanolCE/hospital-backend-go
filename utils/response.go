package utils

import (
  "log"
  "path/filepath"
  "strings"

  "github.com/gofiber/fiber/v2"
  "github.com/xeipuuv/gojsonschema"
)

// buildFileURI convierte una ruta absoluta en un URI compatible con file:///
func buildFileURI(relPath string) string {
  abs, err := filepath.Abs(relPath)
  if err != nil {
    log.Fatalf("Error al resolver ruta absoluta de %s: %v", relPath, err)
  }
  // En Windows abs="C:\..."; convertimos barras y anteponemos file:///
  return "file:///" + strings.ReplaceAll(abs, "\\", "/")
}

// ResponseSuccess valida con JSON-Schema antes de devolver un 2xx
func ResponseSuccess(c *fiber.Ctx, statusCode int, intCode string, data interface{}) error {
  resp := fiber.Map{
    "statusCode": statusCode,
    "intCode":    intCode,
    "data":       data,
  }

  // Carga el URI del schema de respuesta
  schemaURI := buildFileURI("./schemas/response.json")
  //log.Println("📄 Usando response schema URI:", schemaURI)

  schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
  docLoader := gojsonschema.NewGoLoader(resp)

  // Valida
  result, err := gojsonschema.Validate(schemaLoader, docLoader)
  if err != nil {
    log.Printf(" Schema load/parse error: %v", err)
    return c.Status(fiber.StatusInternalServerError).
      JSON(fiber.Map{"error": "Schema validation error", "details": err.Error()})
  }
  if !result.Valid() {
    log.Printf(" Schema validation failed: %v", result.Errors())
    return c.Status(fiber.StatusInternalServerError).
      JSON(fiber.Map{"error": "Invalid response schema", "details": result.Errors()})
  }

  // Envío final
  return c.Status(statusCode).JSON(resp)
}

// ResponseError valida con JSON-Schema antes de devolver un error
func ResponseError(c *fiber.Ctx, statusCode int, intCode string, errMsg string) error {
  resp := fiber.Map{
    "statusCode": statusCode,
    "intCode":    intCode,
    "error":      errMsg,
  }

  // Carga el URI del schema de error
  schemaURI := buildFileURI("./schemas/error_response.json")
  //log.Println("📄 Usando error schema URI:", schemaURI)

  schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
  docLoader := gojsonschema.NewGoLoader(resp)

  result, err := gojsonschema.Validate(schemaLoader, docLoader)
  if err != nil {
    log.Printf(" Error schema load/parse: %v", err)
    return c.Status(fiber.StatusInternalServerError).
      JSON(fiber.Map{"error": "Schema validation error", "details": err.Error()})
  }
  if !result.Valid() {
    log.Printf(" Error response schema validation failed: %v", result.Errors())
    return c.Status(fiber.StatusInternalServerError).
      JSON(fiber.Map{"error": "Invalid error response schema", "details": result.Errors()})
  }

  return c.Status(statusCode).JSON(resp)
}
