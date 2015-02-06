# pkgcloud

[![GoDoc](https://godoc.org/github.com/mlafeldt/pkgcloud?status.svg)](https://godoc.org/github.com/mlafeldt/pkgcloud)

Talk to the [PackageCloud API](https://packagecloud.io/docs/api), in Go.

## Installation

    $ go get github.com/mlafeldt/pkgcloud/...

## API Usage

See [Godoc](https://godoc.org/github.com/mlafeldt/pkgcloud) and
[pkgcloud.go](pkgcloud.go) to learn about the API.

## Client Usage

### Pushing packages

Pushing packages with `pkgcloud-push` is the only operation supported so far.
The tool is a simple and fast replacement for the original `package_cloud push`
command. If you pass more than one package, `pkgcloud-push` will push them in
parallel! Before using it, however, make sure that `PACKAGECLOUD_TOKEN` is set
in your environment.

Usage:

    $ pkgcloud-push user/repo/distro/version /path/to/packages

Examples:

    $ pkgcloud-push mlafeldt/chef-runner/ubuntu/trusty chef-runner_0.8.0-2_amd64.deb

    $ pkgcloud-push mlafeldt/chef-runner/ubuntu/trusty *.deb
