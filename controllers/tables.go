package controllers

import (
	"database/sql"
	"dinerFlow/config"
	"dinerFlow/models"
	"dinerFlow/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

// GetTables godoc
// @Summary List tables
// @Description List diner tables, filtered optionally by status, seats or table name
// @Tags tables
// @Accept json
// @Produce json
// @Param table_name query string false "Partial match of table name"
// @Param seats query int false "Exact number of seats"
// @Param status query string false "Table status (available, occupied, deleted)"
// @Success 200 {array} models.Tables
// @Failure 500 {object} models.ErrorResponse
// @Router /tables [get]
func GetTables(c *gin.Context) {
	tableName := c.Query("table_name")
	seatsStr := c.Query("seats")
	status := c.Query("status")

	query := "SELECT id, table_name, seats, status, created_by, updated_info FROM diner_tables WHERE 1=1"
	args := []interface{}{}

	// If status is not explicitly "deleted", default is to exclude deleted
	if status == "" {
		query += " AND status != ?"
		args = append(args, "deleted")
	} else {
		query += " AND status = ?"
		args = append(args, status)
	}

	if tableName != "" {
		query += " AND table_name LIKE ?"
		args = append(args, "%"+tableName+"%")
	}

	if seatsStr != "" {
		seats, err := strconv.Atoi(seatsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'seats' inválido"})
			return
		}
		query += " AND seats = ?"
		args = append(args, seats)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar as mesas"})
		return
	}
	defer rows.Close()

	var dinerTables []models.Tables
	for rows.Next() {
		var t models.Tables
		if err := rows.Scan(&t.ID, &t.TableName, &t.Seats, &t.Status, &t.CreatedBy, &t.UpdatedInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar as mesas"})
			return
		}
		dinerTables = append(dinerTables, t)
	}

	c.JSON(http.StatusOK, dinerTables)
}

// CreateTable godoc
// @Summary Create new table
// @Description Create a new diner table
// @Tags tables
// @Accept json
// @Produce json
// @Param table body models.Tables true "New table data"
// @Success 201 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tables [post]
func CreateTable(c *gin.Context) {
	var t models.Tables

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	_, err := config.DB.Exec("INSERT INTO diner_tables (table_name, seats, created_by) VALUES (?, ?, ?)", t.TableName, t.Seats, userId)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"error": "Já existe uma mesa com esse nome"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar a mesa", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Mesa criada com sucesso"})
}

// EditTable godoc
// @Summary Edit existing table
// @Description Update diner table info
// @Tags tables
// @Accept json
// @Produce json
// @Param table body models.Tables true "Updated table data"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tables [put]
func EditTable(c *gin.Context) {
	var t models.Tables

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var old models.Tables
	var existingAuditJSON sql.NullString

	err := config.DB.QueryRow(
		"SELECT id, table_name, seats, updated_info FROM diner_tables WHERE id = ?",
		t.ID,
	).Scan(&old.ID, &old.TableName, &old.Seats, &existingAuditJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar dados antigos", "details": err.Error()})
		return
	}

	// Gera JSON atualizado da auditoria
	finalAuditJSON, err := utils.AppendAuditLog(existingAuditJSON, userId, old, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar JSON de auditoria", "details": err.Error()})
		return
	}

	_, err = config.DB.Exec(
		"UPDATE diner_tables SET table_name = ?, seats = ?, updated_info = ? WHERE id = ?",
		t.TableName, t.Seats, finalAuditJSON, t.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao editar a mesa", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesa editada com sucesso"})
}

// DeleteTable godoc
// @Summary Soft delete table
// @Description Set table status to 'deleted'
// @Tags tables
// @Accept json
// @Produce json
// @Param table body object true "Table ID to delete" schema:{ "id": 1 }
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tables [delete]
func DeleteTable(c *gin.Context) {
	type Request struct {
		ID int `json:"id"`
	}

	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var old models.Tables
	var existingAuditJSON sql.NullString
	var status string

	err := config.DB.QueryRow(
		"SELECT id, table_name, seats, status, updated_info FROM diner_tables WHERE id = ?",
		req.ID,
	).Scan(&old.ID, &old.TableName, &old.Seats, &status, &existingAuditJSON)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar dados da mesa", "details": err.Error()})
		return
	}

	if status == "deleted" {
		c.JSON(http.StatusConflict, gin.H{"error": "Mesa já está deletada"})
		return
	}

	// Gera JSON atualizado de auditoria
	finalAuditJSON, err := utils.AppendAuditLog(existingAuditJSON, userId, old, map[string]interface{}{
		"status": "deleted",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar JSON de auditoria", "details": err.Error()})
		return
	}

	_, err = config.DB.Exec("UPDATE diner_tables SET status = 'deleted', updated_info = ? WHERE id = ?", finalAuditJSON, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar a mesa", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesa deletada com sucesso"})
}
