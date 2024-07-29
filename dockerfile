FROM golang:1.22.5-alpine AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos del proyecto
COPY . .

# Compila la aplicación
RUN go build -o server .

# Etapa final
FROM alpine:latest

# Crea un directorio para la aplicación
WORKDIR /app

# Copia el binario compilado desde la etapa de construcción
COPY --from=builder /app/server .

# Crea un directorio para los archivos estáticos
RUN mkdir -p /app/static

# Establecer el punto de entrada
CMD ["./server"]