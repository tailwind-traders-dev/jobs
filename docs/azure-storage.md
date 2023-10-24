# Azure Storage & Data

In addition to deploying one or more of our 3 compute services, Virtual Machines, Container Apps, and Kubernetes Service, we deploy 5 "core" Storage & Data services in Azure. We deploy the first four alongside every deployment, and keep the fifth optional.

We use the term "Storage" in contrast to "Compute" to refer to services that are primarily used to store data, and not to run code. However, they are not limited to a particular Storage service, such as Azure Blob Storage, but persist various types of data and state, from blobs, to secrets, to messages in-flight, container images, and relational data in PostgreSQL.

The hallmark of these core services is that they are fully managed and cost effective. There is often a one-to-many relationship between the service, and the number of apps it serves. For example, we have one Azure Container Registry, but many apps that use it to store container images.

## Azure Blob Storage

Azure Blob Storage is a very frequently used service across all apps to store blobs of all kinds.

## Azure Key Vault

While we primarily use Managed Identity to access services in Azure, including Key Vault, there are often various secrets in our apps that we need our apps to be able to consume efficiently.

## Azure Service Bus

Azure Service Bus is a key service used to enable asynchronous communication between apps. It helps us decouple our applications and effectively scale the compute.

## Azure Container Registry

The majority of our deployed applications will be containerized at some point. This enables us to deploy them consistently across environments, from Azure Kubernetes Service and Azure Container Apps, to Virtual Machines using Docker and/or Docker Compose.

Azure Container Registry [Tasks](https://learn.microsoft.com/en-us/azure/container-registry/container-registry-tasks-overview#quick-task)  allows us to avoid any dependency on a local environment for building container images, allowing us to push our code and Dockerfile to the registry using the Azure CLI (az), and have the image built and stored in the registry.

Azure Container Registry images are private by default, and secured by Role-Based Access Control (RBAC) and Managed Identity.

By comparison, GitHub Container Registry makes it easy to distribute images publicly, and is a great option for open source projects and containers that contain no sensitive code or data and accessible without authentication. While we can use GitHub Container Registry as a private registry, this "only supports authentication using a personal access token (classic)", and we prefer Managed Identity for production workloads.

## Azure Database for Postgres

Postgres is our relational database of choice for our open source applications.

However, because this service carries a cost due to persistent compute and storage, we make this component optional. We do this as some scenarios may not require Postgres up-front, or may utilize an alternative deployment option, such as running a container on Azure Kubernetes Service, the Postgres add-on for Azure Container Apps, or running Postgres in a container on a Virtual Machine.

We optionally deploy an Azure Database for Postgres (Flexible Server) instance.
