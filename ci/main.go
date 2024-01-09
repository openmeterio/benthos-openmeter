package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

const (
	goVersion           = "1.21.5"
	golangciLintVersion = "v1.54.2"

	alpineBaseImage = "alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48"

	helmDocsVersion = "v1.11.3"
	helmVersion     = "3.13.2"
)

const imageRepo = "ghcr.io/openmeterio/benthos-openmeter"

type Ci struct {
	// +private
	RegistryUser string

	// +private
	RegistryPassword *Secret

	// Project source directory
	// This will become useful once pulling from remote becomes available
	//
	// +private
	Source *Directory
}

func New(
	// Checkout the repository (at the designated ref) and use it as the source directory instead of the local one.
	checkout Optional[string],

	// Container registry user (required for pushing images).
	registryUser Optional[string],

	// Container registry password (required for pushing images).
	registryPassword Optional[*Secret],
) (*Ci, error) {
	var source *Directory

	if refName, ok := checkout.Get(); ok {
		source = dag.Git("https://github.com/openmeterio/benthos-openmeter.git", GitOpts{
			KeepGitDir: true,
		}).Branch(refName).Tree()
	} else {
		source = projectDir()
	}

	return &Ci{
		RegistryUser:     registryUser.GetOr(""),
		RegistryPassword: registryPassword.GetOr(nil),
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
		_, err := m.Lint().Sync(ctx)

		return err
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
		_, err := m.Build().HelmChart(Opt("0.0.0")).Sync(ctx)

		return err
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

func (m *Ci) Lint() *Container {
	return dag.GolangciLint(GolangciLintOpts{
		Version:   golangciLintVersion,
		GoVersion: goVersion,
	}).
		Run(m.Source, GolangciLintRunOpts{
			Verbose: true,
		})
}

// Build and publish a snapshot of all artifacts from the current development version.
func (m *Ci) Snapshot(ctx context.Context) error {
	// TODO: capture branch name and push it as tag/version
	// TODO: version should be a combination of branch name and build time?
	return m.pushImages(ctx, "latest", []string{"latest", "main"})
}

// Build and publish all release artifacts.
func (m *Ci) Release(ctx context.Context, version string) error {
	var group errgroup.Group

	group.Go(func() error {
		// Disable pushing images for now
		return nil

		// TODO: refuse to publish release artifacts in a dirty git dir or when there is no tag pointing to the current ref
		return m.pushImages(ctx, version, []string{version})
	})

	group.Go(func() error {
		username, password := m.RegistryUser, m.RegistryPassword

		if username == "" {
			return errors.New("registry user is required to push helm charts to ghcr.io")
		}

		if password == nil {
			return errors.New("registry password is required to push helm charts to ghcr.io")
		}

		chart := m.Build().HelmChart(Opt(version))

		_, err := dag.Helm(HelmOpts{Version: helmVersion}).
			Login("ghcr.io", username, password).
			Push(chart, "oci://ghcr.io/openmeterio/helm-charts").
			Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}

func (m *Ci) pushImages(ctx context.Context, version string, tags []string) error {
	username, password := m.RegistryUser, m.RegistryPassword

	if username == "" {
		return errors.New("registry user is required to push images to ghcr.io")
	}

	if password == nil {
		return errors.New("registry password is required to push images to ghcr.io")
	}

	images := m.Build().containerImages(version)

	var group errgroup.Group

	for _, tag := range tags {
		tag := tag

		group.Go(func() error {
			_, err := dag.Container().
				WithRegistryAuth("ghcr.io", username, password).
				Publish(ctx, fmt.Sprintf("%s:%s", imageRepo, tag), ContainerPublishOpts{
					PlatformVariants: images,
				})
			if err != nil {
				return err
			}

			return nil
		})
	}

	return group.Wait()
}
