#!/bin/bash

# Esperar a que MySQL esté completamente inicializado
sleep 5

# Conectarse a MySQL y ejecutar los comandos necesarios
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" <<EOF
CREATE DATABASE IF NOT EXISTS MoonTech;
CREATE USER IF NOT EXISTS 'apiuser'@'%' IDENTIFIED BY 'skibidi!';
GRANT ALL PRIVILEGES ON MoonTech.* TO 'apiuser'@'%';
FLUSH PRIVILEGES;
EOF

# Ejecutar el script del esquema si existe
mysql -u apiuser -p'skibidi!' MoonTech < /docker-entrypoint-initdb.d/schema.sql
