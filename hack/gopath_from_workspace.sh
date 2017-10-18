#!/bin/sh

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

# Creates a populated GOPATH from the repositories in a bazel workspace.
# Used to ensure that code generation scripts are running against the versions
# of external libraries in the workspace.
#
# Requires a temporary directory to be provided as its first argument:
# ./gopath_from_workspace.sh <tmpdir>
#
# This populates the directory in $1 as a GOPATH.

set -euo pipefail

bazel fetch //:genfiles_deps

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")
TMP_GOPATH="$1"
BAZEL_EXTERNAL_DIR="$(bazel info output_base)/external"
QUERY_RESULT="$(bazel query 'kind(go_repository, //external:*)' --output xml)"
COUNT="$(echo "${QUERY_RESULT}" | xmllint --xpath "count(//query/rule)" - 2>/dev/null)"

# Copy the dependent repos
for i in $(seq 1 "${COUNT}"); do
  NAME="$(echo "${QUERY_RESULT}" | xmllint --xpath "string(//query/rule[$i]/string[@name='name']/@value)" - 2>/dev/null)"
  DIR="$(echo "${QUERY_RESULT}" | xmllint --xpath "string(//query/rule[$i]/string[@name='importpath']/@value)" - 2>/dev/null)"
  SRC_DIR="${BAZEL_EXTERNAL_DIR}/${NAME}"
  if [ ! -d "${SRC_DIR}" ]; then
    continue
  fi
  DST_DIR="${TMP_GOPATH}/src/${DIR}"
  mkdir -p "${DST_DIR}"
  cp -rT "${SRC_DIR}" "${DST_DIR}"
done

# Copy the cluster-registry repo
mkdir -p "$TMP_GOPATH/src/k8s.io/cluster-registry"
cp -r "${SCRIPT_ROOT}/../"* "$TMP_GOPATH/src/k8s.io/cluster-registry"
