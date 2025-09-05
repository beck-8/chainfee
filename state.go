package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func stateHeight(c *gin.Context) {
	head, err := lapi.ChainHead(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		})
		return
	}
	c.String(http.StatusOK, head.Height().String())

}
