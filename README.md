# fl - the experimental CLI for Flintlock

`fl` is an experimental command-line client for [Flintlock](https://github.com/liquidmetal-dev/flintlock), the microVM
control-plane component of the [liquidmetal-dev](https://github.com/liquidmetal-dev) project. It talks to a Flintlock
host over gRPC to create, list, and delete microVMs.

## Prerequisites

- [Go 1.23](https://go.dev/dl/)
- [GoReleaser](https://goreleaser.com/) (used to build/release binaries)

This repo includes a [mise](https://mise.jdx.dev/) config (`mise.toml`) pinning the tool versions above. If you use
mise, just run:

```sh
mise install
```

## Building

```sh
make build
```

This runs `goreleaser release --snapshot --clean`, producing binaries for linux/darwin/windows and apk/deb/rpm
packages under `dist/`.

## Usage

```sh
fl version
```

Displays version information.

```sh
fl microvm create --host <flintlock-host:port> --name <name> [flags]
```

Creates a new microVM on the given Flintlock host. Use `--name-autogenerate` instead of `--name` to have a name
generated for you. See `fl microvm create --help` for the full set of flags (vcpu/memory sizing, kernel and root
images, network interfaces, cloud-init metadata, additional volumes, etc.).

```sh
fl microvm get --host <flintlock-host:port> [vmid]
```

Lists all microVMs on a host, or fetches a single microVM by id.

```sh
fl microvm delete --host <flintlock-host:port> <vmid>
```

Deletes a microVM from a host.

## Releases

Tagged pushes (`v*.*.*`) are built and published automatically via GoReleaser (see
`.github/workflows/release.yaml` and `.goreleaser.yaml`).

## License

[Mozilla Public License 2.0](LICENSE)
