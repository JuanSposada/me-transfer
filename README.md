# Me-Transfer 🚀 | Sistema de Transferencia de Archivos Seguros

Me-Transfer es una aplicación web inspirada en WeTransfer que permite a los usuarios subir archivos y generar enlaces de descarga directa de forma instantánea. El sistema utiliza una arquitectura moderna diseñada para la velocidad y la seguridad.

## 🛠 Tecnologías Utilizadas
Lenguaje: Go (Golang) 1.21+

Framework Web: Gin Gonic (Alto rendimiento)

Base de Datos: PostgreSQL (Alojado en Supabase)

Almacenamiento en la Nube: Supabase Storage (Buckets)

Frontend: HTML5, JavaScript Vanilla y Tailwind CSS para el diseño

Despliegue: Railway (PaaS)

## 🌟 Características Principales
Subida Optimizada: Carga de archivos eficiente hacia el Storage de Supabase.

Enlaces Firmados (Signed URLs): Se generan URLs temporales y seguras que protegen la privacidad de los archivos.

Base de Datos Relacional: Cada archivo cuenta con un UUID único y metadatos registrados en Postgres.

Descarga Directa: El sistema entrega un enlace final que dispara la descarga automáticamente sin necesidad de pasos intermedios.

Arquitectura de un solo paso: El backend procesa la subida, el registro y la generación del link en una sola transacción para mejorar la experiencia de usuario.



## 🔄 Flujo de Funcionamiento
- Selección: El usuario elige un archivo (máximo 5MB).

- Procesamiento: El cliente envía el archivo al endpoint /upload.

- Backend Go:

- Sube el archivo al Bucket de Supabase.

- Genera un link de descarga firmado (Signed URL).

- Guarda los metadatos (nombre, tamaño, tipo) en Postgres.

- Respuesta: El servidor devuelve el link de descarga directa en el JSON de éxito.

- Finalización: El usuario copia el enlace y puede compartirlo. El receptor descarga el archivo al instante al abrir el link.

![alt text](image.png)
### Justificación de la Arquitectura
Usaremos PostgreSQL como fuente de verdad rápida para consultas frecuentes (estado de tokens, expiración) debido a su baja latencia. Supabase se justifica como:

Storage Engine: Para no saturar el servidor local con archivos binarios.

Alta Disponibilidad: Si la DB local falla, la metadata en Supabase sirve como respaldo de desastre para reconstruir el estado del sistema.

### Estructura Base del Proyecto (Go)
Esta estructura sigue los estándares de la comunidad (golang-standards/project-layout) para que sea escalable y fácil de navegar en el repositorio.

```plaintext
go-transfer/
├── cmd/
│   └── api/
│       └── main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── api/
│   │   ├── handlers/        # Lógica de los endpoints (Persona C)
│   │   ├── middleware/      # Seguridad y validaciones (Persona C)
│   │   └── routes.go        # Definición de rutas
│   ├── config/
│   │   └── config.go        # Carga de variables de entorno (.env)
│   ├── models/
│   │   └── file.go          # Estructuras de datos (Persona A)
│   ├── repository/
│   │   ├── postgres/        # Implementación DB Principal (Persona A)
│   │   └── supabase/        # Cliente Supabase y Storage (Persona B)
│   ├── service/
│   │   └── file_service.go  # Lógica de negocio (une repo y cloud)
│   └── worker/
│       └── expiration.go    # Proceso de limpieza automática (Persona A)
├── pkg/
│   └── utils/               # Helpers (Hash, UUID, validadores)
├── .env                     # Credenciales (Postgres_URL, Supabase_Key)
├── go.mod
└── README.md
```

### División de Tareas Equipo de 3

---

### Persona A Backend y Persistencia (PostgreSQL & Core)

**Focus:** El corazón de los datos y la lógica de negocio.

#### Responsabilidades
- **Diseño de base de datos:** Crear el esquema en PostgreSQL.  
  - **Tablas mínimas:** `files`, `tokens`.
- **Capa de Datos:** Implementar el repositorio para PostgreSQL usando **sqlx** o **pgx**.
- **Lógica de Expiración:** Desarrollar un Worker/Job interno en Go (usar `time.Ticker`) que:
  - Busque archivos expirados.
  - Marque registros para eliminación.
- **Seguridad:** Implementar generación de tokens seguros (**UUID v4**) y validación de integridad.

---

### Persona B Especialista en Cloud e Integración (Supabase)

**Focus:** Almacenamiento físico y sincronización externa.

#### Responsabilidades
- **Supabase Storage:** Configurar el bucket y desarrollar el cliente en Go para subir/descargar archivos.
- **Sincronización de Metadata:** Al subir un archivo, guardar un backup de la metadata en Supabase según el requerimiento.
- **Validación de Archivos:** Implementar filtros de seguridad:
  - Validar **MIME types**.
  - Detectar archivos ejecutables.
  - Aplicar límite de tamaño.
- **Logs:** Configurar sistema de logging que persista en Supabase.

---

### Persona C Desarrollador de API y Routing (Transporte)

**Focus:** La cara pública del servicio y la orquestación.

#### Responsabilidades
- **Servidor Web:** Configurar el router (ej. **Gin**, **Echo** o **Chi**).
- **Endpoints:** Implementar handlers:
  - `POST /upload`
  - `GET /download/{token}`
  - `GET /file/{token}`
- **Middleware:** Crear middlewares para:
  - Manejo de errores.
  - Recuperación de pánicos.
  - Límites de tamaño de request.
- **Orquestación:** Integrar y orquestar los servicios provistos por **Persona A** y **Persona B** dentro de los handlers.

---

### Notas Operativas Rápidas
- **Comunicación:** Mantener sincronía diaria breve entre A, B y C para evitar solapamientos.
- **Pruebas:** Definir pruebas end-to-end que cubran flujo de upload → metadata backup → storage → expiración → eliminación.
- **Seguridad:** Revisar permisos de buckets y accesos a la base de datos antes de desplegar a producción.



### Prompt para importar en cuaolquier Modelo de IA para continuar con el mismo hilo:

```plaintext
📋 Resumen del Proyecto: Go-Transfer (Clon WeTransfer)
Objetivo: Desarrollar una app web de transferencia de archivos en Go.
Arquitectura:

DB Principal: PostgreSQL (Información de archivos, tokens, expiración).

DB Secundaria/Cloud: Supabase (Backup de metadatos, Logs y Supabase Storage para archivos físicos).

Justificación: PostgreSQL para latencia mínima en metadata activa; Supabase para escalabilidad de archivos y respaldo de desastre.

División de Tareas (Equipo de 3):

Persona A (Backend & Core): Esquema Postgres, lógica de repositorio, worker de expiración automática y tokens seguros.

Persona B (Cloud & Storage): Integración con Supabase Storage, validación de archivos (MIME, size), logging en la nube.

Persona C (API & Routing): Servidor HTTP (Gin/Echo), middlewares de seguridad, endpoints POST /upload, GET /download/{token}, GET /file/{token}.

Estructura del Proyecto:

go-transfer/
├── cmd/api/main.go          # Punto de entrada
├── internal/
│   ├── api/handlers/        # Persona C
│   ├── models/              # Persona A
│   ├── repository/          # Persona A (PG) y Persona B (Supabase)
│   ├── service/             # Lógica de negocio
│   └── worker/              # Limpieza (Persona A)
├── pkg/utils/               # Helpers
└── .env                     # Configuración


Seguridad Obligatoria: Validación MIME, protección Path Traversal, bloqueo de ejecutables, límite de tamaño por request, UUIDs para tokens.

```


# 🚀 Me-Transfer Backend

Sistema de transferencia de archivos con almacenamiento en Supabase y persistencia en PostgreSQL.

---

## 🛠️ Estado del Proyecto: Hito 1 & 2 (Persona A) (Sebastian)

La **Infraestructura Base** y la **Capa de Datos** ya están configuradas y listas para su integración con los servicios de almacenamiento (Persona B) y la API (Persona C).

### ✅ Logros Actuales:
* **Contenedorización**: PostgreSQL 15 configurado mediante Docker Compose.
* **Esquema de Datos**: Tablas `files` y `tokens` automatizadas con soporte para UUID.
* **Modelos de Dominio**: Structs en Go (`FileMetadata`, `Token`) que actúan como contrato de datos.
* **Capa de Repositorio**: Implementación del patrón *Repository* con `pgxpool` para una gestión eficiente de conexiones.

---

## 📂 Estructura de la Capa de Datos

```text
me-transfer/
├── internal/
│   ├── models/           # Definición de "Objetos" (Structs)
│   │   └── file.go
│   └── repository/       # Interfaces (Contratos) e Implementación
│       ├── repository.go # <-- Persona C: Usar esta interfaz
│       └── postgres/     # Lógica SQL específica
│           ├── db.go
│           └── file_pg.go
├── migrations/           # Scripts de inicialización SQL
└── docker-compose.yml    # Orquestación de servicios
```

## 🚀 Guía de Inicio Rápido
1. **Levantar Base de Datos** 
Asegúrate de tener Docker instalado y ejecuta:

```bash
docker compose up -d

```
2. **Configuración del Entorno (.env)**
Crea un archivo .env basado en el siguiente ejemplo para que el backend pueda conectar con el contenedor:

```env
# Database Config
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin_me_transfer
DB_PASSWORD=admin_pass
DB_NAME=me_transfer_db

# Connection String para el Driver pgx
POSTGRES_URL=postgres://admin_me_transfer:admin_pass@localhost:5432/me_transfer_db?sslmode=disable
```

3. **Verificar Tablas**
Puedes comprobar que el esquema se cargó correctamente ejecutando:
```bash
docker exec -it me_transfer_db psql -U admin_me_transfer -d me_transfer_db -c "\dt"
```

## 📋 Roadmap para el Equipo
### ☁️ Persona B (Supabase & Storage)
- Implementar el cliente de Supabase en internal/storage.

- Objetivo: Crear la función de subida que reciba el buffer del archivo y devuelva el supabase_path.

### 🌐 Persona C (API & Handlers)
- Configurar el servidor (Gin/Echo) y los endpoints.

- Objetivo: Implementar POST /upload.

- Flujo: Recibir MultiPart File -> Llamar a Persona B (Upload) -> Llamar a Persona A (CreateFile).

### 🛠️ Comandos de Desarrollo
- Instalar dependencias: `go mod tidy`

- Correr en modo desarrollo: `go run cmd/api/main.go` (una vez creado)



# Actualizacion 4/8/2026
## 📂 Estructura de Servicios (Actualizada)

Hemos separado la persistencia de datos del almacenamiento de archivos físicos siguiendo principios de Clean Architecture:

```text
internal/
├── repository/       # Capa de DATOS (PostgreSQL) - Metadata y Tokens
└── storage/          # Capa de ARCHIVOS (Supabase) - Blobs/Binarios
    ├── storage.go    # Interfaz/Contrato de Almacenamiento
    └── supabase.go   # Implementación del cliente de Supabase

```

## ☁️ Guía para Persona B (Storage Service)
Se ha dejado preparado el cliente de Supabase para que puedas empezar directamente con la lógica de negocio.

1. Variables de Entorno Necesarias

Asegúrate de configurar estas llaves en tu .env:

Fragmento de código
SUPABASE_URL=[https://tu-proyecto.supabase.co](https://tu-proyecto.supabase.co)
SUPABASE_KEY=tu-service-role-key-secreta
SUPABASE_BUCKET=me-transfer-files

2. Implementación del Cliente
El archivo internal/storage/supabase.go ya cuenta con la inicialización del SDK oficial:

SDK utilizado: github.com/supabase-community/storage-go

Misión: Implementar los métodos Upload y GetSignedURL utilizando el cliente s.client.

3. El Contrato (Interfaz)
Cualquier método que desarrolles debe cumplir con la interfaz StorageService en internal/storage/storage.go. Esto permitirá que la Persona C (API) llame a tus funciones de la siguiente manera:

```Go
// Ejemplo de uso para Persona C
path, err := storageService.Upload(ctx, "archivo.pdf", reader)
```

# 🛠️ Integración del Sistema (Main)
El archivo `cmd/api/main.go` ya realiza el "Wire-up" inicial:

Carga el `.env.`

Inicializa el Pool de PostgreSQL.

(Pendiente) Inicializar el StorageService para pasarlo a los Handlers.


### ¿Por qué agregamos esto?
1. **Claridad de límites**: Le dejas claro a la Persona B que su terreno es `internal/storage`.
2. **Contrato definido**: Ella ya sabe que su función `Upload` debe devolver un `string` (el path) y un `error`.
3. **Dependencias**: Ya sabe que debe usar el paquete de `storage-go`.

---

## 🏗️ Estado de la Infraestructura (Abril  8, 2026)

La fase de cimentación ha sido completada exitosamente. El proyecto cuenta con una arquitectura desacoplada basada en interfaces.

### ⚙️ Servicios Configurados
1.  **Motor de Base de Datos**: PostgreSQL corriendo sobre Docker (Puerto 5432).
2.  **Persistencia**: Capa de repositorio implementada en `internal/repository/postgres`. Soporta operaciones de creación y consulta de metadatos de archivos.
3.  **Storage Service**: Integración con Supabase Storage lista. El cliente se inicializa globalmente y está disponible para la lógica de negocio.

### 🛠️ Roadmap de Integración
- [ ] **Persona B (Storage Specialist)**: Implementar métodos `Upload` y `GetSignedURL` en `internal/storage/supabase.go`.
- [ ] **Persona C (API Developer)**: Crear servidor HTTP y endpoints para recibir archivos (Multipart Form) y coordinar los servicios de Storage y Repository.

### 🚀 Cómo arrancar
1. Levantar DB: `docker compose up -d`
2. Configurar `.env` con las llaves de Supabase.
3. Ejecutar: `go run cmd/api/main.go`