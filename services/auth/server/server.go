package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"nexlab.tech/core/pkg/logging"
	"nexlab.tech/core/services/auth/action"
	"nexlab.tech/core/services/auth/env"
	"nexlab.tech/core/services/auth/event"
	"nexlab.tech/core/services/auth/version"
)

func eventsHandler(cfg *initConfig) gin.HandlerFunc {
	ev, err := event.New(&event.Config{
		Controller: cfg.controller,
		Env:        cfg.env,
	})
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		ev.ServeHTTP(c.Writer, c.Request)
	}
}

func actionsHandler(cfg *initConfig) gin.HandlerFunc {
	act, err := action.New(action.Config{
		Env:        cfg.env,
		Controller: cfg.controller,
		JwtAuth:    cfg.JwtAuth,
	})

	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		act.ServeHTTP(c.Writer, c.Request)
	}
}

func verifyTokenPostHandler(cfg *initConfig) gin.HandlerFunc {
	handler := newAuthHandler(cfg)
	return func(c *gin.Context) {
		handler.Post(c.Writer, c.Request)
	}
}

func main() {
	envVar := env.GetEnv()
	logging.InitLogger(envVar.LogLevel)

	logrus.WithField("variables", envVar).Debugf("environment variables")

	cfg, err := NewInitConfig(envVar)
	if err != nil {
		log.Fatalf("failed to initialize service: %s", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/events", eventsHandler(cfg))
	r.POST("/actions", actionsHandler(cfg))
	r.POST("/verify-token", verifyTokenPostHandler(cfg))
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, version.GetVersion())
	})
	log.Fatal(r.Run("0.0.0.0:" + envVar.Port))
}
