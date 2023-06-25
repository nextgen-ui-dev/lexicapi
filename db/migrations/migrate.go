package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lexica-app/lexicapi/app"
)

var action string
var steps uint
var version uint

func init() {
	flag.StringVar(&action, "action", "up", "run db migrations [up | down]")
	flag.UintVar(&steps, "steps", 0, "amount of migrations run. If not specified, run all")
	flag.UintVar(&version, "version", 0, "version of migration to force to")
	flag.Parse()
}

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	migrate, err := migrate.New("file://db/migrations", config.DbUrl)
	if err != nil {
		log.Fatal("Failed to read migration files:", err)
	}

	if action == "up" {
		if steps != 0 {
			if err = migrate.Steps(int(steps)); err != nil {
				if err.Error() == "no change" {
					log.Println("Nothing to run")
				} else {
					log.Fatal("Failed to run migration:", err.Error())
				}
			}
		} else {
			if err = migrate.Up(); err != nil {
				if err.Error() == "no change" {
					log.Println("Nothing to run")
				} else {
					log.Fatal("Failed to run migration:", err.Error())
				}
			}
		}
	} else if action == "down" {
		if steps != 0 {
			if err = migrate.Steps(-1 * int(steps)); err != nil {
				if err.Error() == "no change" {
					log.Println("Nothing to run")
				} else {
					log.Fatal("Failed to run migration:", err.Error())
				}
			}
		} else {
			if err = migrate.Down(); err != nil {
				if err.Error() == "no change" {
					log.Println("Nothing to run")
				} else {
					log.Fatal("Failed to run migration:", err.Error())
				}
			}
		}
	} else if action == "force" {
		if version != 0 {
			if err = migrate.Force(int(version)); err != nil {
				if err.Error() == "no change" {
					log.Println("Nothing to run")
				} else {
					log.Fatal("Failed to fix/force migration:", err)
				}
			}
		} else {
			log.Fatal("Invalid migration version target")
		}
	} else {
		log.Fatal("Invalid migration action")
	}
}
