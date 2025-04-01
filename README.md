#Vigilant-Web
A prototype solution that integrates Vigilant Onion’s darkweb scanning functionality with a web interface. This repository provides a Flask web app to issue darkweb scans (e.g., searching for your company name) and retrieve scan results stored in a SQLite database—all deployed as a single container using Azure Container Instances.

#Repository Structure
```
vigilant-web/
├── app.py             # Flask web interface for issuing scans and retrieving results
├── observer.py        # Vigilant Onion scanning script (used via the --find flag)
├── config/
│   └── config.yml     # Configuration file (proxy, timeout, logging, etc.)
├── deploy.sh          # Bash script to build the Docker image, push it to ACR, and deploy to ACI
├── Dockerfile         # Builds the container (installs Tor, Python dependencies, and runs both Tor and Flask)
├── requirements.txt   # Python dependencies (Flask, etc.)
├── .gitignore         # Files/folders to ignore in Git
└── LICENSE            # MIT License
```
##Features
Darkweb Scanning: Uses the observer.py script to search for provided keywords on the darkweb via Tor.

###Web Interface: A simple Flask app (app.py) that allows you to issue scans through a web form and view scan results in JSON format.

###Automated Deployment: The deploy.sh script creates necessary Azure resources (Resource Group, Azure Container Registry) and deploys the container instance with a public FQDN.

###Tor Integration: The Dockerfile installs and starts the Tor client so that all darkweb requests are properly routed.

#Prerequisites
##Azure Cloud Shell: You can run the deployment script directly from the Azure Cloud Shell.

Azure CLI: Ensure you’re logged in (az login).

Git: Available in Cloud Shell.

A valid GitHub repository URL for this project.

#Configuration
The sample configuration file is located at config/config.yml. You can adjust settings such as:
```
Proxy settings (for Tor)

Timeouts

Scoring thresholds

Logging settings (e.g., syslog details)

Example content from config/config.yml:

yaml
Copy
debug: False
dbname: "database.db"
dbpath: "utils/database"
server_proxy: localhost
port_proxy: 9050
type_proxy: socks5h
timeout: 30
score_categorie: 20
score_keywords: 40
count_categories: 5
daystime: 10
sendlog: True
logport: 5151
logip: "localhost"
```
#Deployment
The deploy.sh script automates the entire process:

Creates an Azure resource group.

Creates an Azure Container Registry (ACR) and enables its admin account.

Clones this repository.

Patches dependencies if needed.

Builds the Docker image in ACR.

Deploys the container to Azure Container Instances (ACI) with a DNS name label, providing a public FQDN.

One-Liner Deployment Command
You can deploy the entire solution from Azure Cloud Shell with this one-liner:

```bash
git clone https://github.com/DataGuys/vigilant-web.git && cd vigilant-web && chmod +x deploy.sh && ./deploy.sh
```

After the script runs, it will output the FQDN (e.g., vigilant-instance-<timestamp>.eastus.azurecontainer.io) where you can access the web interface.

Usage
Access the Web Interface:
Open the provided FQDN in your browser. You’ll see a simple form to enter a search term (e.g., your company name) and trigger a darkweb scan.

Issue a Scan:
Enter your search term on the homepage and click “Scan.” The app will invoke the darkweb scan via observer.py and store the results.

View Results:
Visit /results on the deployed FQDN (e.g., http://<fqdn>/results) to see a JSON list of all scans with their results and timestamps.

Notes
Security:
For production, consider adding authentication to your Flask app and using a more robust database.

Logging & Monitoring:
Enhance logging as needed for troubleshooting.

Dependencies:
Ensure that any outdated dependencies (e.g., BeautifulSoup) are updated in requirements.txt if required.

License
This project is licensed under the MIT License.
