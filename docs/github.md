# GitHub

In this project we use GitHub to host our code, and GitHub Actions to build our applications, often into container images, and GitHub Packages [Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry) to host public versions of our container images.

Our GitHub Actions workflow, [.github/workflows/build-and-publish.yaml](../.github/workflows/build-and-publish.yaml), will be triggered on any commit to the `latest` branch, and create a container image tagged `latest` (e.g. `ghcr.io/tailwind-traders-dev/jobs:latest`)under the "Packages" tab of our GitHub organization or user.

We use the `latest` image primarily for development purposes.

Sometimes, we may want to push another branch to `latest` in order to test it prior to commiting that code to main. We can do this via `git push origin HEAD:latest`.

We can also use a [workflow_dispatch](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch) trigger to build a container image from a specific branch. We intend to expand upon this workflow, as well as enable proper support for verion tags in our GitHub Actions workflow.

