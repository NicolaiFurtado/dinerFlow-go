package controllers

import (
	"database/sql"
	"dinerFlow/config"
	"dinerFlow/models"
	"dinerFlow/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// GetItems godoc
// @Summary List Items
// @Description List diner items, filtered optionally by name and category.
// @Tags Diner Items
// @Accept json
// @Produce json
// @Param name query string false "Partial match of item name"
// @Param category query string false "Item category"
// @Success 200 {array} models.Items
// @Failure 500 {object} models.ErrorResponse
// @Router /items [get]
func GetItems(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")

	query := "SELECT id, cod, name, description, price, category, created_by, updated_info FROM diner_items WHERE 1=1"
	args := []interface{}{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		query += " AND category LIKE ?"
		args = append(args, "%"+category+"%")
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar os itens"})
		return
	}
	defer rows.Close()

	var dinerItems []models.Items
	for rows.Next() {
		var t models.Items
		var updatedInfo sql.NullString
		if err := rows.Scan(&t.ID, &t.Cod, &t.Name, &t.Description, &t.Price, &t.Category, &t.CreatedBy, &updatedInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar os itens", "details": err.Error()})
			return
		}
		t.UpdatedInfo = updatedInfo.String
		dinerItems = append(dinerItems, t)
	}

	c.JSON(http.StatusOK, dinerItems)
}

// CreateItem godoc
// @Summary Create Items
// @Description Create diner items.
// @Tags Diner Items
// @Accept json
// @Produce json
// @Param table body models.Items true "New table data"
// @Success 201 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /items [post]
func CreateItem(c *gin.Context) {
	var item models.Items
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var generatedCod string
	for {
		codInt := 10000 + rand.Intn(90000) // Random number between 10000 and 99999
		generatedCod = fmt.Sprintf("%05d", codInt)

		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM diner_items WHERE cod = ?)", generatedCod).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar código existente"})
			return
		}

		if !exists {
			break
		}
	}

	query := "INSERT INTO diner_items (cod, name, description, price, category, created_by) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := config.DB.Exec(query, generatedCod, item.Name, item.Description, item.Price, item.Category, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar o item"})
		return
	}

	id, _ := result.LastInsertId()
	item.ID = int(id)

	c.JSON(http.StatusCreated, item)
}

// EditItem godoc
// @Summary Edit existing diner items
// @Description Update diner items
// @Tags Diner Items
// @Accept json
// @Produce json
// @Param table body models.Items true "Updated table data"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /items [put]
func EditItem(c *gin.Context) {
	var i models.Items

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var old models.Items
	var existingAuditJSON sql.NullString

	err := config.DB.QueryRow(
		"SELECT id, cod, name, description, price, category, updated_info FROM diner_items WHERE id = ?",
		i.ID,
	).Scan(&old.ID, &old.Cod, &old.Name, &old.Description, &old.Price, &old.Category, &existingAuditJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar dados antigos", "details": err.Error()})
		return
	}

	// Gera JSON atualizado da auditoria
	finalAuditJSON, err := utils.AppendAuditLog(existingAuditJSON, userId, old, i)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar JSON de auditoria", "details": err.Error()})
		return
	}

	result, err := config.DB.Exec(
		"UPDATE diner_items SET name = ?, description = ?, price = ?, category = ?, updated_info = ? WHERE id = ? AND cod = ?",
		i.Name, i.Description, i.Price, i.Category, finalAuditJSON, i.ID, i.Cod,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao editar o item", "details": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado ou código incorreto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item editada com sucesso"})
}

// DeleteItem godoc
// @Summary Delete an existing diner item
// @Description Delete a diner item by ID and Cod
// @Tags Diner Items
// @Accept json
// @Produce json
// @Param item body models.Items true "Item to delete"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /items [delete]
func DeleteItem(c *gin.Context) {
	var i models.Items

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	var old models.Items
	err := config.DB.QueryRow(
		"SELECT id, cod FROM diner_items WHERE id = ? AND cod = ?",
		i.ID, i.Cod,
	).Scan(&old.ID, &old.Cod)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar dados do item", "details": err.Error()})
		return
	}

	result, err := config.DB.Exec("DELETE FROM diner_items WHERE id = ? AND cod = ?", i.ID, i.Cod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar o item", "details": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado ou código incorreto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deletada com sucesso"})
}
