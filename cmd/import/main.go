package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/darlingson/chitsime/internal/store"
)

func main() {
	dbPath := flag.String("db", "data/chitsime.db", "path to the SQLite database file")
	file := flag.String("file", "data/jobs-structured.json", "path to the structured jobs JSON file")
	flag.Parse()

	st, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := st.Setup(ctx); err != nil {
		log.Fatal(err)
	}

	n, err := st.ImportFromFile(ctx, *file)
	if err != nil {
		log.Fatal(err)
	}

	total, err := st.CountJobs(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("imported %d jobs from %s into %s (%d total)\n", n, *file, *dbPath, total)
}
