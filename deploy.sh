#!/bin/bash
set -euo pipefail

# ------------------------------------------------------------------------------
# DEPLOY.SH
# Builds the container from GitHub, pushes it to ACR, then deploys to ACI
# Usage: ./deploy.sh
# ------------------------------------------------------------------------------

echo "Listing available Azure subscriptions..."
# We pipe to 'nl' to number them for easy selection
az account list --query "[].{Name:name, ID:id}" -o table | nl

# Prompt for subscription number
read -p "Enter the subscription number to use: " subNum

# Convert from 1-based to 0-based index
subIndex=$((subNum - 1))
selectedSub=$(az account list --query "[$subIndex].id" -o tsv)
if [ -z "$selectedSub" ]; then
  echo "Error: invalid subscription selection."
  exit 1
fi

# Set the subscription
az account set --subscription "$selectedSub"
echo "Using subscription ID: $selectedSub"

# Prompt for the Resource Group name (create if missing)
read -p "Enter the Resource Group name: " resourceGroup
if ! az group show --name "$resourceGroup" &>/dev/null; then
  echo "Resource Group '$resourceGroup' does not exist, creating..."
  az group create --name "$resourceGroup" --location eastus
else
  echo "Resource Group '$resourceGroup' already exists. Using it."
fi

# Define ACR name (must be lowercase, 5-50 alphanumeric chars)
# If "myregistry" is not unique, consider parameterizing it.
acrName="myregistry"

# Check if ACR already exists
if ! az acr show --name "$acrName" --resource-group "$resourceGroup" &>/dev/null; then
  echo "Creating Azure Container Registry '$acrName' in RG '$resourceGroup'..."
  # Omit --api-version if it's causing issues
  az acr create --resource-group "$resourceGroup" --name "$acrName" --sku Basic
else
  echo "ACR '$acrName' already exists in RG '$resourceGroup'."
fi

# Build the container image from GitHub and push it to ACR
echo "Building Docker image from GitHub and pushing to ACR..."
az acr build \
  --registry "$acrName" \
  --image vigilant-web:latest \
  https://github.com/DataGuys/vigilant-web.git

# Retrieve ACR login server, username, and password
acrLoginServer=$(az acr show --name "$acrName" --resource-group "$resourceGroup" --query "loginServer" -o tsv)
username=$(az acr credential show --name "$acrName" --resource-group "$resourceGroup" --query "username" -o tsv)
password=$(az acr credential show --name "$acrName" --resource-group "$resourceGroup" --query "passwords[0].value" -o tsv)

# Generate a unique DNS label (you can customize this logic)
uniqueSuffix=$(date +%s) # or use random string
dnsLabel="vigilant-web-$uniqueSuffix"

# Deploy the container image to Azure Container Instances with registry credentials
echo "Deploying container to Azure Container Instances..."
az container create \
  --resource-group "$resourceGroup" \
  --name vigilant-web \
  --image "${acrLoginServer}/vigilant-web:latest" \
  --dns-name-label "$dnsLabel" \
  --ports 443 \
  --registry-username "$username" \
  --registry-password "$password"

# Wait for the container to be in a running state
echo "Waiting for the container instance to be in a 'Running' state..."
az container show \
  --resource-group "$resourceGroup" \
  --name vigilant-web \
  --query "instanceView.state" \
  --output table

echo "Deployment complete!"
fqdn=$(az container show --resource-group "$resourceGroup" --name vigilant-web --query "ipAddress.fqdn" -o tsv)
echo "Your container is accessible at: https://$fqdn"
