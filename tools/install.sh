#!/usr/bin/env sh

echo "Installing tools..."
go install github.com/cosmtrek/air@latest
go install github.com/momaek/formattag@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/swaggo/swag/cmd/swag@latest
