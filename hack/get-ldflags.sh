#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -o nounset
set -o pipefail
set -o errexit

root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )"
VERSION="${EFFECTIVE_VERSION:-$(cat "${1:-${root_dir}/VERSION}")}"
PROGRAM_NAME="${2:-kube-rbac-proxy-watcher}"


MAJOR_VERSION=""
MINOR_VERSION=""

if [[ "${VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
  MAJOR_VERSION=${BASH_REMATCH[1]}
  MINOR_VERSION=${BASH_REMATCH[2]}
  if [[ -n "${BASH_REMATCH[4]}" ]]; then
    MINOR_VERSION+="+"
  fi
fi

# .dockerignore ignores all files unrelevant for build (e.g. docs) to only copy relevant source files to the build
# container. Hence, git will always detect a dirty work tree when building in a container (many deleted files).
# This command filters out all deleted files that are ignored by .dockerignore to only detect changes to relevant files
# as a dirty work tree.
# Additionally, it filters out changes to the `VERSION` file, as this is currently the only way to inject the
# version-to-build in our pipelines (see https://github.com/gardener/cc-utils/issues/431).
TREE_STATE="$([ -z "$(git status --porcelain 2>/dev/null | grep -vf <(git ls-files -o --deleted --ignored --exclude-from=.dockerignore) -e 'VERSION')" ] && echo clean || echo dirty)"

echo "-X k8s.io/component-base/version.gitMajor=$MAJOR_VERSION
      -X k8s.io/component-base/version.gitMinor=$MINOR_VERSION
      -X k8s.io/component-base/version.gitVersion=$VERSION
      -X k8s.io/component-base/version.gitTreeState=$TREE_STATE
      -X k8s.io/component-base/version.gitCommit=$(git rev-parse --verify HEAD)
      -X k8s.io/component-base/version.buildDate=$(date '+%Y-%m-%dT%H:%M:%S%z' | sed 's/\([0-9][0-9]\)$/:\1/g')
      -X k8s.io/component-base/version/verflag.programName=$PROGRAM_NAME"
