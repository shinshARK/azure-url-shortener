// File ini dijalankan di level Subscription untuk membuat Resource Group.
targetScope = 'subscription'

// === Parameter ===
param location string = 'southeastasia'
param projectPrefix string = 'us'
param principalId string // ID pengguna yang menjalankan deployment.

// === Variabel ===
var resourceGroupName = 'rg-${projectPrefix}-prod'

// === Resource Group ===
// Membuat Resource Group utama untuk proyek ini.
resource rg 'Microsoft.Resources/resourceGroups@2024-03-01' = {
  name: resourceGroupName
  location: location
}

// === Module ===
// Mendeploy file 'core.bicep' ke dalam Resource Group yang baru dibuat.
module coreResources 'core.bicep' = {
  name: 'CoreResourcesDeployment'
  scope: rg // Menentukan target deployment adalah Resource Group 'rg'.
  params: {
    location: location
    projectPrefix: projectPrefix
    principalId: principalId
  }
}

// === Outputs ===
// Menampilkan informasi penting setelah deployment selesai.
output keyVaultName string = coreResources.outputs.keyVaultName
output aksClusterName string = coreResources.outputs.aksClusterName
output aksResourceId string = coreResources.outputs.aksResourceId