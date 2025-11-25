// File: bicep/core.bicep

// === Parameter ===
param location string
param projectPrefix string
param principalId string
@secure()
param sqlAdminPassword string

// === Variabel ===
var uniqueSuffix = uniqueString(resourceGroup().id)
var keyVaultName = '${projectPrefix}-kv-${uniqueSuffix}'
var aksClusterName = '${projectPrefix}-aks-${uniqueSuffix}'
var sqlServerName = '${projectPrefix}-sql-${uniqueSuffix}'
var cosmosAccountName = '${projectPrefix}-cosmos-${uniqueSuffix}'
var serviceBusNamespaceName = '${projectPrefix}-sb-${uniqueSuffix}'

// === 1. Azure Key Vault ===
resource kv 'Microsoft.KeyVault/vaults@2023-07-01' = {
  name: keyVaultName
  location: location
  properties: {
    sku: { family: 'A', name: 'standard' }
    tenantId: subscription().tenantId
    accessPolicies: [
      {
        tenantId: subscription().tenantId
        objectId: principalId
        permissions: { secrets: [ 'get', 'list', 'set', 'delete' ] }
      }
    ]
  }
}

// === 2. Azure Kubernetes Service (AKS) ===
resource aks 'Microsoft.ContainerService/managedClusters@2024-01-01' = {
  name: aksClusterName
  location: location
  identity: { type: 'SystemAssigned' }
  properties: {
    dnsPrefix: '${projectPrefix}-aks'
    agentPoolProfiles: [
      {
        name: 'agentpool'
        count: 2
        vmSize: 'Standard_B2s' 
        mode: 'System'
        osType: 'Linux'
      }
    ]
    networkProfile: { networkPlugin: 'azure', loadBalancerSku: 'standard' }
  }
}

// === 3. Azure SQL Server & Database (Relational) ===
// CHANGED VERSION TO 2021-11-01 (Stable)
resource sqlServer 'Microsoft.Sql/servers@2021-11-01' = {
  name: sqlServerName
  location: location
  properties: {
    administratorLogin: 'sqladmin'
    administratorLoginPassword: sqlAdminPassword
  }
}

// CHANGED VERSION TO 2021-11-01 (Stable)
resource sqlDb 'Microsoft.Sql/servers/databases@2021-11-01' = {
  parent: sqlServer
  name: 'links-db'
  location: location
  sku: { name: 'Basic', tier: 'Basic' } 
}

// CHANGED VERSION TO 2021-11-01 (Stable)
resource sqlFirewall 'Microsoft.Sql/servers/firewallRules@2021-11-01' = {
  parent: sqlServer
  name: 'AllowAzureServices'
  properties: {
    startIpAddress: '0.0.0.0'
    endIpAddress: '0.0.0.0'
  }
}

// === 4. Azure Cosmos DB (NoSQL Analytics) ===
resource cosmosAccount 'Microsoft.DocumentDB/databaseAccounts@2024-05-15' = {
  name: cosmosAccountName
  location: location
  kind: 'GlobalDocumentDB'
  properties: {
    databaseAccountOfferType: 'Standard'
    locations: [ { locationName: location, failoverPriority: 0 } ]
    capabilities: [ { name: 'EnableServerless' } ] 
  }
}

resource cosmosDb 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2024-05-15' = {
  parent: cosmosAccount
  name: 'analytics-db'
  properties: { resource: { id: 'analytics-db' } }
}

resource cosmosContainer 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2024-05-15' = {
  parent: cosmosDb
  name: 'clicks'
  properties: {
    resource: {
      id: 'clicks'
      partitionKey: { paths: [ '/short_code' ], kind: 'Hash' }
    }
  }
}

// === 5. Azure Service Bus (Messaging) ===
resource serviceBus 'Microsoft.ServiceBus/namespaces@2022-10-01-preview' = {
  name: serviceBusNamespaceName
  location: location
  sku: { name: 'Basic' }
}

resource sbQueue 'Microsoft.ServiceBus/namespaces/queues@2022-10-01-preview' = {
  parent: serviceBus
  name: 'analytics-queue'
}

// === Outputs ===
output keyVaultName string = kv.name
output aksClusterName string = aks.name
output sqlServerName string = sqlServer.name
output cosmosAccountName string = cosmosAccount.name

