// File: bicep/main.bicep
targetScope = 'subscription'

// === Parameter ===
param location string = 'southeastasia'
param projectPrefix string = 'us'
param principalId string 
@secure()
param sqlAdminPassword string // RECEIVE the password here

// === Variabel ===
var resourceGroupName = 'rg-${projectPrefix}-prod'

// === Resource Group ===
resource rg 'Microsoft.Resources/resourceGroups@2024-03-01' = {
  name: resourceGroupName
  location: location
}

// === Module ===
module coreResources 'core.bicep' = {
  name: 'CoreResourcesDeployment'
  scope: rg 
  params: {
    location: location
    projectPrefix: projectPrefix
    principalId: principalId
    sqlAdminPassword: sqlAdminPassword // PASS it down here
  }
}

// === Outputs ===
output keyVaultName string = coreResources.outputs.keyVaultName
output aksClusterName string = coreResources.outputs.aksClusterName
output sqlServerName string = coreResources.outputs.sqlServerName
output cosmosAccountName string = coreResources.outputs.cosmosAccountName
