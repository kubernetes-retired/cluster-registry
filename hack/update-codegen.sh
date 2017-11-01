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

set -euo pipefail

SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCRIPT_BASE=$(cd ${SCRIPT_ROOT}/../..; pwd)
REPO_DIRNAME=$(basename $(dirname "${SCRIPT_ROOT}"))
TMP_GOPATH="$(mktemp -d /tmp/gopathXXXXXXXX)"
GEN_TMPDIR="$(mktemp -d /tmp/genXXXXXXXX)"

# Called on EXIT after the temporary directories are created.
function clean_up() {
  if [[ "${TMP_GOPATH}" == "/tmp/gopath"* ]]; then
    rm -rf "${TMP_GOPATH}"
  fi
  if [[ "${GEN_TMPDIR}" == "/tmp/gen"* ]]; then
    rm -rf "${GEN_TMPDIR}"
  fi
}
trap clean_up EXIT

# Generates code for the provided groupname ($1) and version ($2) using $3
# as the --output-base flag for all generation commands.
# To verify instead of generating, pass "--verify-only" as $4.
function generate_group() {
  local GROUP_NAME=$1
  local VERSION=$2
  local OUTPUT_BASE=$3
  local CLIENT_PKG=k8s.io/cluster-registry/pkg/client
  local CLIENTSET_PKG=${CLIENT_PKG}/clientset_generated
  local LISTERS_PKG=${CLIENT_PKG}/listers_generated
  local INFORMERS_PKG=${CLIENT_PKG}/informers_generated
  local APIS_PKG=k8s.io/cluster-registry/pkg/apis

  echo "generating clientset for group ${GROUP_NAME} and version ${VERSION} at ${SCRIPT_BASE}/${CLIENT_PKG}"
  bazel run //vendor/k8s.io/code-generator/cmd/client-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-base ${APIS_PKG} \
    --input ${GROUP_NAME} \
    --clientset-path ${CLIENTSET_PKG} \
    --output-base "${OUTPUT_BASE}" \
    --clientset-name "internalclientset" \
    "$4"
  bazel run //vendor/k8s.io/code-generator/cmd/client-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-base ${APIS_PKG} \
    --input ${GROUP_NAME}/${VERSION} \
    --clientset-path ${CLIENTSET_PKG} \
    --output-base "${OUTPUT_BASE}" \
    --clientset-name "clientset" \
    "$4"

  echo "generating listers for group ${GROUP_NAME} and version ${VERSION} at ${SCRIPT_BASE}/${LISTERS_PKG}"
  bazel run //vendor/k8s.io/code-generator/cmd/lister-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-dirs ${APIS_PKG}/${GROUP_NAME},${APIS_PKG}/${GROUP_NAME}/${VERSION} \
    --output-package ${LISTERS_PKG} \
    --output-base "${OUTPUT_BASE}" \
    "$4"

  echo "generating informers for group ${GROUP_NAME} and version ${VERSION} at ${SCRIPT_BASE}/${INFORMERS_PKG}"
  bazel run //vendor/k8s.io/code-generator/cmd/informer-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-dirs ${APIS_PKG}/${GROUP_NAME},${APIS_PKG}/${GROUP_NAME}/${VERSION} \
    --versioned-clientset-package ${CLIENT_PKG}/clientset_generated/clientset \
    --internal-clientset-package ${CLIENT_PKG}/clientset_generated/internalclientset \
    --listers-package ${LISTERS_PKG} \
    --output-package ${INFORMERS_PKG} \
    --output-base "${OUTPUT_BASE}" \
    "$4"

  echo "generating deep copies"
  bazel run //vendor/k8s.io/code-generator/cmd/deepcopy-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-dirs ${APIS_PKG}/${GROUP_NAME},${APIS_PKG}/${GROUP_NAME}/${VERSION} \
    --output-base "${OUTPUT_BASE}" \
    --output-file-base zz_generated.deepcopy \
    "$4"

  echo "generating defaults"
  bazel run //vendor/k8s.io/code-generator/cmd/defaulter-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-dirs ${APIS_PKG}/${GROUP_NAME},${APIS_PKG}/${GROUP_NAME}/${VERSION} \
    --output-base "${OUTPUT_BASE}" \
    --output-file-base zz_generated.defaults \
    "$4"

  echo "generating conversions"
  bazel run //vendor/k8s.io/code-generator/cmd/conversion-gen -- \
    --go-header-file "${SCRIPT_ROOT}/boilerplate/boilerplate.go.txt" \
    --input-dirs ${APIS_PKG}/${GROUP_NAME},${APIS_PKG}/${GROUP_NAME}/${VERSION} \
    --extra-peer-dirs "k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/conversion,k8s.io/apimachinery/pkg/runtime" \
    --output-base "${OUTPUT_BASE}" \
    --output-file-base zz_generated.conversion \
    "$4"
}

# Set up the temporary GOPATH with necessary dependencies.
mkdir -p "${TMP_GOPATH}/src/k8s.io/cluster-registry"
mkdir -p "${TMP_GOPATH}/src/k8s.io/apimachinery"
cp -r "${SCRIPT_ROOT}/../"* "${TMP_GOPATH}/src/k8s.io/cluster-registry"
cp -r "${SCRIPT_ROOT}/../vendor/k8s.io/apimachinery/"* "${TMP_GOPATH}/src/k8s.io/apimachinery"

# In verify mode, generate into the temporary GOPATH.
OUTPUT_BASE="${GEN_TMPDIR}"
if [ -n "$@" ]; then
  OUTPUT_BASE="${TMP_GOPATH}/src"
fi

# Perform the code generation.
export GOPATH="${TMP_GOPATH}"
generate_group clusterregistry v1alpha1 "${OUTPUT_BASE}" "${@-}"

# In generate mode, copy the generated files back into the tree.
if [ -n "$@" ]; then
  cp -r "${OUTPUT_BASE}/k8s.io/cluster-registry/"* "${SCRIPT_BASE}/${REPO_DIRNAME}"
fi
