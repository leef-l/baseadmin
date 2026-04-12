package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"gbaseadmin/app/upload/internal/ops/backfill"
	"gbaseadmin/internal/cmdutil"
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		configPath string
		batchSize  int
		dryRun     bool
	)

	flag.StringVar(&configPath, "config", "", "GoFrame config file path")
	flag.IntVar(&batchSize, "batch-size", 500, "batch size")
	flag.BoolVar(&dryRun, "dry-run", false, "scan only, do not update dir_id")
	flag.Parse()

	if err := cmdutil.UseConfigFile(configPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	result, err := backfill.Run(context.Background(), backfill.Options{
		BatchSize: batchSize,
		DryRun:    dryRun,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Printf(
		"upload dir backfill done: scanned=%d updated=%d no_relative_dir=%d no_match=%d ambiguous=%d dry_run=%t\n",
		result.Scanned,
		result.Updated,
		result.NoRelativeDir,
		result.NoMatch,
		result.Ambiguous,
		dryRun,
	)
	return 0
}
