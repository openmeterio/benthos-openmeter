package main

import (
	"path/filepath"
	"strings"
	"time"
)

// Build individual artifacts. (Useful for testing and development)
func (m *Ci) Build() *Build {
	return &Build{
		Source: m.Source,
	}
}

type Build struct {
	// +private
	Source *Directory
}

func (m *Build) containerImages(version string) []*Container {
	platforms := []Platform{
		"linux/amd64",
		"linux/arm64",
	}

	variants := make([]*Container, 0, len(platforms))

	for _, platform := range platforms {
		variants = append(variants, m.containerImage(platform, Opt(version)))
	}

	return variants
}

// Build a container image.
func (m *Build) ContainerImage(
	// Platform in the format of OS/ARCH[/VARIANT] (eg. "darwin/arm64/v7")
	platform Optional[Platform],
) *Container {
	return m.containerImage(platform.GetOr(""), OptEmpty[string]())
}

func (m *Build) containerImage(platform Platform, version Optional[string]) *Container {
	return dag.Container(ContainerOpts{Platform: platform}).
		From(alpineBaseImage).
		WithLabel("org.opencontainers.image.title", "benthos-openmeter").
		WithLabel("org.opencontainers.image.description", "Ingest events into OpenMeter from everywhere").
		WithLabel("org.opencontainers.image.url", "https://github.com/openmeterio/benthos-openmeter").
		WithLabel("org.opencontainers.image.created", time.Now().String()). // TODO: embed commit timestamp
		WithLabel("org.opencontainers.image.source", "https://github.com/openmeterio/benthos-openmeter").
		WithLabel("org.opencontainers.image.licenses", "Apache-2.0").
		With(func(c *Container) *Container {
			if v, ok := version.Get(); ok {
				c = c.WithLabel("org.opencontainers.image.version", v)
			}

			return c
		}).
		WithExec([]string{"apk", "add", "--update", "--no-cache", "ca-certificates", "tzdata", "bash"}).
		WithWorkdir("/etc/openmeter").
		WithFile("/etc/openmeter/cloudevents.spec.json", m.Source.File("cloudevents.spec.json")).
		WithFile("/etc/openmeter/examples/http-server/input.yaml", m.Source.File("examples/http-server/input.yaml")).
		WithFile("/etc/openmeter/examples/http-server/output.yaml", m.Source.File("examples/http-server/output.yaml")).
		WithFile("/etc/openmeter/examples/kubernetes-pod-exec-time/config.yaml", m.Source.File("examples/kubernetes-pod-exec-time/config.yaml")).
		WithFile("/usr/local/bin/benthos", m.binary(platform, version))
}

// Build a binary.
func (m *Build) Binary(
	// Platform in the format of OS/ARCH[/VARIANT] (eg. "darwin/arm64/v7")
	platform Optional[Platform],
) *File {
	return m.binary(platform.GetOr(""), OptEmpty[string]())
}

func (m *Build) binary(platform Platform, version Optional[string]) *File {
	return dag.Go(GoOpts{
		Version: goVersion,
	}).
		WithPlatform(string(platform)).
		WithCgoDisabled().
		WithSource(m.Source).
		Build(GoWithSourceBuildOpts{
			Trimpath: true,
			RawArgs: []string{
				"-ldflags",
				"-s -w -X main.version=" + version.GetOr("unknown"),
			},
		})
}

func (m *Build) HelmChart(version Optional[string]) *File {
	chart := m.helmChartDir()

	var opts HelmPackageOpts

	if v, ok := version.Get(); ok {
		opts.Version = strings.TrimPrefix(v, "v")
		opts.AppVersion = v
	}

	return dag.Helm(HelmOpts{Version: helmVersion}).Package(chart, opts)
}

func (m *Build) helmChartDir() *Directory {
	chart := dag.Host().Directory(filepath.Join(root(), "deploy/charts/benthos-openmeter"), HostDirectoryOpts{
		Exclude: []string{"charts"}, // exclude dependencies
	})

	readme := dag.HelmDocs(HelmDocsOpts{Version: helmDocsVersion}).Generate(chart, HelmDocsGenerateOpts{
		Templates: []*File{
			dag.Host().File(filepath.Join(root(), "deploy/charts/template.md")),
			dag.Host().File(filepath.Join(root(), "deploy/charts/benthos-openmeter/README.tmpl.md")),
		},
		SortValuesOrder: "file",
	})

	return chart.WithFile("README.md", readme)
}
