package handler

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	kind *Kind
}

func New() Handler {
	// It assumes kind only managed by `kind-manager`
	return Handler{
		kind: &Kind{Status: StatusNotExist},
	}
}

const (
	StatusNotExist = "Not exist"
	StatusCreating = "Creating"
	StatusDeleting = "Deleting"
	StatusRunning = "Running"
)

type Kind struct {
	Status string `json:"status"`
}

// status: not-exist, creating, deleting, running,
// /kind
// /kind/create
// /kind/delete
func (h Handler) KindGet(c *gin.Context) {
	c.JSON(http.StatusOK, h.kind)
}

func (h *Handler) KindCreatePut(c *gin.Context) {
	if h.kind.Status != StatusNotExist {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Couldn't create kind cluster since it is already exist",
			"status": h.kind.Status,
		})
		return
	}
	configPath := c.PostForm("config_path")
	log.Println("config path is: " + configPath)

	cmd := exec.Command("/usr/local/bin/kind",
		"create",
				"cluster",
				"--config",
				configPath)

	go func() {
		h.kind.Status = StatusCreating
		err := cmd.Start()
		if err != nil {
			log.Println("Failed to start to creating kind cluster: " + err.Error())
			return
		}
		err = cmd.Wait()
		if err != nil {
			log.Println("error occurred during creating kind cluster: " + err.Error())
			h.kind.Status = StatusNotExist

		} else {
			// TODO: Print elapsed time to create cluster
			log.Println("kind cluster created!!")
			h.kind.Status = StatusRunning
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully triggered to creating kind cluster",
		"status": h.kind.Status,
	})
}

func (h *Handler) KindDestroyDelete(c *gin.Context) {
	if h.kind.Status != StatusRunning {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Couldn't delete kind cluster since it does not exist",
			"status": h.kind.Status,
		})
		return
	}

	cmd := exec.Command("/usr/local/bin/kind",
		"delete",
		"cluster",
		)

	go func() {
		h.kind.Status = StatusDeleting
		err := cmd.Start()
		if err != nil {
			log.Println("Failed to start to deleting kind cluster: " + err.Error())
			return
		}
		err = cmd.Wait()
		if err != nil {
			log.Println("error occurred during deleting kind cluster: " + err.Error())
			// FIXME: It doesn't seems to be appropriate status
			h.kind.Status = StatusRunning

		} else {
			// TODO: Print elapsed time to delete cluster
			log.Println("kind cluster deleted!!")
			h.kind.Status = StatusNotExist
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully triggered to deleting kind cluster",
		"status": h.kind.Status,
	})
}
