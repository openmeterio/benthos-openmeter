package main

import (
	"context"

	"golang.org/x/sync/errgroup"
)

const (
	goVersion           = "1.21.5"
	golangciLintVersion = "v1.54.2"

	alpineBaseImage = "alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48"

	helmDocsVersion = "v1.11.3"
	helmVersion     = "3.13.2"

	ghVersion = "2.42.1"
)

const imageRepo = "ghcr.io/openmeterio/benthos-openmeter"

type Ci struct {
	// +private
	RegistryUser string

	// +private
	RegistryPassword *Secret

	// +private
	GitHubToken *Secret

	// Project source directory
	// This will become useful once pulling from remote becomes available
	//
	// +private
	Source *Directory
}

func New(
	// Checkout the repository (at the designated ref) and use it as the source directory instead of the local one.
	// +optional
	checkout string,

	// Container registry user (required for pushing images).
	// +optional
	registryUser string,

	// Container registry password (required for pushing images).
	// +optional
	registryPassword *Secret,

	// GitHub token.
	// +optional
	githubToken *Secret,
) (*Ci, error) {
	var source *Directory

	if checkout != "" {
		source = dag.Git("https://github.com/openmeterio/benthos-openmeter.git", GitOpts{
			KeepGitDir: true,
		}).Branch(checkout).Tree()
	} else {
		source = projectDir()
	}

	return &Ci{
		RegistryUser:     registryUser,
		RegistryPassword: registryPassword,
		GitHubToken:      githubToken,
		Source:           source,
	}, nil
}

// Run all checks and build all artifacts.
func (m *Ci) Ci(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		_, err := m.Test().Sync(ctx)

		return err
	})

	group.Go(func() error {
		return m.Lint().All(ctx)
	})

	// TODO: run trivy scan on container(s?)
	// TODO: version should be the commit hash (if any?)?
	group.Go(func() error {
		images := m.Build().containerImages("ci")

		for _, image := range images {
			_, err := image.Sync(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	// TODO: run trivy scan on helm chart
	group.Go(func() error {
		_, err := m.Build().HelmChart("0.0.0").Sync(ctx)

		return err
	})

	group.Go(func() error {
		files := m.releaseAssets("ci")

		for _, file := range files {
			_, err := file.Sync(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return group.Wait()
}

func (m *Ci) Test() *Container {
	return dag.Go(GoOpts{
		Version: goVersion,
	}).
		WithSource(m.Source).
		Exec([]string{"go", "test", "-v", "./..."})
}
