package main

import (
	"os"
	"path/filepath"
	"slices"
)

func root() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(wd, "..")
}

// paths to exclude from all contexts
var excludes = []string{
	".direnv",
	".devenv",
	"ci",
	"deploy/charts/**/charts",
}

func exclude(paths ...string) []string {
	return append(slices.Clone(excludes), paths...)
}

func projectDir() *Directory {
	return dag.Host().Directory(root(), HostDirectoryOpts{
		Exclude: exclude(),
	})
}
