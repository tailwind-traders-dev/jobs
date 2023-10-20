param principalID string

var rand = substring(uniqueString(resourceGroup().id), 0, 6)

resource serviceBus 'Microsoft.ServiceBus/namespaces@2021-06-01-preview' existing = {
  name: 'servicebus${rand}'
}

var roleDefinitionId = {
  Owner: '8e3af657-a8ff-443c-a75c-2fe8c4bcb635'
  Contributor: 'b24988ac-6180-42a0-ab88-20f7382dd24c'
  Reader: 'acdd72a7-3385-48ef-bd42-f606fba81ae7'
  KeyVaultAdministrator: '00482a5a-887f-4fb3-b363-3b7fe8e74483'
  AcrPull: '7f951dda-4ed3-4680-a7ca-43fe172d538d'
  ServiceBusDataOwner: '090c5cfd-751d-490a-894a-3ce6f1109419'
  StorageBlobDataContributor: 'ba92f5b4-2d11-453d-a403-e96b0029c9fe'
  KubernetesServiceClusterUserRole: '4abbcc35-e782-43d8-92c5-2d3f1bd2253f'
}

var roleAssignmentServiceBusDefinition = 'ServiceBusDataOwner'
resource roleAssignmentServiceBus 'Microsoft.Authorization/roleAssignments@2020-08-01-preview' = {
  name: guid(serviceBus.id, principalID)
  scope: serviceBus
  properties: {
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', roleDefinitionId[roleAssignmentServiceBusDefinition])
    principalId: principalID
  }
}
