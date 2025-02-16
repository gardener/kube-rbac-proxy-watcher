#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

set -e

root_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )"

gosec_report="false"
gosec_report_parse_flags=""

parse_flags() {
  while test $# -gt 1; do
    case "$1" in
      --gosec-report)
        shift; gosec_report="$1"
        ;;
      *)
        echo "Unknown argument: $1"
        exit 1
        ;;
    esac
    shift
  done
}

parse_flags "$@"
if [[ "$gosec_report" != "false" ]]; then
  echo "Exporting report to $root_dir/gosec-report.sarif"
  gosec_report_parse_flags="-track-suppressions -fmt=sarif -out=gosec-report.sarif -stdout"
fi

go tool -modfile=${root_dir}/go.mod gosec -exclude-generated -exclude-dir=hack $gosec_report_parse_flags ./...