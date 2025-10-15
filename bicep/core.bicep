// File ini berisi sumber daya yang akan di-deploy di dalam Resource Group.

// === Parameter ===
// Nilai-nilai ini diterima dari file 'main.bicep'.
param location string
param projectPrefix string
param principalId string

// === Variabel ===
var containerRegistryName = '${projectPrefix}acr${uniqueString(resourceGroup().id)}'
var keyVaultName = '${projectPrefix}-kv-${uniqueString(resourceGroup().id)}'

// === Azure Container Registry (ACR) ===
// Tempat untuk menyimpan Docker image Anda.
resource acr 'Microsoft.ContainerRegistry/registries@2023-07-01' = {
  name: containerRegistryName
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    adminUserEnabled: true
  }
}

// === Azure Key Vault ===
// Brankas aman untuk menyimpan semua secret aplikasi.
resource kv 'Microsoft.KeyVault/vaults@2023-07-01' = {
  name: keyVaultName
  location: location
  properties: {
    sku: {
      family: 'A'
      name: 'standard'
    }
    tenantId: subscription().tenantId
    accessPolicies: [
      {
        tenantId: subscription().tenantId
        objectId: principalId
        permissions: {
          secrets: [ 'get', 'list', 'set', 'delete', 'purge', 'recover' ]
        }
      }
    ]
  }
}

// === Outputs ===
// Mengirimkan nilai kembali ke 'main.bicep' untuk ditampilkan.
output acrLoginServer string = acr.properties.loginServer
output keyVaultName string = kv.name
