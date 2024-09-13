package handler

import (
	"net/http"

	"github.com/Mamvriyskiy/database_course/main/logger"
	"github.com/Mamvriyskiy/database_course/main/pkg"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createHome(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		logger.Log("Warning", "Get", "Error get userID from context", nil, "userId")
		return
	}

	var input pkg.Home
	if err := c.BindJSON(&input); err != nil {
		logger.Log("Error", "c.BindJSON()", "Error bind json:", err, "")
		return
	}

	var intVal float64
	if val, ok := id.(float64); ok {
		intVal = val
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"errors": "Ошибка создания дома",
		})
		logger.Log("Error", "userID.(float64)", "Error:", ErrNoFloat64Interface, "")
		return
	}

	homeID, err := h.services.IHome.CreateHome(input)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"errors": "Ошибка создания дома",
		})
		logger.Log("Error", "CreateHome", "Error create home:", err, intVal, input)
		return
	}

	_, err = h.services.IAccessHome.AddOwner(int(intVal), homeID)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"errors": "Ошибка добавления хозяина дома",
		})
		logger.Log("Error", "AddOwner", "Error add owner:", err, intVal, homeID)
		return
	}

	c.Set("homeID", homeID)
	c.Next()

	c.JSON(http.StatusOK, map[string]interface{}{
		"homeId": homeID,
	})

	logger.Log("Info", "", "A home has been created", nil)
}

func (h *Handler) deleteHome(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		logger.Log("Warning", "Get", "Error get userID from context", nil, "userID")
		return
	}

	var input pkg.Home
	if err := c.BindJSON(&input); err != nil {
		logger.Log("Error", "c.BindJSON()", "Error bind json:", err, "")
		return
	}

	var intVal float64
	if val, ok := id.(float64); ok {
		intVal = val
	} else {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errors": "Ошибка обновления",
		})
		logger.Log("Error", "userID.(float64)", "Error:", ErrNoFloat64Interface, "")
		return
	}

	err := h.services.IHome.DeleteHome(int(intVal), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errors": "дом не найден",
		})
		logger.Log("Error", "DeleteHome", "Error delete home:", err, id.(int))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{})

	logger.Log("Info", "", "A home has been deleted", nil)
}

type getAllListHomeResponse struct {
	Data []pkg.Home `json:"data"`
}

func (h *Handler) listHome(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		logger.Log("Warning", "Get", "Error get userID from context", nil, "userID")
		return
	}

	var intVal float64
	if val, ok := id.(float64); ok {
		intVal = val
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"errors": "Ошибка получения списка домов",
		})
		logger.Log("Error", "userID.(float64)", "Error:", ErrNoFloat64Interface, "")
		return 
	}

	homeListUser, err := h.services.IHome.ListUserHome(int(intVal))
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"errors": "Ошибка получения списка домов",
		})
		logger.Log("Error", "ListUserHome", "Error get user:", err, id.(int))
		return
	}

	c.JSON(http.StatusOK, getAllListHomeResponse{
		Data: homeListUser,
	})

	logger.Log("Info", "", "The list of users has been received", nil)
}

func (h *Handler) updateHome(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		logger.Log("Warning", "Get", "Error get userID from context", nil, "userID")
		return
	}

	var input pkg.UpdateNameHome
	err := c.BindJSON(&input)
	if err != nil {
		logger.Log("Error", "c.BindJSON()", "Error bind json:", err, "")
		return
	}

	var intVal float64
	if val, ok := id.(float64); ok {
		intVal = val
	} else {
		logger.Log("Error", "userID.(float64)", "Error:", ErrNoFloat64Interface, "")
		c.JSON(http.StatusInternalServerError, map[string]string{
			"errors": "Ошибка обновления",
		})
		return
	}
	input.UserID = int(intVal)

	err = h.services.IHome.UpdateHome(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": "дом не найден",
		})
		logger.Log("Error", "UpdateHome", "Error update home:", err, "")
		return
	}

	c.JSON(http.StatusOK, getAllListHomeResponse{})

	logger.Log("Info", "", "A home has been update", nil)
}
