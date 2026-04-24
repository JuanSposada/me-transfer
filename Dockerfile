# Fase de construcción
FROM golang:1.25.0-alpine AS builder

# Instalamos git por si alguna dependencia lo necesita
RUN apk add --no-cache git

WORKDIR /app

# Copiamos los archivos de módulos primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario apuntando a tu main.go
RUN go build -o main ./cmd/api/main.go

# Fase de ejecución final (imagen ligera)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Traemos el ejecutable desde la fase builder
COPY --from=builder /app/main .

# Exponemos el puerto
EXPOSE 8080

# Comando para arrancar la app
CMD ["./main"]