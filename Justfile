# Justfile

# Load environment variables
set dotenv-load

default:
    just --list

run:
    go run cmd/ronbunkun/main.go

build:
    go build -o bin/ronbunkun cmd/ronbunkun/main.go

