package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Repository struct {
	Locked bool `json:"locked"`
}

var repositories = make(map[string]*Repository)
var mu sync.Mutex

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/pipeline/status", getAllRepos)
	r.GET("/pipeline/status/:repoName", getRepo)
	r.POST("/pipeline/lock", lockRepo)
	r.POST("/pipeline/unlock", unlockRepo)
	r.POST("/pipeline/unlockAll", unlockAllRepos)

	r.Run(":8081")
}

type LockRequest struct {
	RepoName string `json:"repoName"`
}

func lockRepo(c *gin.Context) {
	var request LockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	repo, ok := repositories[request.RepoName]
	if !ok {
		repo = &Repository{}
		repositories[request.RepoName] = repo
	}
	repo.Locked = true

	c.JSON(http.StatusOK, gin.H{"message": "Repository locked successfully"})
}

func unlockRepo(c *gin.Context) {
	var request LockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	repo, ok := repositories[request.RepoName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}
	repo.Locked = false

	c.JSON(http.StatusOK, gin.H{"message": "Repository unlocked successfully"})
}

func unlockAllRepos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	for _, repo := range repositories {
		repo.Locked = false
	}

	c.JSON(http.StatusOK, gin.H{"message": "All repositories unlocked successfully"})
}

func getAllRepos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	response := make(map[string]bool)
	for repoName, repo := range repositories {
		response[repoName] = repo.Locked
	}
	c.JSON(http.StatusOK, response)
}

func getRepo(c *gin.Context) {
	repoName := c.Param("repoName")
	mu.Lock()
	defer mu.Unlock()

	repo, ok := repositories[repoName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{repoName: repo.Locked})
}
