package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/database"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

func main() {
	apply := flag.Bool("apply", false, "commit recalculated readiness (default: rollback/dry-run)")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected positional arguments")
		os.Exit(1)
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := database.Open(ctx, dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "database connection failed")
		os.Exit(1)
	}
	defer db.Close()
	n, err := tasks.NewStore(db).RecalculateAll(ctx, *apply)
	if err != nil {
		fmt.Fprintln(os.Stderr, "recalculation failed; transaction not committed")
		os.Exit(1)
	}
	fmt.Printf("Tasks: %d; applied: %t\n", n, *apply)
}
