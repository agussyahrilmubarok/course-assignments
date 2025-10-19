#!/bin/bash

# Nama container Docker PostgreSQL Anda
DOCKER_CONTAINER=forum-api-postgres

# Nama database yang ingin Anda buat
DB_NAME=forumapi_test

DB_USER=developer
DB_PASSWORD=supersecretpassword
# Mengeksekusi perintah untuk membuat database
docker exec -it $DOCKER_CONTAINER psql -U $DB_USER -c "CREATE DATABASE $DB_NAME;"

echo "Database $DB_NAME telah berhasil dibuat."