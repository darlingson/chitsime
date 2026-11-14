package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/darlingson/chitsime/internal/store"
)

func main() {
	dbPath := flag.String("db", "data/chitsime.db", "path to the SQLite database file")
	fresh := flag.Bool("fresh", false, "delete an existing database before creating it")
	flag.Parse()

	if *fresh {
		for _, p := range []string{*dbPath, *dbPath + "-wal", *dbPath + "-shm"} {
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				log.Fatalf("remove %s: %v", p, err)
			}
		}
	}

	st, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := st.Setup(ctx); err != nil {
		log.Fatal(err)
	}

	n, err := st.CountJobs(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("database ready: %s (%d jobs)\n", *dbPath, n)
}
