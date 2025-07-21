
package handlers

import (
    "context"
    "unicode"

    "strings"
    //"time"

    "hospital-back/config"
    "hospital-back/models"
    "hospital-back/utils"
    

    "github.com/gofiber/fiber/v2"

    "github.com/xeipuuv/gojsonschema"

    "golang.org/x/crypto/bcrypt"

    "github.com/pquerna/otp"
    "github.com/pquerna/otp/totp"
    
)

const createUserSchema = `{
  "type":"object",
  "required":["nombre","apellido","correo","password","tipo_usuario"],
  "properties":{
    "nombre":{"type":"string"},
    "apellido":{"type":"string"},
    "correo":{"type":"string","format":"email"},
    "password":{"type":"string","minLength":12},
    "tipo_usuario":{"type":"string","enum":["paciente","medico","enfermera","admin"]}
  },
  "additionalProperties":false
}`



// CreateUser crrea un nuevo usuario en la bd
func CreateUser(c *fiber.Ctx) error {

    // 1) Validar JSON Schema
    body := c.Body()
    schemaLoader := gojsonschema.NewStringLoader(createUserSchema)
    docLoader := gojsonschema.NewBytesLoader(body)
    result, err := gojsonschema.Validate(schemaLoader, docLoader)
    if err != nil {
        return utils.ResponseError(c, 400, "F01", "Error validando JSON")
    }
    if !result.Valid() {
        errs := make([]string, 0, len(result.Errors()))
        for _, desc := range result.Errors() {
            errs = append(errs, desc.String())
        }
        return utils.ResponseError(c, 400, "F02", "Esquema inválido: "+strings.Join(errs, "; "))
    }

    // 2) Parsear al struct
    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return utils.ResponseError(c, 400, "F03", "Datos inválidos")
        
    }
    

    // 3)  password
    if len(user.Password) < 12 {
        return utils.ResponseError(c, 400, "F04", "Password requiere ≥12 caracteres")
    }
    var hasNum, hasSym bool
    for _, r := range user.Password {
        switch {
        case '0' <= r && r <= '9':
            hasNum = true
        case strings.ContainsRune("!@#$%^&*()", r):
            hasSym = true
        }
    }
    if !hasNum || !hasSym {
        return utils.ResponseError(c, 400, "F05", "Password requiere número y símbolo")
    }

    // 4) Hashear y guardar el password 
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return utils.ResponseError(c, 500, "F06", "Error al procesar password")
    }
    user.Password = string(hash)

     // 5) Generar secreto TOTP
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "HospitalApp",
        AccountName: user.Correo,
        Algorithm:   otp.AlgorithmSHA1,
        Digits:      otp.DigitsSix,
        Period:      30,
    })
    if err != nil {
        return utils.ResponseError(c, 500, "F07", "Error generando MFA secret")
    }
    user.MFASecret = key.Secret()
    // user.MFAEnabled queda en false por defecto


    // 6) Insertar usuario con MFA en la BD
    _, err = config.DB.Exec(context.Background(),
        `INSERT INTO usuarios 
         (nombre, apellido, correo, password, tipo_usuario, mfa_secret, mfa_enabled)
         VALUES ($1,$2,$3,$4,$5,$6,$7)`,
        user.Nombre, user.Apellido, user.Correo, user.Password, user.Tipo,
        user.MFASecret, false,
    )
    if err != nil {
        return utils.ResponseError(c, 500, "F08", "Error al insertar usuario")
    }

    /* // 5) Responder éxito
    return utils.ResponseSuccess(c, 201, "S01", []fiber.Map{
        {"message": "Usuario creado correctamente"},
    }) */

    // 7) Devolver al cliente la URL OTP para que la escanee
    return utils.ResponseSuccess(c, 201, "S01", []fiber.Map{
        {
            "message":     "Usuario creado correctamente",
            "otpauth_url": key.URL(),  // ej: otpauth://totp/HospitalApp:email?secret=XXX...
        },
    })
     

    
}

// Para la contraseña acepte 12 carcateres con simbolo ynumeros
func esPasswordValida(pass string) bool {
    if len(pass) < 12 {
        return false
    }
    var tieneNumero, tieneSimbolo bool
    for _, c := range pass {
        switch {
        case unicode.IsNumber(c):
            tieneNumero = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            tieneSimbolo = true
        }
    }
    return tieneNumero && tieneSimbolo

}

// GetUsers devuelve la lista de usuarios registrados en la base de datos
func GetUsers(c *fiber.Ctx) error {
        rows, err := config.DB.Query(context.Background(), 
        "SELECT id_usuario, nombre, apellido, correo, password, tipo_usuario FROM usuarios")
        if err != nil {
            return utils.ResponseError(c, 500, "F08", "Error al obtener usuarios")
        }
        defer rows.Close()

        var usuarios []models.User

        for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Correo, &user.Password, &user.Tipo)
        if err != nil {
            return utils.ResponseError(c, 500, "F09", "Error al leer datos")
        }
        usuarios = append(usuarios, user)
    }
        return utils.ResponseSuccess(c, 200, "S03", usuarios)
}

// Obtener usuario por ID
func GetUserByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User

   err := config.DB.QueryRow(context.Background(),
    "SELECT id_usuario, nombre, apellido, correo, password, tipo_usuario FROM usuarios WHERE id_usuario=$1", id).
    Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Correo, &user.Password, &user.Tipo)

    if err != nil {
        return utils.ResponseError(c, 404, "F10", "Usuario no encontrado")
    }
     return utils.ResponseSuccess(c, 200, "S04", user)
}

// Actualizar usuario por ID
func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user := new(models.User)

    if err := c.BodyParser(user); err != nil {
        return utils.ResponseError(c, 400, "F11", "Datos inválidos")
    }

    // Validación de contraseña
    if !esPasswordValida(user.Password) {
        return utils.ResponseError(c, 400, "F12", "Contraseña debe tener mínimo 12 caracteres, incluir símbolo y número")
    }

    // Hasheo de contraseña
    hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
         return utils.ResponseError(c, 500, "F13", "Error al procesar la contraseña")
    }
    user.Password = string(hash)

    // Actualizar en la base de datos
    _, err = config.DB.Exec(context.Background(),
        "UPDATE usuarios SET nombre=$1, apellido=$2, correo=$3, password=$4, tipo_usuario=$5 WHERE id_usuario=$6",
        user.Nombre, user.Apellido, user.Correo, user.Password, user.Tipo, id)

    if err != nil {
        return utils.ResponseError(c, 500, "F14", "Error al actualizar usuario")
    }

    return utils.ResponseSuccess(c, 200, "S06", []fiber.Map{{"message": "Usuario eliminado correctamente"}})
}


// Eliminar usuario por ID
func DeleteUser(c *fiber.Ctx) error {
    id := c.Params("id")

    _, err := config.DB.Exec(context.Background(),
        "DELETE FROM usuarios WHERE id_usuario=$1", id)

    if err != nil {
       return utils.ResponseError(c, 500, "F15", "Error al eliminar usuario")
    }

   return utils.ResponseSuccess(c, 200, "S06", []fiber.Map{{"message": "Usuario eliminado correctamente"}})
}




// Login permite a un usuario autenticarse y devuelve un token JWT
func Login(c *fiber.Ctx) error {
    var datos struct {
        Correo   string `json:"correo"`
        Password string `json:"password"`
        OTP      string `json:"otp"`  
    }

    if err := c.BodyParser(&datos); err != nil {
        return utils.ResponseError(c, 401, "F16", "Credenciales inválidas")
    }

    // 1. Obtener id, hash y rol
    var id int
    var hashedPassword, rolNombre string
    if err := config.DB.QueryRow(context.Background(),
        "SELECT id_usuario, password, tipo_usuario FROM usuarios WHERE correo=$1",
        datos.Correo).Scan(&id, &hashedPassword, &rolNombre); err != nil {
        return utils.ResponseError(c, 401, "F17", "Credenciales incorrectas")
    }

    // 2. Validar hash correcto
    if len(hashedPassword) < 4 || hashedPassword[:4] != "$2a$" {
        return utils.ResponseError(c, 401, "F18", "Usuario inválido, recrea tu cuenta")
    }
    if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(datos.Password)) != nil {
         return utils.ResponseError(c, 401, "F19", "Credenciales incorrectas")
    }

    // 2b) Si ya había MFA configurado, exigir OTP
    var mfaEnabled bool
    var mfaSecret  string
    config.DB.QueryRow(context.Background(),
        "SELECT mfa_enabled, mfa_secret FROM usuarios WHERE id_usuario=$1", id,
    ).Scan(&mfaEnabled, &mfaSecret)

    // Si MFA aún no está activado pero el usuario manda OTP, lo ignoramos pero lo advertimos (opcional)
    // Puedes simplemente documentarlo o hacer un log
    if !mfaEnabled && datos.OTP != "" {
        // fmt.Println("MFA no está activado, OTP ignorado")
    }

    if mfaEnabled {
        if datos.OTP == "" {
            return utils.ResponseError(c, 401, "MFA01", "OTP requerido")
        }
        if !totp.Validate(datos.OTP, mfaSecret) {
            return utils.ResponseError(c, 401, "MFA02", "OTP inválido")
        }
    }

    

    // 3. Cargar permisos desde la BBD 
    rows, err := config.DB.Query(context.Background(),
        `SELECT p.nombre
         FROM permisos p
         JOIN rol_permisos rp ON p.id = rp.permiso_id
         JOIN roles r        ON r.id = rp.rol_id
         WHERE r.nombre = $1`, rolNombre)
    if err != nil {
        return utils.ResponseError(c, 500, "F20", "Error al obtener permisos")
    }
    defer rows.Close()

    var permisos []string
    for rows.Next() {
        var p string
        if rows.Scan(&p) == nil {
            permisos = append(permisos, p)
        }
    }

    // 4. Generar Access Token (ahora con tipo de usuario)
    accessToken, err := utils.GenerarToken(id, permisos, rolNombre)
    if err != nil {
        return utils.ResponseError(c, 500, "F21", "Error al generar AccessToken")
    }

    //  Intentamos eliminar tokens anteriores; si falla, solo lo logueamos y seguimos.
     if _, err := config.DB.Exec(context.Background(),
        `DELETE FROM refresh_tokens WHERE user_id = $1`, id); err != nil {
        return utils.ResponseError(c, 500, "F22", "Error al limpiar tokens anteriores")
    }


    // 5. Generar Refresh Token y tmabien guardarlo en BD
    refreshToken, expiresAt := utils.GenerateRefreshToken()
    if _, err := config.DB.Exec(context.Background(),
        `INSERT INTO refresh_tokens (token, user_id, expires_at)
         VALUES ($1, $2, $3)`,
        refreshToken, id, expiresAt,
    ); err != nil {
         return utils.ResponseError(c, 500, "F23", "Error al generar RefreshToken")
    }

    // 6. Devolver ambos
   return utils.ResponseSuccess(c, 200, "S02", []fiber.Map{
    {
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    },
})
}





