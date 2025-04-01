#!/bin/bash
# List available subscriptions with numbers
echo "Available Subscriptions:"
az account list --query "[].{Name:name, ID:id}" -o table | nl

# Prompt for subscription number
read -p "Enter the subscription number to use: " subNum

# Get the selected subscription ID (note: converting the number to zero-index)
selectedSub=$(az account list --query "[$((subNum-1))].id" -o tsv)

# Set the subscription
az account set --subscription "$selectedSub"
echo "Using subscription: $selectedSub"

# Prompt for the resource group name
read -p "Enter the Resource Group name: " resourceGroup

# Define ACR name (adjust as needed)
acrName="myRegistry"

# Create the Azure Container Registry
az acr create --resource-group "$resourceGroup" --name "$acrName" --sku Basic

# Build the container image from GitHub and push it to ACR
az acr build --registry "$acrName" --image vigilant-web:latest https://github.com/DataGuys/vigilant-web.git

# Deploy the container image to an Azure Container Instance
az container create \
  --resource-group "$resourceGroup" \
  --name vigilant-web \
  --image "${acrName}.azurecr.io/vigilant-web:latest" \
  --dns-name-label vigilant-web-unique \
  --ports 443
