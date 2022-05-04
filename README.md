<p align="center">
  <h2 align="center">Analytical Platform Tools Operator</h2>
  <p align="center">The intention of this repository is to act as a placeholder for the Tools Operator.</p>
</p>

---

[![Go Report Card](https://goreportcard.com/badge/github.com/ministryofjustice/analytical-platform-tools-operator)](https://goreportcard.com/report/github.com/ministryofjustice/analytical-platform-tools-operator)
[![GoDoc](https://godoc.org/github.com/ministryofjustice/analytical-platform-tools-operator?status.svg)](https://godoc.org/github.com/ministryofjustice/analytical-platform-tools-operator)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/ministryofjustice/analytical-platform-tools-operator.svg)](https://github.com/ministryofjustice/analytical-platform-tools-operator)
[![GitHub issues](https://img.shields.io/github/issues/ministryofjustice/analytical-platform-tools-operator.svg)](https://GitHub.com/ministryofjustice/analytical-platform-tools-operator/issues/)
[![GitHub release](https://img.shields.io/github/release/ministryofjustice/analytical-platform-tools-operator.svg)](https://GitHub.com/ministryofjustice/analytical-platform-tools-operator/releases/)

## Development practices

### Commit messages

Please if you can use [conventional commits](https://conventionalcommits.org/) to make your commits and follow the [conventional commit guidelines](https://conventionalcommits.org/en/v1.0.0/guidelines.html).

### Versioning

Please use [semantic versioning](https://semver.org/) for versioning when releasing. This will make it easier to track changes and to make it easier to find the latest version.

## Local development

### Pre-reqs

- you should install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

### How to use

Using a local Kubernetes cluster, you should be able to deploy and perform actions locally that should emulate the dev and production clusters.

To spin up a local cluster with a local container registry run the following in the root of the repository:

```bash
make cluster up
```

When you're ready you can build and tag for local push:

```bash
make local-docker-build
```

And then push to the local image registry hosted in the kind cluster:

```bash
make local-docker-push
```

Finally deploy:

```bash
make deploy
```

## GitHub actions

The following actions run in this repository and perform the following:

### Tests

`make test` is run on each push to a branch.
