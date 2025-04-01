#!/bin/bash
# Variables – customize these as needed.
RESOURCE_GROUP="vigilantRG"
LOCATION="eastus"
# ACR name must be globally unique; appending a timestamp for uniqueness.
ACR_NAME="vigilantacr$(date +%s)"
CONTAINER_INSTANCE_NAME="vigilant-instance"
IMAGE_NAME="vigilantonion:latest"
# Replace with your repository URL that contains the full solution (app.py, Dockerfile, etc.)
GITHUB_REPO_URL="https://github.com/yourusername/vigilant-web.git"
# Generate a unique DNS label for the container instance
DNS_LABEL="vigilant-instance-$(date +%s)"

echo "Creating resource group: $RESOURCE_GROUP in $LOCATION"
az group create --name $RESOURCE_GROUP --location $LOCATION

echo "Creating Azure Container Registry: $ACR_NAME"
az acr create --resource-group $RESOURCE_GROUP --name $ACR_NAME --sku Basic --location $LOCATION

echo "Enabling admin account for ACR: $ACR_NAME"
az acr update -n $ACR_NAME --admin-enabled true

echo "Cloning repository from GitHub"
git clone $GITHUB_REPO_URL
cd $(basename $GITHUB_REPO_URL .git)

# Optional: Patch requirements.txt if needed (e.g., replace 'beautifulsoup' with 'beautifulsoup4')
if grep -qi "beautifulsoup" requirements.txt; then
    echo "Patching requirements.txt..."
    sed -i 's/beautifulsoup/beautifulsoup4/Ig' requirements.txt
fi

# Verify Dockerfile exists
if [ ! -f Dockerfile ]; then
    echo "Dockerfile not found. Exiting."
    exit 1
fi

echo "Building Docker image in ACR..."
az acr build --registry $ACR_NAME --image $IMAGE_NAME .

# Retrieve ACR credentials
ACR_USERNAME=$(az acr credential show --name $ACR_NAME --query "username" -o tsv)
ACR_PASSWORD=$(az acr credential show --name $ACR_NAME --query "passwords[0].value" -o tsv)
FULL_IMAGE_NAME="$ACR_NAME.azurecr.io/$IMAGE_NAME"

echo "Deploying container instance: $CONTAINER_INSTANCE_NAME"
az container create \
  --resource-group $RESOURCE_GROUP \
  --name $CONTAINER_INSTANCE_NAME \
  --image $FULL_IMAGE_NAME \
  --cpu 1 --memory 1.5 \
  --os-type Linux \
  --registry-login-server "$ACR_NAME.azurecr.io" \
  --registry-username "$ACR_USERNAME" \
  --registry-password "$ACR_PASSWORD" \
  --restart-policy OnFailure \
  --ip-address Public \
  --ports 8080 \
  --dns-name-label $DNS_LABEL

echo "Deployment complete. Your container instance is accessible at:"
az container show --resource-group $RESOURCE_GROUP --name $CONTAINER_INSTANCE_NAME --query ipAddress.fqdn -o tsv
