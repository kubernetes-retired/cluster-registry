# Development

## Prerequisites

You'll need to clone the repository before doing any work. Make sure to clone
into $GOPATH/src/k8s.io/cluster-registry, since much of the tooling expects
this.

Before doing any development work, you must (in order, from the repository root
directory, after cloning):

1.  run `bazel run //:gazelle`

### bazel

You must have a recent version of [bazel](https://bazel.io) installed. Bazel is
the recommended way to build and test the cluster registry. Bazel is designed to
maintain compatibility with standard Go tooling, but this is not tested on a
regular basis, and some scripts/tooling in this project are built around Bazel.

NOTE: There is an issue with bazel 0.6.x. As a workaround, use 0.5.x, or pass
the flag `--incompatible_comprehension_variables_do_not_leak=false` to bazel
0.6.x invocations.

### `docker-credential-gcr`

To push an image to Google Container Registry you'll also have to have
`docker-credential-gcr` installed and configured. This allows for Docker clients
v1.11+ to easily make authenticated requests to GCR's repositories (gcr.io,
eu.gcr.io, etc.):

1.  Run `gcloud components install docker-credential-gcr`
1.  Run `docker-credential-gcr configure-docker`

### dep

This repository maintains its `vendor` directory with
[dep](https://github.com/golang/dep). You must have v0.3.2 or newer of the tool
installed if you intend to update the vendored dependencies.

## Building `clusterregistry`

`clusterregistry` is the binary for the Kubernetes API server that serves the
cluster registry API.

To build it, from the root of the repository:

1.  Run `bazel build //cmd/clusterregistry`. (This may take a while the first
    time you run it.)
1.  If you want to build a docker image, run
    `bazel build //cmd/clusterregistry:clusterregistry-image`
1.  To push an image to Google Container registry, you'll need to run
    `bazel run //cmd/clusterregistry:push-clusterregistry-image --define project=<your_project_id>`

## Building `crinit`

`crinit` is a command-line tool to boostrap a cluster registry into a Kubernetes
cluster.

To build it, from the root of the repository:

1.  Run `bazel build //cmd/crinit`. (This may take a while the first time you
    run it.)

## Run all tests

You can run all the unit tests by running
`bazel test -- //cmd/... //pkg/...  -//cmd/clusterregistry:push-clusterregistry-image -//pkg/client/...`
from the repository root. (This may take a while the first time you run it.)

## Updating Bazel files

You will need to update the BUILD and BUILD.bazel files when making changes that
cause the Go imports to change.

1.  Run `./hack/update-bazel.sh`
1.  Add the updated `BUILD.bazel` and `BUILD` files to your commit.

## Updating vendored dependencies

The `dep` tool is currently only marginally supported by k/ repos. There are
some warts.

1.  Use the `dep` tool and/or modify the `Gopkg.toml` file to reference the
    new dependency versions. Refer to the [dep
    docs](https://github.com/golang/dep#usage) for more info.
1.  Run `dep prune`.
1.  At this point, there may be other modifications necessary either to the
    vendored dependencies or the `Gopkg.toml` file. The known ones are noted
    below. Make these and any additional necessary ones, and add them to this
    list.
    -   As of #58, it was necessary to modify the BUILD file in
        `vendor/k8s.io/client-go/util/cert` to have the go_library not reference
        testdata.
1.  [Update the BUILD and BUILD.bazel files](#updating-bazel-files).
1.  Run the tests and fix any breakages.
1.  When sending out a PR, please put the handmade changes in one commit and the
    generated updates in another commit so that it's easier for reviewers to
    see what's been done.

## Updating generated code

If you modify any files in `pkg/apis`, you will likely need to regenerate the
generated clients and other generated files.

1.  Run `./hack/update-codegen.sh` to update the files.
1.  Add the generated files to your PR, preferably in a separate,
    generated-only commit so that they are easier to review.
