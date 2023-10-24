param location string = resourceGroup().location
param clusterName string = 'aks1'
param nodeCount int = 1
param vmSize string = 'standard_d2s_v5'
param gpu1VmSize string = 'Standard_NC4as_T4_v3'
param gpu2VmSize string = 'Standard_NC6s_v3'
//param gpu3VmSize string = 'Standard_NC24ads_A100_v4'

param kubernetesVersion string = '1.27.3'

var rand = substring(uniqueString(resourceGroup().id), 0, 6)

resource managedIdentity 'Microsoft.ManagedIdentity/userAssignedIdentities@2018-11-30' = {
  name: '${resourceGroup().name}-identity'
  location: location
}

resource aks 'Microsoft.ContainerService/managedClusters@2023-07-02-preview' = {
  name: clusterName
  location: location
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${managedIdentity.id}': {}
    }
  }
  properties: {
    kubernetesVersion: kubernetesVersion
    dnsPrefix: clusterName
    enableRBAC: true
    agentPoolProfiles: [
      {
        name: 'nodepool1'
        count: nodeCount
        vmSize: vmSize
        mode: 'System'
        osType: 'Linux'
        osSKU: 'AzureLinux'
      }
      {
        name: 'gpu1'
        count: 1
        vmSize: gpu1VmSize
        mode: 'User'
        osType: 'Linux'
        osSKU: 'AzureLinux'
        nodeTaints: [
          'sku=gpu:NoSchedule'
        ]
        enableAutoScaling: true
        minCount: 0
        maxCount: 1
      }
      {
        name: 'gpu2'
        count: 1
        vmSize: gpu2VmSize
        mode: 'User'
        osType: 'Linux'
        osSKU: 'AzureLinux'
        nodeTaints: [
          'sku=gpu:NoSchedule'
        ]
        enableAutoScaling: true
        minCount: 0
        maxCount: 1
      }
    ]
  }
}
