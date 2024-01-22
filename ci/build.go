package main

import (
	"fmt"
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
		variants = append(variants, m.containerImage(platform, version))
	}

	return variants
}

// Build a container image.
func (m *Build) ContainerImage(
	// Target platform in "[os]/[platform]/[version]" format (e.g., "darwin/arm64/v7", "windows/amd64", "linux/arm64").
	// +optional
	platform Platform,
) *Container {
	return m.containerImage(platform, "")
}

func (m *Build) containerImage(platform Platform, version string) *Container {
	return dag.Container(ContainerOpts{Platform: platform}).
		From(alpineBaseImage).
		WithLabel("org.opencontainers.image.title", "benthos-openmeter").
		WithLabel("org.opencontainers.image.description", "Ingest events into OpenMeter from everywhere").
		WithLabel("org.opencontainers.image.url", "https://github.com/openmeterio/benthos-openmeter").
		WithLabel("org.opencontainers.image.created", time.Now().String()). // TODO: embed commit timestamp
		WithLabel("org.opencontainers.image.source", "https://github.com/openmeterio/benthos-openmeter").
		WithLabel("org.opencontainers.image.licenses", "Apache-2.0").
		With(func(c *Container) *Container {
			if version != "" {
				c = c.WithLabel("org.opencontainers.image.version", version)
			}

			return c
		}).
		WithExec([]string{"apk", "add", "--update", "--no-cache", "ca-certificates", "tzdata", "bash"}).
		WithWorkdir("/etc/benthos").
		WithFile("/etc/benthos/cloudevents.spec.json", m.Source.File("cloudevents.spec.json")).
		WithFile("/etc/benthos/examples/http-server/input.yaml", m.Source.File("examples/http-server/input.yaml")).
		WithFile("/etc/benthos/examples/http-server/output.yaml", m.Source.File("examples/http-server/output.yaml")).
		WithFile("/etc/benthos/examples/kubernetes-pod-exec-time/config.yaml", m.Source.File("examples/kubernetes-pod-exec-time/config.yaml")).
		WithFile("/usr/local/bin/benthos", m.binary(platform, version))
}

// Build a binary.
func (m *Build) Binary(
	// Target platform in "[os]/[platform]/[version]" format (e.g., "darwin/arm64/v7", "windows/amd64", "linux/arm64").
	// +optional
	platform Platform,
) *File {
	return m.binary(platform, "")
}

func (m *Build) binary(platform Platform, version string) *File {
	if version == "" {
		version = "unknown"
	}

	return dag.Go(GoOpts{
		Version: goVersion,
	}).
		WithModuleCache(dag.CacheVolume("benthos-openmeter-go-mod")).
		WithBuildCache(dag.CacheVolume("benthos-openmeter-go-build")).
		WithPlatform(string(platform)).
		WithCgoDisabled().
		WithSource(m.Source).
		Build(GoWithSourceBuildOpts{
			Name:     "benthos",
			Trimpath: true,
			RawArgs: []string{
				"-ldflags",
				"-s -w -X main.version=" + version,
			},
		})
}

func (m *Build) binaryArchives(version string) []*File {
	platforms := []Platform{
		"linux/amd64",
		"linux/arm64",

		"darwin/amd64",
		"darwin/arm64",
	}

	archives := make([]*File, 0, len(platforms))

	for _, platform := range platforms {
		archives = append(archives, m.binaryArchive(version, platform))
	}

	return archives
}

func (m *Build) binaryArchive(version string, platform Platform) *File {
	var archiver interface {
		Archive(name string, source *Directory) *File
	} = dag.Archivist().TarGz()

	if strings.HasPrefix(string(platform), "windows/") {
		archiver = dag.Archivist().Zip()
	}

	return archiver.Archive(
		fmt.Sprintf("benthos_%s", strings.ReplaceAll(string(platform), "/", "_")),
		dag.Directory().
			WithFile("", m.binary(platform, version)).
			WithFile("", m.Source.File("README.md")).
			WithFile("", m.Source.File("LICENSE")),
	)
}

func (m *Build) checksums(files []*File) *File {
	return dag.Container().
		From(alpineBaseImage).
		WithWorkdir("/work").
		With(func(c *Container) *Container {
			dir := dag.Directory()

			for _, file := range files {
				dir = dir.WithFile("", file)
			}

			return c.WithMountedDirectory("/work", dir)
		}).
		WithExec([]string{"sh", "-c", "sha256sum $(ls) > checksums.txt"}).
		File("/work/checksums.txt")
}

func (m *Build) HelmChart(
	// Release version.
	// +optional
	version string,
) *File {
	chart := helmChartDir(m.Source)

	var opts HelmPackageOpts

	if version != "" {
		opts.Version = strings.TrimPrefix(version, "v")
		opts.AppVersion = version
	}

	return dag.Helm(HelmOpts{Version: helmVersion}).Package(chart, opts)
}

func helmChartDir(source *Directory) *Directory {
	chart := source.Directory("deploy/charts/benthos-openmeter")

	readme := dag.HelmDocs(HelmDocsOpts{Version: helmDocsVersion}).Generate(chart, HelmDocsGenerateOpts{
		Templates: []*File{
			source.File("deploy/charts/template.md"),
			source.File("deploy/charts/benthos-openmeter/README.tmpl.md"),
		},
		SortValuesOrder: "file",
	})

	return chart.WithFile("README.md", readme)
}
