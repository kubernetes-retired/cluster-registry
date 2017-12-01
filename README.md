# Cluster Registry

A lightweight tool for maintaining a list of clusters and associated metadata.

# What is it?

The cluster registry helps you keep track of and perform operations on your
clusters. This repository contains an implementation of the cluster API
([code](https://github.com/kubernetes/cluster-registry/tree/master/pkg/apis/clusterregistry),
[design](https://github.com/kubernetes/cluster-registry/tree/master/docs/api_design.md))
backed by a Kubernetes-style API server, which is the canonical implementation
of the cluster registry.

# Documentation

Documentation is in the
[`/docs`](https://github.com/kubernetes/cluster-registry/tree/master/docs/)
directory.

# Getting involved

The cluster registry is still a young project, but we welcome your
contributions, suggestions and input! Please reach out to the
[kubernetes-sig-multicluster](https://groups.google.com/forum/#!forum/kubernetes-sig-multicluster)
mailing list, or find us on
[Slack](https://github.com/kubernetes/community/blob/master/communication.md#social-media)
in [#sig-multicluster](https://kubernetes.slack.com/messages/sig-multicluster/).

## Maintainers

-   [@perotinus](https://github.com/perotinus)
-   [@font](https://github.com/font)
-   [@madhusudancs](https://github.com/madhusudancs)

# Development

There is a [nightly build here](https://k8s-testgrid.appspot.com/sig-multicluster-cluster-registry)

Basic instructions for working in the cluster-registry repo are
[here](https://github.com/kubernetes/cluster-registry/tree/master/docs/development.md).
