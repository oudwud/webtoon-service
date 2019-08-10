package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oudwud/webtoon-service/pkg/config"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func Run(conf *config.Config) error {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.MaxMultipartMemory = 32 << 20 // 32 MiB

	v1 := router.Group("/v1")
	{
		v1.POST("/merge", mergeImagesHandler)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http listen error: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Warn("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "fail to shutdown gracefully")
	}

	log.Warn("server shutdown successfully")
	return nil
}

func writeErrorResp(c *gin.Context, err error) {
	if httpErr, ok := err.(*httpError); ok {
		c.String(httpErr.Code(), httpErr.Error())
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func mergeImagesHandler(c *gin.Context) {
	if err := mergeImages(c); err != nil {
		log.Error("fail to merge images: ", err)
		writeErrorResp(c, err)
		return
	}
	c.String(http.StatusOK, "done")
}
