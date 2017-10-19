# Development

## Prerequisites

You must have a recent version of [bazel](https://bazel.io) installed. Bazel is
the recommended way to build and test the cluster registry. Bazel is designed to
maintain compatibility with standard Go tooling, but this is not tested on a
regular basis, and some scripts/tooling in this project are built around Bazel.

NOTE: There is an issue with bazel 0.6.x. As a workaround, use 0.5.x, or pass
the flag `--incompatible_comprehension_variables_do_not_leak=false` to bazel
0.6.x invocations.

To push an image to Google Container Registry you'll also have to have 
`docker-credential-gcr` installed and configured. This allows for Docker clients 
v1.11+ to easily make authenticated requests to GCR's repositories (gcr.io, 
eu.gcr.io, etc.):

1. Run `gcloud components install docker-credential-gcr`
1. Run `docker-credential-gcr configure-docker`

You'll need to clone the repository before doing any work. It's expedient to
clone into $GOPATH/src/k8s.io/cluster-registry, since some Kubernetes and go
tooling expect this, but the repository itself is location-agnostic.

Before doing any development work, you must (in order, from the repository root
directory, after cloning):

1.  run `./update-codegen.sh`
1.  run `bazel run //:gazelle`

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
`bazel test -- ... -//pkg/client/... -//cmd/clusterregistry:push-clusterregistry-image -//pkg/crinit/testing/...`
from the repository root. (This may take a while the first time you run it.)
