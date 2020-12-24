package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gurrpi/kind-manager/server/handler"
)

const Port string = ":15050"

type Server struct {
	instance *gin.Engine
}

func New() Server {
	r := gin.Default()

	h := handler.New()
	r.GET("/kind", h.KindGet)
	r.PUT("/kind", h.KindCreatePut)
	r.DELETE("/kind", h.KindDestroyDelete)
	return Server{
		instance: r,
	}
}

func (s Server) Run() error {
	err := s.instance.Run(Port)
	if err != nil {
		return err
	}
	return nil
}
