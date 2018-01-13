#!/usr/bin/env bash
# Copyright 2018 The Kubernetes Authors.
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

# This script runs all of the Go source verification scripts inside of the
# ./hack/go-tools directory. The success or failure of each script is outputted
# in green or red colored text, respectively. If any script fails, an error is
# returned, otherwise returns 0.

set -euo pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_ROOT}/.." && pwd)"
TMP_GOPATH="$(mktemp -d /tmp/gopathXXXXXXXX)"
TMP_REPO_ROOT="${TMP_GOPATH}/src/k8s.io/cluster-registry"
TMP_GO_TOOLS_DIR="${TMP_GOPATH}/src/k8s.io/cluster-registry/hack/go-tools"

# Called on EXIT after the temporary directories are created.
function clean_up() {
  if [[ "${TMP_GOPATH}" == "/tmp/gopath"* ]]; then
    rm -rf "${TMP_GOPATH}"
  fi
}
trap clean_up EXIT

function run-checks {
  local -r pattern=$1

  for t in $(ls ${pattern})
  do
    echo -e "Verifying ${t}"
    local start=$(date +%s)
    cd ${TMP_REPO_ROOT} && "${t}" && tr=$? || tr=$?
    local elapsed=$(($(date +%s) - ${start}))
    if [[ ${tr} -eq 0 ]]; then
      echo -e "${color_green}SUCCESS${color_norm}  ${t}\t${elapsed}s"
    else
      echo -e "${color_red}FAILED${color_norm}   ${t}\t${elapsed}s"
      ret=1
    fi
  done
}

# Set up the temporary GOPATH. This helps to run this script in a vanilla
# environment, e.g. Prow, because the verify-govet.sh script called by this one
# runs 'go list' which will print the import path for each package (e.g.
# k8s.io/cluster-registry/...) to pass into 'go vet'. This results in go not
# finding the package correctly if either GOPATH doesn't exist (in the case of
# the current bazelbuild image), or the package exists but in a different path.
# So this temporary GOPATH is set up to replicate the same import path
# location.
mkdir -p "${TMP_REPO_ROOT}"
cp -r "${REPO_ROOT}/"* "${TMP_REPO_ROOT}"
export GOPATH="${TMP_GOPATH}"

echo "Working directory: ${TMP_REPO_ROOT}"

# Some useful colors.
if [[ -z "${color_start-}" ]]; then
  declare -r color_start="\033["
  declare -r color_red="${color_start}0;31m"
  declare -r color_yellow="${color_start}0;33m"
  declare -r color_green="${color_start}0;32m"
  declare -r color_norm="${color_start}0m"
fi

ret=0
run-checks "${TMP_GO_TOOLS_DIR}/*.sh"
exit ${ret}
