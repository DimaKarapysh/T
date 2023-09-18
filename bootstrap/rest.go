package bootstrap

import (
	"T/delivery"
	"T/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func InitRest() error {
	if os.Getenv("APP_DEBUG") != "false" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		return err
	}
	var n int
	log.Println("count procedure")
	_, err = fmt.Scan(&n)
	s := service.NewService(3)

	delivery.NewHandler(s).Router(r.Group("/task"))
	err = r.Run(":" + os.Getenv("REST_PORT"))
	if err != nil {
		return err
	}
	return nil
}
