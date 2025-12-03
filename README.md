# Veritas Backend - Sistema de Facturación Multi-tenant

Backend desarrollado en Go con Clean Architecture, utilizando EntGO como ORM y PostgreSQL como base de datos.

## 🚀 Características

- ✅ Arquitectura limpia (Clean Architecture)
- ✅ ORM EntGO para gestión de base de datos
- ✅ Autenticación JWT
- ✅ Multi-tenancy
- ✅ CRUD completo de productos
- ✅ Dashboard con métricas y reportes
- ✅ Carga masiva de productos (CSV)
- ✅ API RESTful

## 📋 Requisitos Previos

- Go 1.18+
- PostgreSQL 12+
- Variables de entorno configuradas

## 🛠️ Instalación

1. Clonar el repositorio:
```bash
git clone <repository-url>
cd Veritasbackend
```

2. Instalar dependencias:
```bash
go mod download
```

3. Configurar variables de entorno:
```bash
# Copiar el archivo de ejemplo (o crear .env manualmente)
cp env.example.txt .env
# O crear .env manualmente con el contenido de env.example.txt
```

Editar `.env` con tus credenciales. **Nota: La base de datos está corriendo en Docker** (ver `dbveritas/docker-compose.yml`):
```env
PORT=8080
GIN_MODE=debug

# Database Configuration (Docker Compose)
# Puerto mapeado: 5434:5432 (puerto externo: 5434)
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=veritas_db
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=24h

CORS_ALLOWED_ORIGINS=http://localhost:3000
```

4. Asegurarse de que la base de datos Docker esté corriendo:
```bash
cd ../dbveritas
docker-compose up -d
```

La base de datos se creará automáticamente al iniciar el contenedor.

5. Generar código de EntGO:
```bash
go generate ./ent
```

6. Ejecutar migraciones (se ejecutan automáticamente al iniciar el servidor)

7. Poblar la base de datos con datos de prueba:
```bash
go run cmd/seed/main.go
```

8. Iniciar el servidor:
```bash
go run cmd/server/main.go
```

El servidor estará disponible en `http://localhost:8080`

## 📁 Estructura del Proyecto

```
Veritasbackend/
├── cmd/
│   ├── server/
│   │   └── main.go              # Punto de entrada del servidor
│   └── seed/
│       └── main.go              # Script de seeding
├── internal/
│   ├── domain/                  # Capa de dominio
│   │   └── repositories/       # Interfaces de repositorios
│   ├── usecase/                 # Casos de uso
│   │   ├── auth/
│   │   ├── dashboard/
│   │   └── stock/
│   ├── handler/                 # Handlers HTTP
│   └── infrastructure/          # Infraestructura
│       ├── config/              # Configuración
│       ├── database/            # Cliente de base de datos
│       ├── middleware/          # Middlewares
│       └── seeder/              # Seeder de datos
├── ent/                         # Código generado por EntGO
│   ├── schema/                  # Schemas de EntGO
│   └── ...
├── pkg/                         # Paquetes compartidos
│   ├── jwt/                     # Utilidades JWT
│   ├── errors/                  # Errores personalizados
│   └── validator/               # Validación
└── go.mod
```

## 🔌 Endpoints de la API

### Autenticación

#### `POST /api/auth/login`
Iniciar sesión.

**Request:**
```json
{
  "email": "admin@demo.veritas.com",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "jwt-token",
  "user": {
    "id": "user-id",
    "email": "admin@demo.veritas.com",
    "name": "Administrador Demo",
    "role": "admin"
  },
  "tenantId": "tenant-id"
}
```

#### `GET /api/auth/me`
Obtener usuario actual.

**Headers:**
- `Authorization: Bearer <token>`
- `X-Tenant-ID: <tenant-id>`

**Response:**
```json
{
  "user": {
    "id": "user-id",
    "email": "admin@demo.veritas.com",
    "name": "Administrador Demo",
    "role": "admin"
  }
}
```

### Dashboard

#### `GET /api/dashboard/metrics`
Obtener métricas del dashboard.

**Headers:**
- `Authorization: Bearer <token>`
- `X-Tenant-ID: <tenant-id>`

**Response:**
```json
{
  "totalProducts": 10,
  "totalInvoices": 5,
  "revenue": 1500.50,
  "lowStockItems": 2
}
```

#### `GET /api/dashboard/reports?period=daily&startDate=2024-01-01&endDate=2024-01-31`
Obtener reportes.

**Query Parameters:**
- `period`: daily, weekly, monthly
- `startDate`: YYYY-MM-DD
- `endDate`: YYYY-MM-DD

**Headers:**
- `Authorization: Bearer <token>`
- `X-Tenant-ID: <tenant-id>`

### Stock

#### `GET /api/stock?page=1&limit=20`
Listar productos.

**Headers:**
- `Authorization: Bearer <token>`
- `X-Tenant-ID: <tenant-id>`

#### `POST /api/stock`
Crear producto.

**Request:**
```json
{
  "name": "Producto Ejemplo",
  "description": "Descripción del producto",
  "price": 99.99,
  "stock": 50,
  "sku": "PROD-001"
}
```

#### `PUT /api/stock/:id`
Actualizar producto.

#### `DELETE /api/stock/:id`
Eliminar producto.

#### `POST /api/stock/upload`
Carga masiva de productos (CSV).

**Request:** `multipart/form-data` con campo `file`

**CSV Format:**
```csv
name,description,price,stock,sku
Producto 1,Descripción 1,10.50,100,SKU-001
Producto 2,Descripción 2,20.75,50,SKU-002
```

## 👥 Usuarios de Prueba

Después de ejecutar el seeder, tendrás los siguientes usuarios:

### Tenant: "Empresa Demo" (slug: demo)
- `admin@demo.veritas.com` / `admin123` (Admin)
- `user@demo.veritas.com` / `user123` (User)

### Tenant: "Acme Corporation" (slug: acme)
- `admin@acme.com` / `admin123` (Admin)
- `manager@acme.com` / `manager123` (Manager)

### Tenant: "Tech Solutions" (slug: tech)
- `admin@techsolutions.com` / `admin123` (Admin)
- `user@techsolutions.com` / `user123` (User)

## 🛠️ Comandos Útiles

```bash
# Generar código de EntGO
go generate ./ent

# Ejecutar seeder
go run cmd/seed/main.go

# Ejecutar servidor
go run cmd/server/main.go

# Compilar
go build -o bin/server cmd/server/main.go

# Ejecutar tests
go test ./...
```

## 🔐 Seguridad

- Passwords hasheados con bcrypt
- Autenticación JWT
- Validación de tenant en cada request
- CORS configurado
- Headers de seguridad

## 📝 Notas

- Las migraciones se ejecutan automáticamente al iniciar el servidor
- El seeder es idempotente (puede ejecutarse múltiples veces)
- Los IDs de tenant se manejan mediante el header `X-Tenant-ID`

## 🚧 Próximas Mejoras

- [ ] Refresh token automático
- [ ] Tests unitarios y de integración
- [ ] Documentación Swagger/OpenAPI
- [ ] Rate limiting
- [ ] Logging estructurado
- [ ] Métricas y monitoreo

## 📄 Licencia

Este proyecto es parte del curso de Sistemas 3.

