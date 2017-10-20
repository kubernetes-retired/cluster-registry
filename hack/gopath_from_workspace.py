#!/usr/bin/env python

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
# Requires an empty temporary directory to be provided as its first argument:
# ./gopath_from_workspace.sh <tmpdir>
#
# This populates the provided directory as a GOPATH.

import os.path
import argparse
import shutil
import string
import subprocess
import xml.etree.ElementTree as ElementTree


def main(tmpdir):
  subprocess.check_call(["bazel", "fetch", "//:genfiles_deps"])

  bazel_external_dir = os.path.join(
      string.strip(subprocess.check_output(["bazel", "info", "output_base"])),
      "external")
  workspace_dir = string.strip(
      subprocess.check_output(["bazel", "info", "workspace"]))

  query_result = subprocess.check_output([
      "bazel", "query", "kind(go_repository, //external:*)", "--output", "xml"
  ])
  xml = ElementTree.fromstring(query_result)
  elements = xml.findall("./rule")
  for e in elements:
    name = e.find("./string[@name='name']").attrib["value"]
    importpath_element = e.find("./string[@name='importpath']")
    if importpath_element is not None:
      import_path = importpath_element.attrib["value"]
      srcdir = os.path.join(bazel_external_dir, name)
      if os.path.exists(srcdir):
        shutil.copytree(
            srcdir, os.path.join(tmpdir, "src", import_path), symlinks=True)

  shutil.copytree(
      workspace_dir,
      os.path.join(tmpdir, "src", "k8s.io", "cluster-registry"),
      symlinks=True)


if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument("tmpdir")
  args = parser.parse_args()
  main(string.strip(args.tmpdir))
