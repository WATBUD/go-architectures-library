# Print message
Write-Output "ps1:Updating Swagger docs..."

# Run Swagger initialization
swag init -g ./cmd/web/main.go

# Check if the command was successful
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to generate Swagger docs."
    exit $LASTEXITCODE
}

# Print message
Write-Output "Building and running application..."

# # Run Go application
# go run ./cmd/web/main.go

# # Check if the command was successful
# if ($LASTEXITCODE -ne 0) {
#     Write-Error "Failed to run the Go application."
#     exit $LASTEXITCODE
# }
