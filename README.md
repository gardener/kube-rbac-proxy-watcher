# kube-rbac-proxy-watcher

[![REUSE status](https://api.reuse.software/badge/github.com/gardener/kube-rbac-proxy-watcher)](https://api.reuse.software/info/github.com/gardener/kube-rbac-proxy-watcher)
[![Build](https://github.com/gardener/kube-rbac-proxy-watcher/actions/workflows/non-release.yaml/badge.svg)](https://github.com/gardener/kube-rbac-proxy-watcher/actions/workflows/non-release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gardener/kube-rbac-proxy-watcher)](https://goreportcard.com/report/github.com/gardener/kube-rbac-proxy-watcher)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](LICENSE) [![Release](https://img.shields.io/github/v/release/gardener/kube-rbac-proxy-watcher.svg?style=flat)](https://github.com/gardener/kube-rbac-proxy-watcher) [![Go Reference](https://pkg.go.dev/badge/github.com/gardener/kube-rbac-proxy-watcher.svg)](https://pkg.go.dev/github.com/gardener/kube-rbac-proxy-watcher)

## Usage

This utility serves the need of managing the lifecycle of a child process in a container environment. It is the container `entrypoint`, which later starts a dependent child process and watches for changes on a particular location on the mounted filesystem. If there are changes the main process restarts the child process by sending SIGTERM signal.

Usually this scenario is beneficial in cases where an application does not support configuration hot reloading and needs to be restarted to load changes from the configuration. Such application is [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy), which needs to be restarted to reflect configuration changes.

## Feedback and Support

Feedback and contributions are always welcome!

Please report bugs or suggestions as [GitHub issues](https://github.com/gardener/kube-rbac-proxy-watcher/issues) or reach out on [Slack](https://gardener-cloud.slack.com/) (join the workspace [here](https://gardener.cloud/community/community-bio/)).

## Learn more

Please find further resources about out project here:

* [Our landing page gardener.cloud](https://gardener.cloud/)
* ["Gardener, the Kubernetes Botanist" blog on kubernetes.io](https://kubernetes.io/blog/2018/05/17/gardener/)
* ["Gardener Project Update" blog on kubernetes.io](https://kubernetes.io/blog/2019/12/02/gardener-project-update/)
* [Gardener Extensions Golang library](https://godoc.org/github.com/gardener/gardener/extensions/pkg)
* [GEP-1 (Gardener Enhancement Proposal) on extensibility](https://github.com/gardener/gardener/blob/master/docs/proposals/01-extensibility.md)
* [Extensibility API documentation](https://github.com/gardener/gardener/tree/master/docs/extensions)
