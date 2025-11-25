// File ini berisi sumber daya yang akan di-deploy di dalam Resource Group.

// === Parameter ===
// Nilai-nilai ini diterima dari file 'main.bicep'.
param location string
param projectPrefix string
param principalId string

// === Variabel ===
var keyVaultName = '${projectPrefix}-kv-${uniqueString(resourceGroup().id)}'
var aksClusterName = '${projectPrefix}-aks-${uniqueString(resourceGroup().id)}'

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

// === Azure Kubernetes Service (AKS) ===
// Cluster Kubernetes untuk menjalankan container aplikasi.
resource aks 'Microsoft.ContainerService/managedClusters@2024-01-01' = {
  name: aksClusterName
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    dnsPrefix: '${projectPrefix}-aks'
    agentPoolProfiles: [
      {
        name: 'agentpool'
        count: 2
        vmSize: 'Standard_B2s'  // Small VM for demo/dev
        mode: 'System'
        osType: 'Linux'
      }
    ]
    networkProfile: {
      networkPlugin: 'azure'
      loadBalancerSku: 'standard'
    }
  }
}

// === Outputs ===
// Mengirimkan nilai kembali ke 'main.bicep' untuk ditampilkan.
output keyVaultName string = kv.name
output aksClusterName string = aks.name
output aksResourceId string = aks.id