package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"kvStorage/internal/usecase"
	"kvStorage/types"
	"net/http"
)

type Server struct {
	uc  *usecase.UseCase
	log *logrus.Logger
}

func New(uc *usecase.UseCase, log *logrus.Logger) *Server {
	return &Server{
		uc:  uc,
		log: log,
	}
}

func (s *Server) GetValue(c *gin.Context) {
	key := c.Param("key")
	val, err := s.uc.GetValue(key)
	if err != nil {
		if errors.Is(err, types.ErrKeyNotFound) {
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Key %s not found", key)})
			return
		}

		s.log.WithError(err).Errorln("Error getting value")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": val})
	return
}

func (s *Server) CreateValue(c *gin.Context) {
	var pair types.KeyValue
	err := c.Bind(&pair)
	if err != nil {
		s.log.Error("Invalid pair of key value ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(pair)
	if err != nil {
		s.log.Error("Invalid key or value ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = s.uc.PutValue(pair)
	if err != nil {
		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	return
}

func (s *Server) EditeValue(c *gin.Context) {
	var pair types.KeyValue
	err := c.Bind(&pair)
	if err != nil {
		s.log.Error("Invalid pair of key value ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(pair)
	if err != nil {
		s.log.Error("Invalid key or value ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = s.uc.SetValue(pair)
	if err != nil {
		if errors.Is(err, types.ErrKeyNotFound) {
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Key %s not found", pair.Key)})
			return
		}

		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	return
}

func (s *Server) DeleteValue(c *gin.Context) {
	key := c.Param("key")

	err := s.uc.DeleteValue(key)
	if err != nil {
		if errors.Is(err, types.ErrKeyNotFound) {
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Key %s not found", key)})
			return
		}

		s.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	return
}
