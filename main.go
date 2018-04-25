package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/TV4/graceful"
	"github.com/frozzare/statscoll/api"
	"github.com/frozzare/statscoll/config"
	"github.com/frozzare/statscoll/db"
	"github.com/spf13/pflag"
)

func main() {
	var (
		configFile string
	)

	pflag.StringVarP(&configFile, "config", "c", "config.yml", "sets the config file")
	pflag.Parse()

	// Try to read the config file.
	c, err := config.ReadFile(configFile)
	if err != nil {
		log.Printf("statscoll: %s\n", err.Error())
		return
	}

	db, err := db.Open("mysql", c.DSN)
	if err != nil {
		log.Fatalf("statscoll: %s\n", err)
	}

	handler, err := api.NewHandler(db)
	if err != nil {
		log.Fatalf("statscoll: %s\n", err)
	}

	graceful.LogListenAndServe(&http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: handler,
	})
}
