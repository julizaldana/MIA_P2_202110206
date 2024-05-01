FROM golang:alpine

# Instalar graphviz
RUN apk add --no-cache graphviz

# Establecer el directorio de trabajo
WORKDIR /Go

# Copiar los archivos de la aplicación
COPY . .

# Ejecutar go mod download para descargar las dependencias del módulo
RUN go mod download

# Exponer el puerto en el que se ejecuta la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["go", "run", "main.go"]
