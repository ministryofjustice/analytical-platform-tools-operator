<p align="center">
  <h2 align="center">Analytical Platform Tools Operator</h2>
  <p align="center">This repository contains the necessary code to deploy a "Tools" Kubernetes operator into the Analytical Platform ecosystem.</p>
</p>

---

[![Go Report Card](https://goreportcard.com/badge/github.com/ministryofjustice/analytical-platform-tools-operator)](https://goreportcard.com/report/github.com/ministryofjustice/analytical-platform-tools-operator)
[![GoDoc](https://godoc.org/github.com/ministryofjustice/analytical-platform-tools-operator?status.svg)](https://godoc.org/github.com/ministryofjustice/analytical-platform-tools-operator)
[![GitHub release](https://img.shields.io/github/release/ministryofjustice/analytical-platform-tools-operator.svg)](https://GitHub.com/ministryofjustice/analytical-platform-tools-operator/releases/)

The "Tools" operator will enable the [control panel](https://github.com/ministryofjustice/analytics-platform-control-panel/tree/main/controlpanel) to communicate with it using standard REST API calls. It's function will be to list, create and delete tools required by the MoJ Data Analysts. The tools at the time of writing this were limited to:

- [JupyterLab](https://jupyter.org/)

- [Airflow](https://airflow.apache.org/)

- [Rstudio](https://www.rstudio.com/)

Using a [Kubernetes client](https://kubernetes.io/docs/reference/using-api/client-libraries/) we can do things like:

- Show what the current user has deployed in their namespace:

  ```bash
  > kubectl get tool -n user-namespace

  NAME         AGE
  airflow      17h
  jupyterlab   17h
  rstudio      17h
  ```

- Install/start/stop and delete a tool for the current user.

  ```bash
  > kubectl apply -f ./config/samples/tools_v1alpha1_jupyterlab.yaml
  ```

  An example of a manifest used to create JupyterLab would be:

  ```bash
  apiVersion: tools.analytical-platform.justice/v1alpha1
  kind: JupyterLab
  metadata:
    name: jupyterlab-sample
    namespace: user-namespace
  spec:
    image: jupyter/minimal-notebook
    version: python-3.9.10
  ```

- List all available versions of each tools

  TODO: Add global tool command - perhaps another api

This operator is a small step towards easing development of features and further tooling.

## Development practices

#### Commit messages

Please if you can use [conventional commits](https://conventionalcommits.org/) to make your commits and follow the [conventional commit guidelines](https://www.conventionalcommits.org/en/v1.0.0/#specification).

#### Versioning

Please use [semantic versioning](https://semver.org/) for versioning when releasing. This will make it easier to track changes and to make it easier to find the latest version.

## Local development

#### Pre-reqs

- install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

#### How to use

Using a local Kubernetes cluster, you should be able to deploy and perform actions locally that should emulate the dev and production clusters.

To spin up a local cluster with a local container registry run the following in the root of the repository:

```bash
make cluster up
```

To run the controller outside the cluster, run:

```bash
make install run
```

In another window you can create a sample resource from here: `./config/samples/.`

## Deployment

The deployment in this repository isn't ideal. At the time of its creation, the Analytical Platform only allows installs via [flux](https://github.com/moj-analytical-services/analytical-platform-flux/tree/main/clusters/development/apps/tools-operator) - this flux installation is also fairly flaky and required a bit of hacking to allow us control of the image it installs. This next section should clear up how deployment works.

#### Development cluster

The development cluster is an EKS cluster in the dev AWS account. Deployment to this cluster should be restricted to anything in the main branch of this repository.

##### GitHub Action

A GitHub Acton named `build-test-build-dev` triggers on each push to main. The intention of this pipeline is to test (`make test`), build (`make docker-build`) and push to dockerhub (`make docker-push`) using a combination of `branch`-`gitSHA`-`timestamp` as the image tag. The pipeline will then deploy the controller to the Development cluster by amending the image tag in `config/manager`.

The less than ideal part of this pipeline is we have to amend the kustomize manifest file as a step in the pipeline. This means the pipeline creates a commit and push to main to deploy.

##### Manual

###### Pre-reqs

You must have push permissions to the `ministryofjustice` dockerhub repository.

###### Build, push and deploy

This will test and build locally, pushing the [ministryofjustice/analytical-platform-tools-operator](https://hub.docker.com/r/ministryofjustice/analytical-platform-tools-operator) dockerhub repository.

```bash
make docker-build docker-push deploy
```

#### Production cluster

The intention of the "Deploy" to production stage is to only push on release. So all changes to main get pushed to dev, but only release get pushed to prod.

## Test

The code in this repository uses [envtest](https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html) to perform go tests against mock ectd and api components. To run these tests:

```bash
make test
```

## How do I add another API

This operator was bootstrapped using the [operator-sdk](https://sdk.operatorframework.io/docs/overview/) binary. To add another API, controller and custom resource:

```bash
> brew install operator-sdk
> operator-sdk create api --group tools --version v1alpha1 --kind NEW_RESOURCE_NAME --resource --controller
> make generate manifests
```

A decent tutorial can be found [here](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)
