package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/cverish/go-jwt-rest-api/internal/models"
)

// Main defines the models to migrate, loads the models with Atlas,
// and writes the migration to `internal/database/migrations`
func main() {
	models := []interface{}{
		&models.User{},
		&models.InvitedUser{},
	}

	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
