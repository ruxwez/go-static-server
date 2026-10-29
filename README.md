# go-static-server

Imagen Docker mínima para servir archivos estáticos: sitios SPA (React + React Router) o un CDN de imágenes, videos o cualquier asset. Basada en Go + [Fiber v3](https://gofiber.io), sin shell y ejecutándose como usuario no root.

## Uso

Copia tus archivos estáticos en `/app/static`:

```dockerfile
FROM tu-usuario/go-static-server:latest
COPY dist/ /app/static/
```

```bash
docker build -t mi-sitio .
docker run -p 8080:8080 mi-sitio
```

Para un CDN de assets, copia la carpeta de medios igual:

```dockerfile
FROM tu-usuario/go-static-server:latest
COPY media/ /app/static/
```

### Comportamiento

- `/` sirve `index.html`.
- Las rutas sin extensión (`/productos/42`) devuelven `index.html`, así que funcionan las rutas de React Router.
- Un archivo que no existe (`/imagen.png`) devuelve `404`.
- Soporta `Range` requests, necesarios para reproducir y hacer seek en videos y audio.
- Comprime con gzip/brotli los archivos de texto (HTML, CSS, JS, JSON, SVG…). Imágenes y videos no se comprimen.

### Variables de entorno

| Variable | Por defecto | Descripción |
|----------|-------------|-------------|
| `PORT`   | `8080`      | Puerto en el que escucha el servidor dentro del contenedor. |

## Desarrollo

```bash
go test ./...
go run .            # sirve ./static en http://localhost:8080
```

## Publicar la imagen

El workflow `.github/workflows/docker-publish.yml` ejecuta los tests y publica la imagen para `linux/amd64` y `linux/arm64`:

- Push a `main`: tag `main`.
- Push de un tag `v1.2.3`: tags `1.2.3`, `1.2` y `latest`.
- Pull requests: solo compila, no publica.

Configura estos secretos en el repositorio de GitHub (Settings → Secrets and variables → Actions):

- `DOCKERHUB_USERNAME`: tu usuario de Docker Hub.
- `DOCKERHUB_TOKEN`: un access token de Docker Hub.

La imagen se publica como `<DOCKERHUB_USERNAME>/go-static-server`.
