package main

import (
	"database/sql"
	"flag"
	"log"
	"sync"

	"github.com/claytonphenderson/scope/internal/data"
	"github.com/claytonphenderson/scope/internal/ingress"
	"github.com/claytonphenderson/scope/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	var watch = flag.Bool("watch", false, "Enable to have scope run until you Ctrl+C")
	flag.Parse()

	logger, _ := zap.NewProduction()
	ingressCh := make(chan models.DbRecord, 10_000)
	tables := make(map[string]bool)
	var wg sync.WaitGroup
	db, dbErr := sql.Open("sqlite3", "/tmp/scope.db")

	if dbErr != nil {
		log.Fatal("error opening scope db", dbErr)
	}
	defer db.Close()

	go func() {
		for {
			record := <-ingressCh
			data.Write(record, tables, db, logger)
			if !*watch {
				wg.Done()
			}
		}
	}()

	if *watch {
		ingress.ReadFromStdIn(ingressCh, nil)
		select {}
	} else {
		ingress.ReadFromStdIn(ingressCh, &wg)
		wg.Wait()
	}

	logger.Info("Done!")
}
