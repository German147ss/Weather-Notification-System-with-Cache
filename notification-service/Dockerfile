# Dockerfile para Users Service
FROM golang:1.22 AS build

WORKDIR /app

# Copiar los archivos del módulo Go y las dependencias vendorizadas
COPY go.mod go.sum ./

# Copiar el código fuente
COPY . .

# Compilar el servicio usando las dependencias de vendor/
RUN go build -mod=vendor -o notification-service .

# Exponer el puerto
EXPOSE 8082

# Comando para ejecutar el servicio
CMD ["./notification-service"]
