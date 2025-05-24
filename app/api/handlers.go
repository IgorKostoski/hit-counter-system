package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HitRequest struct {
	Key string `json:"key" binding:"required"`
}

func HitHandler(c *gin.Context) {
	var req HitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	if req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key cannot be empty "})
		return
	}

	count, err := IncrementHit(db, req.Key)
	if err != nil {
		hitsProcessed.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment hit count" + err.Error()})
		return
	}

	hitsProcessed.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{"key": req.Key, "new_count": count})
}

func CountHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key parameter is required"})
		return
	}

	count, err := GetCount(db, key)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve count: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "count": count})
}
