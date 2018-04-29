package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/facebookgo/httpdown"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/frozzare/statscoll/api"
	"github.com/frozzare/statscoll/cache"
	"github.com/frozzare/statscoll/config"
	"github.com/frozzare/statscoll/stat"
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

	// Create database connection.
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		log.Fatalf("statscoll: %s\n", err)
	}
	defer db.Close()

	// Atuo migrate database.
	db.AutoMigrate(&stat.Stat{})

	// Create the cache.
	cache, err := cache.New()
	if err != nil {
		log.Fatalf("statscoll: %s\n", err)
	}
	defer cache.Close()

	// Create api handler.
	handler, err := api.NewHandler(cache, db)
	if err != nil {
		log.Fatalf("statscoll: %s\n", err)
	}

	hd := &httpdown.HTTP{
		StopTimeout: 10 * time.Second,
		KillTimeout: 1 * time.Second,
	}

	fmt.Printf("Listening on http://0.0.0.0:%d\n", c.Port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: handler,
	}

	if err := httpdown.ListenAndServe(server, hd); err != nil {
		log.Fatal(err)
	}
}
