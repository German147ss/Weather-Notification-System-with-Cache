FROM golang:1.22 as builder

WORKDIR /app

# Copiar los archivos del módulo Go y las dependencias vendorizadas
COPY go.mod go.sum ./

# Copiar el código fuente
COPY . .

# Compilar el servicio
RUN go build -mod=vendor -o weather-service .

# Exponer el puerto en el que el servicio escucha
EXPOSE 8083

# Comando para ejecutar el servicio
CMD ["./weather-service"]
