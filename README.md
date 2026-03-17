# Tecnology Inventory — Backend

API REST para el sistema de gestión de inventario tecnológico. Construida con Go, Fiber v3 y PostgreSQL.

---

## Stack

| Tecnología | Versión | Uso |
|---|---|---|
| Go | 1.24+ | Lenguaje principal |
| Fiber | v3 | Framework HTTP |
| PostgreSQL | 15 | Base de datos |
| lib/pq | 1.11 | Driver PostgreSQL |
| golang-jwt/jwt | v5 | Autenticación JWT |
| google/uuid | 1.6 | Generación de UUIDs |
| Docker + Compose | — | Contenedorización |

---

## Estructura del proyecto

```
inventory-juanfe/
├── cmd/
│   └── main.go                    # Entry point
├── config/
│   ├── config.go                  # Variables de entorno
│   └── database.go                # Conexión a PostgreSQL
├── dtos/
│   ├── request/                   # DTOs de entrada (bind desde HTTP)
│   └── response/                  # DTOs de salida (JSON al frontend)
├── handlers/                      # Capa HTTP — recibe requests, llama services
│   ├── asset_handler.go
│   ├── assignment_handler.go
│   ├── auth_handler.go
│   ├── catalog_handler.go
│   ├── dashboard_handler.go
│   └── inventory_handler.go
├── middleware/
│   └── auth.go                    # Validación JWT
├── models/                        # Structs que mapean a tablas de DB
│   ├── asset.go
│   ├── assignment.go
│   ├── accounting.go
│   ├── catalog.go
│   ├── inventory.go
│   └── user.go
├── repositories/                  # Queries SQL con database/sql puro
│   ├── asset_repository.go
│   ├── assignment_repository.go
│   ├── accounting_group_repository.go
│   ├── city_repository.go
│   ├── area_repository.go
│   ├── category_repository.go
│   ├── dashboard_repository.go
│   ├── inventory_repository.go
│   ├── status_history_repository.go
│   └── user_repository.go
├── routers/                       # Registro de rutas por módulo
│   ├── asset_router.go
│   ├── assignment_router.go
│   ├── auth_router.go
│   ├── catalog_router.go
│   ├── dashboard_router.go
│   └── inventory_router.go
├── services/                      # Lógica de negocio
│   ├── asset_service.go
│   ├── assignment_service.go
│   ├── auth_service.go
│   ├── dashboard_service.go
│   └── inventory_service.go
├── sql/
│   └── init/
│       ├── 001_schema.sql         # Tablas, enums, índices, triggers
│       └── 002_seed_data.sql      # Datos iniciales (363 activos)
├── utils/
│   ├── hash.go                    # bcrypt
│   ├── jwt.go                     # Generar y parsear tokens
│   ├── response.go                # Helpers de respuesta HTTP
│   └── validators.go              # Validaciones de campos
├── .env                           # Variables de entorno (no commitear)
├── .env.example                   # Plantilla de variables
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Arquitectura por capas

```
HTTP Request
    ↓
Handler       — valida el request, llama al service
    ↓
Service       — lógica de negocio, transacciones, reglas
    ↓
Repository    — queries SQL, acceso a DB
    ↓
PostgreSQL
```

Cada capa solo se comunica con la inmediatamente inferior. Los handlers no tocan la DB. Los repositories no contienen lógica de negocio.

---

## Configuración

### Variables de entorno

Crea un archivo `.env` en la raíz del proyecto:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=inventory_db
JWT_SECRET_KEY=tu_clave_secreta_muy_larga
ALLOWED_ORIGINS=http://localhost:5173
```

### Con Docker (recomendado)

```bash
# Levantar base de datos + API
docker compose up --build

# Solo la base de datos
docker compose up db
```

El Docker Compose aplica automáticamente `001_schema.sql` y `002_seed_data.sql` al iniciar por primera vez.

### Sin Docker

```bash
# Instalar dependencias
go mod download

# Ejecutar
go run cmd/main.go
```

---

## Base de datos

### Tablas principales

| Tabla | Descripción |
|---|---|
| `users` | Usuarios administradores del sistema |
| `assets` | Activos tecnológicos (equipos) |
| `asset_categories` | Categorías: Laptop, Desktop, Tablet, etc. |
| `asset_accounts` | Subcuentas contables |
| `accounting_groups` | Grupos contables padre |
| `cities` | Ciudades: Bogotá, Medellín, Urabá, Cartagena |
| `areas` | Áreas: CES, CIDI, ADMIN, CEO, HILTON, etc. |
| `assignments` | Asignaciones de activos a responsables |
| `status_history` | Historial de cambios de estado de activos |
| `inventory_periods` | Períodos de inventario mensual |
| `inventory_records` | Registros de revisión por período |

### Enums

```sql
logical_status_enum  → active | inactive | written_off
physical_status_enum → optimal | good | fair | deteriorated | out_of_service
assignment_status_enum → active | released | written_off
period_status_enum   → open | closed
```

### Crear usuario administrador

```sql
INSERT INTO users (id, name, email, password_hash, is_active)
VALUES (
    gen_random_uuid(),
    'Administrador',
    'admin@inventario.com',
    '$2a$10$HASH_GENERADO_CON_BCRYPT',
    true
);
```

Para generar el hash de la contraseña:

```go
// cmd/tools/hash/main.go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("tu_password"), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```

```bash
go run cmd/tools/hash/main.go
```

---

## Endpoints

Todos los endpoints protegidos requieren header:

```
Authorization: Bearer <token>
```

### Auth

| Método | Ruta | Descripción | Auth |
|---|---|---|---|
| POST | `/api/v1/auth/login` | Iniciar sesión | No |
| GET | `/api/v1/auth/me` | Usuario actual | Sí |

### Assets

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/assets` | Listar con filtros y paginación |
| POST | `/api/v1/assets` | Crear activo |
| GET | `/api/v1/assets/:id` | Detalle del activo |
| PUT | `/api/v1/assets/:id` | Editar activo |
| PATCH | `/api/v1/assets/:id/status` | Cambiar estado |
| GET | `/api/v1/assets/:id/history` | Historial de estados |
| GET | `/api/v1/assets/:id/assignments` | Asignaciones del activo |

#### Filtros disponibles — GET /assets

```
?city_id=1
?area_id=2
?category_id=3
?asset_account_id=1
?logical_status=active
?physical_status=optimal
?from=2022-01-01
?to=2024-12-31
?search=dell
?page=1
?limit=20
```

### Assignments

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/v1/assignments` | Crear asignación |
| PATCH | `/api/v1/assignments/:id/release` | Liberar asignación |

### Inventory

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/inventory/periods` | Listar períodos |
| POST | `/api/v1/inventory/periods` | Abrir nuevo período |
| PATCH | `/api/v1/inventory/periods/:id/close` | Cerrar período |
| GET | `/api/v1/inventory/periods/:id/assets` | Activos con estado del período |
| GET | `/api/v1/inventory/periods/:id/progress` | Progreso del período |
| POST | `/api/v1/inventory/records` | Confirmar o dar de baja activo |

### Catálogos

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/cities` | Ciudades |
| GET | `/api/v1/areas` | Áreas |
| GET | `/api/v1/categories` | Categorías |
| GET | `/api/v1/accounting-groups` | Grupos contables con subcuentas |
| PATCH | `/api/v1/accounting-groups/:id` | Editar nombre del grupo |

### Dashboard

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/v1/dashboard` | Estadísticas generales |

---

## Convenciones

- Los modelos Go usan `PascalCase` — solo para la capa de DB
- Los DTOs de respuesta definen el JSON con tags `json:"snake_case"` — son la única capa que el frontend consume
- Errores HTTP siguen el formato: `{"success": false, "error": "mensaje"}`
- Respuestas exitosas: `{"success": true, "data": {...}}`
- Todos los IDs son UUID v4

---

## Datos iniciales

El archivo `002_seed_data.sql` carga:

- 4 ciudades (Bogotá, Medellín, Urabá, Cartagena)
- 10 áreas (CES, CIDI, ACADEMICO, PSICOSOCIAL, ADMIN, CEO, GFDO, HILTON, PROYECTO, CONTABILIDAD)
- 11 grupos contables + 11 subcuentas
- 13 categorías de activos
- 363 activos reales migrados desde el inventario original en Excel
