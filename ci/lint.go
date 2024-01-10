package main

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Run various linters against the source code.
func (m *Ci) Lint() *Lint {
	return &Lint{
		Source: m.Source,
	}
}

type Lint struct {
	// +private
	Source *Directory
}

func (m *Lint) All(ctx context.Context) error {
	var group errgroup.Group

	group.Go(func() error {
		_, err := m.Go().Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	group.Go(func() error {
		_, err := m.Helm().Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}

func (m *Lint) Go() *Container {
	return dag.GolangciLint(GolangciLintOpts{
		Version:   golangciLintVersion,
		GoVersion: goVersion,
	}).
		Run(m.Source, GolangciLintRunOpts{
			Verbose: true,
		})
}

func (m *Lint) Helm() *Container {
	return dag.Helm(HelmOpts{Version: helmVersion}).Lint(helmChartDir(m.Source))
}
