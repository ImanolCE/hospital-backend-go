# hospital-backend-go

# Sistema de Citas y Reportes Hospitalarios - Backend en Go + Fiber

Este proyecto, trato de implementar un backend funcional para un sistema de gestión de citas médicas y reportes de un hospital. La aplicación está desarrollada en Go utilizando el framework Fiber, y se conecta a Supabase como base de datos.

---

## Funcionalidades principales

Gestión de usuarios (pacientes, médicos, enfermeros, administradores)  
Sistema de autenticación con JWT  
CRUD de consultorios médicos  
CRUD de horarios disponibles  
CRUD de consultas médicas  
Protección de rutas mediante middleware de autenticación  
Conexión a Supabase utilizando cadena de conexión segura  

---

## Estructura del proyecto

```bash

hospital-back/
├── config/           # Configuración y la conexión a Supabase
├── handlers/         # Lógica de cada endpoint (usuarios, consultorios, etc.)
├── middleware/       # Middleware de autenticación con un JWT
├── models/           # Modelos de los datos
├── routes/           # Definición de las rutas y la agrupación de endpoints
├── utils/            # Funciones auxiliares (generación de tokens)
├── .env              # Variables de entorno (lo que no se sube a Git)
├── go.mod / go.sum   # Dependencias del proyecto
└── main.go           # Archivo principal

```
---

##  Requisitos

- Go versión 1.22 o superior

- Crear/Tner una cuenta en Supabase 

- Fiber framework, instalar: (go get github.com/gofiber/fiber/v2)

- Biblioteca de JWT, instalar: (go get github.com/golang-jwt/jwt/v4)

- Controlador PostgreSQL, instalar (go get github.com/jackc/pgx/v5), el controlador pgx/v5 se obtiene directamente desde Supabase al configurar el acceso a la base de datos con Golang.

---

## Ejecucion para el proyecto 

1.- Clonar el repositorio:

git clone https://github.com/TU_USUARIO/hospital-back.git
cd hospital-back

2.- Crear archivo .env con la siguiente variable:

SUPABASE_URL=postgresql://usuario:contraseña@host:puerto/base_de_datos

Nota: se debera cambiar la contrseña con la contraseña de su proyecto en supabase

3.- Instalar dependencias:

go mod tidy
go get github.com/gofiber/fiber/v2
go get github.com/golang-jwt/jwt/v4

4.- Para ejecutar la aplicación:

go run main.go

---

## En la seguridad 

-Se tiene un JWT, que se utiliza para proteger rutas y gestionar sesiones.
-Los tokens expiran tras un periodo de tiempo.

---

## Endpoints implementados

**Para la Autenticación**

POST /usuarios -> Crear usuario

POST /login -> Login y generación de token

**Usuarios (se necesita token)**

GET /usuarios → Para listar todos los usuarios

GET /usuarios/:id -> Detalle por el ID del usuario

PUT /usuarios/:id -> Para actualizar un usuario

DELETE /usuarios/:id -> Para eliminar un usuario

**Consultorios (se necesita token)**

POST /consultorios -> Para ingresar un consultorio 

GET /consultorios -> Para listar los consultorios 

GET /consultorios/:id -> Para detallar los consultorios por id

PUT /consultorios/:id -> Para actualizar un consultorio 

DELETE /consultorios/:id -> Para eliminar un consultorio

**Horarios (se necesita token)**

POST /horarios -> Para ingresar un horario

GET /horarios -> Para listar los horarios

GET /horarios/:id -> Para detallar los horarios por id

PUT /horarios/:id -> Para actualizar un horario

DELETE /horarios/:id -> Para eliminar un horario

**Consultas (se necesita token)**

POST /consultas -> Para ingresar una consulta 

GET /consultas -> Para lsitar las consultas 

GET /consultas/:id -> Para detallar las consultas por id

PUT /consultas/:id -> pPara actualizar una consulta 

DELETE /consultas/:id -> Para eliminar una consulta 

---

## Notas

Este backend que se relizo, se puede ser extendido con:

- CRUD de recetas
- CRUD de expedientes clínicos
- Reportes avanzados
- Podria implementarse un hash de contraseñas