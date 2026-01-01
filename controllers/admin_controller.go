package controllers

import (
	"net/http"
	feedclient "trading-bot/internal/client/feed_client"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	FeedClient feedclient.IFeedClient
}

func NewAdminController(feedClient feedclient.IFeedClient) *AdminController {
	return &AdminController{
		FeedClient: feedClient,
	}
}

func (ac *AdminController) StartAuth(c *gin.Context) {
	ctx := c.Request.Context()
	url, err := ac.FeedClient.StartAuth(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Feed Client config error: " + err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, url)
}

func (ac *AdminController) HandleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	err := ac.FeedClient.HandleCallback(ctx, c.Request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while generating session:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful!",
	})
}
