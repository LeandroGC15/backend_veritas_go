# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependencias necesarias para compilar
RUN apk add --no-cache git

# Copiar go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Crear main.go en la raíz si no existe
# Primero intentar copiar desde cmd/server/main.go, si no existe, copiar desde el main.go que debería estar en el repo
RUN if [ ! -f "./main.go" ]; then \
        if [ -f "./cmd/server/main.go" ]; then \
            echo "Copiando cmd/server/main.go a ./main.go"; \
            cp ./cmd/server/main.go ./main.go; \
        else \
            echo "ERROR: No se encontró main.go ni cmd/server/main.go"; \
            echo "Por favor, asegúrate de que main.go esté en la raíz del proyecto o que cmd/server/main.go exista"; \
            exit 1; \
        fi; \
    else \
        echo "main.go ya existe en la raíz"; \
    fi && \
    echo "Verificando main.go:" && \
    ls -la ./main.go

# Generar código de EntGO (si existe)
RUN if [ -d "./ent" ]; then go generate ./ent || echo "Warning: EntGO generation skipped"; fi

# Compilar la aplicación
RUN echo "Compilando desde ./main.go" && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copiar el binario desde el builder
COPY --from=builder /app/server .

# Exponer el puerto
EXPOSE 8080

# Ejecutar el servidor
CMD ["./server"]
