package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Build and publish a snapshot of all artifacts from the current development version.
func (m *Ci) Snapshot(ctx context.Context) error {
	// TODO: capture branch name and push it as tag/version
	// TODO: version should be a combination of branch name and build time?
	return m.pushImages(ctx, "latest", []string{"latest", "main"})
}

// Build and publish all release artifacts.
func (m *Ci) Release(
	ctx context.Context,

	// Release version.
	version string,
) error {
	group, ctx := errgroup.WithContext(ctx)

	// Container images
	group.Go(func() error {
		// Disable pushing images for now
		return nil

		// TODO: refuse to publish release artifacts in a dirty git dir or when there is no tag pointing to the current ref
		return m.pushImages(ctx, version, []string{version})
	})

	// Binaries
	group.Go(func() error {
		if m.GitHubToken == nil {
			return errors.New("GitHub token is required to publish a release")
		}

		releaseAssets := m.releaseAssets(version)

		_, err := dag.Gh(GhOpts{
			Version: "",
			Token:   m.GitHubToken,
			Repo:    "openmeterio/benthos-openmeter",
		}).Release().Create(version, version, GhReleaseCreateOpts{
			Files:         releaseAssets,
			GenerateNotes: true,
			Latest:        true,
			VerifyTag:     true,
		}).Sync(ctx)

		return err
	})

	// Helm chart
	group.Go(func() error {
		username, password := m.RegistryUser, m.RegistryPassword

		if username == "" {
			return errors.New("registry user is required to push helm charts to ghcr.io")
		}

		if password == nil {
			return errors.New("registry password is required to push helm charts to ghcr.io")
		}

		chart := m.Build().HelmChart(version)

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

func (m *Ci) releaseAssets(version string) []*File {
	binaryArchives := m.Build().binaryArchives(version)
	checksums := m.Build().checksums(binaryArchives)

	return append(binaryArchives, checksums)
}
