# PowerShell script to run Postman tests
param(
    [Parameter(Mandatory=$true)]
    [string]$BaseUrl,
    
    [Parameter(Mandatory=$true)]
    [string]$ApiKey
)

# Exit on any error
$ErrorActionPreference = "Stop"

Write-Host "Checking if newman is installed..."
try {
    newman --version | Out-Null
} catch {
    Write-Host "Installing newman..."
    npm install -g newman
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to install newman"
        exit 1
    }
}

Write-Host "Running Postman tests against $BaseUrl..."

# Set environment variables
$env:CI_BASE_URL = $BaseUrl
$env:CI_API_KEY = $ApiKey

# Run the tests
newman run ".\postman\RemindMeBot.postman_collection.json" `
    -e ".\postman\RemindMeBot.ci.postman_environment.json" `
    --env-var "CI_BASE_URL=$BaseUrl" `
    --env-var "CI_API_KEY=$ApiKey" `
    --reporters 'cli,junit' `
    --reporter-junit-export "test-results.xml" `
    --bail

if ($LASTEXITCODE -ne 0) {
    Write-Error "Postman tests failed"
    exit $LASTEXITCODE
}

Write-Host "All tests passed successfully!"