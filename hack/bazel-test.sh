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

# Generates all necessary files and runs bazel test against the targets in
# the repository.
#
# This is meant to be used from a Kubernetes test-infra execute scenario, as
# defined here:
# https://github.com/kubernetes/test-infra/blob/master/scenarios/execute.py

set -euo pipefail

SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "${SCRIPT_ROOT}/.."
hack/update-codegen.sh
bazel run //:gazelle
# TODO: This is brittle. Find a better way to reference the test-infra repo.
/workspace/test-infra/scenarios/kubernetes_bazel.py --test="//... -//pkg/client/... -//cmd/clusterregistry:push-clusterregistry-image"
