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

# Generar código de EntGO
RUN go generate ./ent

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

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

