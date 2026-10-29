# Compila el servidor
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

# Imagen final: sin shell, sin root
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server

# Sirve los archivos que copies en /app/static
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/server"]
