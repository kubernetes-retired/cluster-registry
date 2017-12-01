#!/usr/bin/env bash
# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

ROOT_RELEASE_DIR="${ROOT_RELEASE_DIR:-nightly}"
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TMPDIR="$(mktemp -d /tmp/crrelXXXXXXXX)"

function clean_up() {
  if [[ "${TMPDIR}" == "/tmp/crrel"* ]]; then
    rm -rf "${TMPDIR}"
  fi
}
trap clean_up EXIT

BUILD_DATE="$(TZ=Etc/UTC date +%Y%m%d)"
RELEASE_DIR="${TMPDIR}/${BUILD_DATE}"

# Check for and install necessary dependencies.
command -v bazel >/dev/null 2>&1 || { echo >&2 "Please install bazel before running this script."; exit 1; }
command -v gcloud >/dev/null 2>&1 || { echo >&2 "Please install gcloud before running this script."; exit 1; }
gcloud components install gsutil docker-credential-gcr
docker-credential-gcr configure-docker

cd "${SCRIPT_ROOT}/.."

# Build the tools.
bazel build //cmd/crinit //cmd/clusterregistry

# Create the archives.
mkdir -p "${RELEASE_DIR}"
tar czf "${RELEASE_DIR}/clusterregistry-client.tar.gz" -C bazel-bin/cmd/crinit crinit
tar czf "${RELEASE_DIR}/clusterregistry-server.tar.gz" -C bazel-bin/cmd/clusterregistry clusterregistry

# Create the `latest` file
echo "${BUILD_DATE}" > "${TMPDIR}/latest"

pushd "${RELEASE_DIR}"

# Create the SHAs.
sha256sum clusterregistry-client.tar.gz > clusterregistry-client.tar.gz.sha
sha256sum clusterregistry-server.tar.gz > clusterregistry-server.tar.gz.sha

popd

# Upload the files to GCS.
gsutil cp -r "${TMPDIR}"/* "gs://crreleases/${ROOT_RELEASE_DIR}"

# Push the server container image.
bazel run //cmd/clusterregistry:push-clusterregistry-image --define repository=crreleases/nightly/clusterregistry

# Adjust the tags on the image. The `push-clusterregistry-image` rule tags the
# pushed image with the `dev` tag by default. This consistent tag allows the
# tool to easily add other tags to the image. The tool then removes the `dev`
# tag since this is not a development image.
gcloud container images add-tag --quiet \
  gcr.io/crreleases/nightly/clusterregistry:dev \
  gcr.io/crreleases/nightly/clusterregistry:${BUILD_DATE}
gcloud container images add-tag --quiet \
  gcr.io/crreleases/nightly/clusterregistry:dev \
  gcr.io/crreleases/nightly/clusterregistry:latest_nightly
gcloud container images untag --quiet \
  gcr.io/crreleases/nightly/clusterregistry:dev
