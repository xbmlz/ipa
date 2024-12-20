package ipa

import (
	"github.com/gin-gonic/gin"
	"github.com/xbmlz/ipa/logger"
)

type Context struct {
	// gin context
	*gin.Context
	// logger
	Log logger.Logger
}
