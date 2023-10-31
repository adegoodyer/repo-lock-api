# Repo lock API

## Overview
- simple Golang API for locking/unlocking pipelines during execution
- routing via [Gin](https://github.com/gin-gonic/gin)

### Todo
- add tests
- containerize
- config (port, set Gin mode etc)

## API
```bash
GET /pipeline/status
GET /pipeline/status/repoName
POST /pipeline/lock
POST /pipeline/unlock
POST /pipeline/unlockAll
```

## Curl commands
```bash
# check status of all repos
curl http://localhost:8081/pipeline/status
# {"repo_example":false,"repo_example2":true}

# check status of single repo
curl http://localhost:8081/pipeline/status/example_repo
# {"repo_example":false}

# lock a repo
curl -X POST -H "Content-Type: application/json" -d '{"repo_name": "example_repo"}' http://localhost:8081/pipeline/lock
# {"message":"Repository locked successfully"}

# unlock a repo
curl -X POST -H "Content-Type: application/json" -d '{"repo_name": "example_repo"}' http://localhost:8081/pipeline/unlock
# {"message":"Repository unlocked successfully"}

# unlock all repos
curl -X POST http://localhost:8081/pipeline/unlockAll
# {"message":"All repositories unlocked successfully"}
```

## Pipeline Usage

### Check lock
```bash
repoName="exampleRepoName"

# Function to check repository lock status
check_repo_lock() {
    local repoName="$1"

    # GET request to API endpoint to check repository status
    response=$(curl -s http://localhost:8081/pipeline/status/$repoName)

    # Extract locked status from response
    lockedStatus=$(echo $response | jq -r '.locked')

    if [ "$lockedStatus" == "true" ]; then
        # Repository is locked, halt pipeline
        echo "Error: Repository $repoName is locked. Halting pipeline."
        exit 1
    else
        # Repository isn't locked, continue with pipeline
        echo "Repository $repoName is unlocked. Continuing with pipeline."
    fi
}

# Call function to check repository status
check_repo_lock "$repoName"

# Rest of pipeline logic...
```

### Add lock
```bash
repoName="exampleRepoName"

# Lock repository at start of pipeline
lock_repo() {
    echo "Locking repository $repoName..."
    curl -X POST -H "Content-Type: application/json" -d "{\"repoName\": \"$repoName\"}" http://localhost:8081/pipeline/lock
    echo "Repository $repoName locked successfully."
}

# Call function to lock repository at start of pipeline
lock_repo

# Rest of pipeline logic...
```

### Remove lock
```bash
repoName="exampleRepoName"

# Unlock repository at end of pipeline
unlock_repo() {
    echo "Unlocking repository $repoName..."
    curl -X POST -H "Content-Type: application/json" -d "{\"repoName\": \"$repoName\"}" http://localhost:8081/pipeline/unlock
    echo "Repository $repoName unlocked successfully."
}

# Rest of pipeline logic...

# Call function to unlock repository at end of pipeline
unlock_repo
```
