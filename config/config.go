package config

import (
	"log"
	"time"
)

const Path = "./Files"
const Layout = "02/Jan/06 15:04"

var Location *time.Location

func init() {
	var err error
	Location, err = time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatalf("Error al cargar la zona horaria: %v", err)
	}
}