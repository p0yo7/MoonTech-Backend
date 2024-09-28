#!/bin/bash
# Funcion para instalar los paquetes necesarios
# Leer cada línea del archivo packages.txt
while IFS= read -r package
do
    echo "Instalando: $package"
    go get "$package"
    if [ $? -eq 0 ]; then
        echo "Instalado con éxito: $package"
    else
        echo "Error al instalar el paquete: $package"
    fi
done < packages.txt
