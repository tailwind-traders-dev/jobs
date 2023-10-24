# Containers

This application can be run locally or on a VM with Go and Mage, or built as a container.

Our `jobs` container is built using the Dockerfiles, [Dockerfile](../Dockerfile) and [dev.Dockerfile](../dev.Dockerfile).

The [build-and-publish.yaml](../.github/workflows/build-and-publish.yaml) GitHub Action builds and publishes the `jobs` container, from the `latest` branch, to GitHub Container Registry.

We support both [GitHub Container Registry](../docs/github.md#github-container-registry) and [Azure Container Registry](../docs/azure-storage.md#azure-container-registry) for public and private images respectively.
