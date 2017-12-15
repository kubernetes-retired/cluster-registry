# Creating a release

> This release process is subject to change as the cluster-registry evolves.

Please see the [development doc](development.md#release-and-build-versioning)
for some more information about the release tools.

## Release process

You will need to have permissions to create a release on the cluster-registry
repo, as well as permissions for the `crreleases` GCP project, in order to run
this release process. We are working on determining how to limit the amount of
special privilege necessary to do a release.

1. Create a
   [new release](https://github.com/kubernetes/cluster-registry/releases/new)
   on the GitHub Releases page for the cluster registry. Choose the latest
   commit (or another commit if you have a particular reason not to choose the
   latest commit) and a tag name with the scheme vX.Y.Z. Name the release
   `vX.Y.Z`. Leave the body empty; it will be added later.
2. Pull the latest version of the cluster-registry repo locally, with the tag
   you just created. Check out that tag: `git checkout tags/vX.Y.Z`
3. Run `hack/release.sh vX.Y.Z >relnotes`. This will require permissions for the
   `crreleases` GCP project, which you may not have. We are working on
   automating this step so that it does not require anything to be done on a
   local machine.
4. Paste the contents of the `relnotes` file into the body of the release.
5. Send an announcement to
   [kubernetes-sig-multicluster](https://groups.google.com/forum/#!forum/kubernetes-sig-multicluster).

## Notes

- The cluster-registry does not use branches for its releases. As it becomes
  necessary, we will evaluate branching strategies.
- There is no verification process for releases. Since each commit is currently
  checked by per-PR tests that run the full suite of tests we have, we expect
  all commits to be green and suitable for release.
