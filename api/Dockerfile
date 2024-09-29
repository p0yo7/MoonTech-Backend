# Usar la imagen oficial de Golang como base
FROM golang:1.23 AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el resto de los archivos del proyecto al contenedor
COPY . .

# Asegurarse de que install_packages.sh tenga el formato correcto para Linux
# Esto evitará problemas con los finales de línea en Windows
RUN apt-get update && apt-get install -y dos2unix \
    && dos2unix /app/install_packages.sh

# Ejecutar el script install_packages.sh
RUN /bin/sh /app/install_packages.sh

# Compilar la aplicación
RUN go build -o myapp .

# Exponer el puerto 8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./myapp"]
