#!/bin/bash

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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCRIPT_BASE=$(cd ${SCRIPT_ROOT}/../..; pwd)
REPO_DIRNAME=$(basename $(dirname "${SCRIPT_ROOT}"))
TMP_GOPATH="$(mktemp -d /tmp/gopathXXXXXXXX)"

# Called on EXIT after the temporary directory is created.
function clean_up() {
  if [[ "${TMP_GOPATH}" == "/tmp/gopath"* ]]; then
    rm -rf "${TMP_GOPATH}"
  fi
}
trap clean_up EXIT

# Generates code for the provided groupname ($1) and version ($2).
function generate_group() {
  local GROUP_NAME=$1
  local VERSION=$2
  local CLIENT_PKG=k8s.io/cluster-registry/pkg/client
  local LISTERS_PKG=${CLIENT_PKG}/listers_generated
  local INFORMERS_PKG=${CLIENT_PKG}/informers_generated
  local APIS_PKG=k8s.io/cluster-registry/pkg/apis
  local INPUT_APIS=(
    ${GROUP_NAME}/
    ${GROUP_NAME}/${VERSION}
  )

  local GEN_TMPDIR="$(mktemp -d /tmp/genXXXXXXXX)"

  echo "generating clientset for group ${GROUP_NAME} and version ${VERSION} at ${SCRIPT_BASE}/${CLIENT_PKG}"
  bazel run @io_k8s_code_generator//cmd/client-gen -- --input-base ${APIS_PKG} --input ${INPUT_APIS[@]} --clientset-path ${CLIENT_PKG}/clientset_generated --output-base ${GEN_TMPDIR}
  bazel run @io_k8s_code_generator//cmd/client-gen -- --clientset-name "clientset" --input-base ${APIS_PKG} --input ${GROUP_NAME}/${VERSION} --clientset-path ${CLIENT_PKG}/clientset_generated --output-base ${GEN_TMPDIR}

  echo "generating listers for group ${GROUP_NAME} and version ${VERSION} at ${SCRIPT_BASE}/${LISTERS_PKG}"
  bazel run @io_k8s_code_generator//cmd/lister-gen -- --input-dirs ${APIS_PKG}/${GROUP_NAME},${APIS_PKG}/${GROUP_NAME}/${VERSION} --output-package ${LISTERS_PKG} --output-base ${GEN_TMPDIR}

  echo "generating informers for group ${GROUP_NAME} and version ${VERSION} at ${SCRIPT_BASE}/${INFORMERS_PKG}"
  bazel run @io_k8s_code_generator//cmd/informer-gen -- \
    --input-dirs ${APIS_PKG}/${GROUP_NAME} --input-dirs ${APIS_PKG}/${GROUP_NAME}/${VERSION} \
    --versioned-clientset-package ${CLIENT_PKG}/clientset_generated/clientset \
    --internal-clientset-package ${CLIENT_PKG}/clientset_generated/internalclientset \
    --listers-package ${LISTERS_PKG} \
    --output-package ${INFORMERS_PKG} \
    --output-base ${GEN_TMPDIR}

  # The following generators do not respect the --output-package, so generate
  # their output into a temporary location and rsync it over.

  echo "generating deep copies"
  bazel run @io_k8s_code_generator//cmd/deepcopy-gen -- --input-dirs ${APIS_PKG}/${GROUP_NAME} --input-dirs ${APIS_PKG}/${GROUP_NAME}/${VERSION} --output-base ${GEN_TMPDIR}

  echo "generating defaults"
  bazel run @io_k8s_code_generator//cmd/defaulter-gen -- --input-dirs ${APIS_PKG}/${GROUP_NAME} --input-dirs ${APIS_PKG}/${GROUP_NAME}/${VERSION} --output-base ${GEN_TMPDIR}

  echo "generating conversions"
  bazel run @io_k8s_code_generator//cmd/conversion-gen -- --input-dirs ${APIS_PKG}/${GROUP_NAME} --input-dirs ${APIS_PKG}/${GROUP_NAME}/${VERSION} --output-base ${GEN_TMPDIR}

  rsync -a "${GEN_TMPDIR}/k8s.io/cluster-registry/" "${SCRIPT_BASE}/${REPO_DIRNAME}"

  if [[ "${GEN_TMPDIR}" == "/tmp/gen"* ]]; then
    rm -rf "${GEN_TMPDIR}"
  fi
}

echo "Creating temporary GOPATH from bazel workspace"
${SCRIPT_ROOT}/gopath_from_workspace.sh "${TMP_GOPATH}"
export GOPATH="${TMP_GOPATH}"
generate_group clusterregistry v1alpha1
