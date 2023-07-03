# Lexica API

Backend for Lexica. Built in Go, using Chi and Postgres.

## Prerequisites

Make sure you have Go installed, at least version 1.20.

We recommend you install several tools to help manage the code and to ease your development experience.

1. Air to utilize hot reloads (Optional)

```bash
go install github.com/cosmtrek/air@latest
```

2. Migrate CLI to handle database migrations

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

3. Formattag to align Go struct tags (Optional)

```bash
go install github.com/momaek/formattag@latest
```

4. Swag to generate and configure documentation based on the OpenAPI spec

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Alternatively, you can run the `tools/install.sh` script to do this for you.

```bash
bash ./tools/install.sh
```

The build system is managed using Make as seen on the `Makefile` at the root of this repository, however you can run all the necessary commands by yourself. We use PostgreSQL as the database, specifically versions >= 15.0.

To help manage the infrastructure dependencies, we will use Docker to spin up containers quickly as seen in the `docker-compose.yml` file.

To sum it up:

- Go >= 1.20
- PostgreSQL >= 15.0
- Migrate CLI
- Swag CLI
- Docker (highly recommended)
- Air (optional, but recommended)
- Formattag (optional, but recommended)

## Database Migrations

As noted, we will be using the Migrate CLI with Make to create the migrations. We created a `migrate.go` file inside of `db/migrations` to handle the migrating process. You can also use the targets in `Makefile` which handles a bit of the flags for the `migrate.go` file for you.

```bash
# Create a migration
make migration [name_of_migration]

# Apply all migrations
make migrate

# Apply n migrations, which correlates to number of up migrations to apply
make migrate -- -steps 5

# Revert all migrations
make migrate-down

# Revert n migrations
make migrate-down -- -steps 5

# Force migration to the nth version. Should only be used to fix dirty migrations
make migrate-fix 5
```

## Installation

1. Clone this repository
2. Run `go mod tidy` or `make install` to install the dependencies.
3. Create a `.env` file based off the `.env.example` and populate it with your desired values.

## Running

We will assume that you'll be using Docker to start up several dependencies.

1. Run `docker compose up -d` to start the dependency containers (e.g. local database).
2. Run `air` or `make dev` to run the server with hot reloads. Alternatively, you can use `make run` or `go run main.go` too if you don't need hot reload.
3. For production, we will use `make build` to create the executable. Run `make run-build` to run the executable. You can see how to do this without Make by looking at the related targets in `Makefile`.
