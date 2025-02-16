#!/usr/bin/env bash

# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -e
root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )"

echo "> Adding Apache License header to all go files where it is not present"

temp_file=$(mktemp)
trap "rm -f $temp_file" EXIT
sed 's|^// *||' ${root_dir}/hack/LICENSE_BOILERPLATE.txt > $temp_file

go tool -modfile=${root_dir}/go.mod addlicense \
  -f $temp_file \
  -y "$(date +"%Y")" \
  -l apache \
  -ignore ".idea/**" \
  -ignore ".vscode/**" \
  -ignore "dev/**" \
  -ignore "**/*.md" \
  -ignore "**/*.html" \
  -ignore "**/*.yaml" \
  -ignore "**/Dockerfile" \
  -ignore "pkg/component/**/*.sh" \
  -ignore "third_party/gopkg.in/yaml.v2/**" \
  .
