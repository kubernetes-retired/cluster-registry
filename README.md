# cluster-registry

This repository contains the code for the cluster registry. The cluster registry
is an effort being led under the auspices of sig-multicluster.

This is currently in the prototyping stage, and is not yet ready for use except
by contributors and experimenters.

If you have questions, please reach out to
[kubernetes-sig-federation](https://groups.google.com/forum/#!forum/kubernetes-sig-federation).

[Cluster Registry API
design](https://docs.google.com/document/d/1Oi9EO3Jwtp69obakl-9YpLkP764GZzsz95XJlX1a960/edit)

# Development

## Prerequisites

You must have a recent verzion of [bazel](https://bazel.io) installed. Bazel is
the recommended way to build and test the cluster registry. Bazel is designed to
maintain compatibility with standard Go tooling, but this is not tested on a
regular basis, and some scripts/tooling in this project are built around Bazel.

Before doing any development work, you must (in order, from the repository root
directory, after cloning):

*   run `update-genfiles.sh`
*   run `bazel run //:gazelle`

## Building crinit

From the root of the repository:

1.  Run `bazel build //cmd/crinit`.

## Building clusterregistry

1.  Run `bazel build //cmd/clusterreigstry`
1.  If you want to build a docker image, run `blaze build
    //cmd/clusterregistry-image`
1.  To push an image to Google Container registry, you'll need to run `blaze run
    //cmd/push-clusterregistry-image --define project=<your_project_id>`

## Run all tests

You can run all the unit tests by running `bazel test ...` from the repository
root.
