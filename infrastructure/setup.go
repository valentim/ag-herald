package infrastructure

import (
	"github.com/valentim/ag-herald/infrastructure/database"
)

// Setup will contain all the necessary pre actions before run the code
func Setup() {
	d := database.Database{
		Name: "herald.db",
	}

	d.Setup()
}
