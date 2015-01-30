# pkgcloud

A collection of Go command-line tools to talk to the [PackageCloud API].

## Installation

    $ go get github.com/mlafeldt/pkgcloud/...

## Pushing packages

This is the only operation supported so far.

Make sure that `PACKAGECLOUD_TOKEN` is set in the environment.

    $ pkgcloud-push user/repo/distro/version /path/to/packages

Examples:

    $ pkgcloud-push mlafeldt/chef-runner/ubuntu/trusty chef-runner_0.8.0-2_amd64.deb

    $ pkgcloud-push mlafeldt/chef-runner/ubuntu/trusty *.deb


[PackageCloud API]: https://packagecloud.io/docs/api
