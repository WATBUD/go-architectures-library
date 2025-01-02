

# Create a Folder
curl -X POST http://localhost:8080/folders ^
     -H "Content-Type: application/json" ^
     -d "{\"name\": \"Folder1\", \"color\": \"blue\", \"index\": 1, \"parentIndex\": 0}"

# Get All Folders
curl -X GET http://localhost:8080/folders

# Get a Specific Folder
curl -X GET http://localhost:8080/folders/{id}

# Replace {id} with the actual folder ID, for example:
curl -X GET http://localhost:8080/folders/60c72b2f9f1b2c6d88f5a123

# Update a Folder
curl -X PUT http://localhost:8080/folders/{id} ^
     -H "Content-Type: application/json" ^
     -d "{\"name\": \"UpdatedFolder\", \"color\": \"red\", \"index\": 2, \"parentIndex\": 1}"

# Delete a Folder
curl -X DELETE http://localhost:8080/folders/{id}

# Update Folder Index
curl -X PUT http://localhost:8080/folders/{id}/index ^
     -H "Content-Type: application/json" ^
     -d "{\"index\": 3}"

# Update Folder Parent
curl -X PUT http://localhost:8080/folders/{id}/index ^
     -H "Content-Type: application/json" ^
     -d "{\"index\": 3}"

# Offline Read
curl -X GET http://localhost:8080/sync/offline







