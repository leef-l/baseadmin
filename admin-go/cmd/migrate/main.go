package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"gbaseadmin/internal/dbmigrate"
)

func main() {
	os.Exit(run())
}

func run() int {
	opts := dbmigrate.Options{}
	flag.StringVar(&opts.Dir, "dir", "database/migrations", "migration directory")
	flag.StringVar(&opts.DSN, "dsn", "", "mysql dsn, highest priority")
	flag.StringVar(&opts.Link, "link", "", "GoFrame database link, for example mysql:user:pass@tcp(host:3306)/db")
	flag.StringVar(&opts.ConfigPath, "config", "", "GoFrame config file path")
	flag.StringVar(&opts.DatabaseNode, "database-node", "default", "database node name inside config file")
	flag.IntVar(&opts.Steps, "steps", 1, "steps for down action")
	flag.IntVar(&opts.Version, "version", -1, "version for force action")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <up|down|version|force|create> [name]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return 2
	}
	opts.Action = args[0]
	if opts.Action == "create" && len(args) > 1 {
		opts.Name = args[1]
	}

	if err := dbmigrate.Run(context.Background(), opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
