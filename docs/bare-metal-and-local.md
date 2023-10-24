# Bare Metal & Local Deployment

The `jobs` application is tailored to provide equal support for the following deployment scenarios:

- Containerized and non-containerized.
- Cloud and non-cloud.
- Local and remote.

This means that by installing Go and mage on a Virtual Machine in the cloud will allow you to run your application with the same experience as you may use to develop it locally.

We do not depend on containers during the development process. We also support non-containerized deployment deployment for production.

Jobs should run as easily on bare metal, taking advantage of a local MacBook running on Apple Silicon, with access to a GPU for AI scenarios, or a Raspberry Pi for less resource-intensive scenarios.

Jobs should run on any platform that runs [Go](https://golang.org/) and [Mage](https://magefile.org/) which includes Linux, macOS, Windows operating systems, across Intel, ARM, and Apple Silicon CPU architectures.
