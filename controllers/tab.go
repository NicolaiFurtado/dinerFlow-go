package controllers

import (
	"database/sql"
	"dinerFlow/config"
	"dinerFlow/models"
	"dinerFlow/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenTab godoc
// @Summary Open a new tab for a table
// @Description Opens a new diner tab if the table exists and is not already associated with an open tab
// @Tags Diner Tabs
// @Accept json
// @Produce json
// @Param tab body models.Tab true "Tab data (table ID, client name, order)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tab [post]
func OpenTab(c *gin.Context) {
	var i models.Tab

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var tableID int
	err := config.DB.QueryRow(
		"SELECT id FROM diner_tables WHERE id = ?",
		i.TableId,
	).Scan(&tableID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mesa não encontrada"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar mesa", "details": err.Error()})
		return
	}

	var existingTableID int
	err = config.DB.QueryRow(
		"SELECT table_id FROM diner_tab WHERE table_id = ? AND status = 'open'",
		i.TableId,
	).Scan(&existingTableID)
	if err != sql.ErrNoRows {
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Já existe uma comanda aberta para esta mesa"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar comandas abertas", "details": err.Error()})
		}
		return
	}

	// Decode i.Order JSON and validate each item_cod in its items slice
	var order struct {
		Items []struct {
			ItemCod int `json:"item_cod"`
			// Add other fields here if needed
		} `json:"items"`
	}

	for _, item := range order.Items {
		var exists bool
		err := config.DB.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM diner_items WHERE cod = ?)",
			item.ItemCod,
		).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar item no banco", "details": err.Error()})
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item não encontrado no cardápio", "item_cod": item.ItemCod})
			return
		}
	}
	_, err = config.DB.Exec(
		"INSERT INTO diner_tab (table_id, client_name, `order`, status, created_at, created_by) VALUES (?, ?, ?, 'open', ?, ?)",
		i.TableId, i.ClientName, i.Order, time.Now(), userId,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"error": "Já existe uma comanda aberta para esta mesa"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar a comanda", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comanda aberta com sucesso para o cliente: " + i.ClientName, "table_id": i.TableId})
}

// UpdateOrder godoc
// @Summary Add new items to an open diner tab
// @Description Merges new items into the current order of an open tab. Preserves existing items.
// @Tags Diner Tabs
// @Accept json
// @Produce json
// @Param order body models.Tab true "Tab ID and new order items to be merged"
// @Success 200 {object} map[string]string "Pedidos atualizados com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos ou erro de parsing"
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Failure 409 {object} map[string]string "Comanda não encontrada ou já fechada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /tab [put]
func UpdateOrder(c *gin.Context) {
	var i models.Tab

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var existingTableID int
	err := config.DB.QueryRow(
		"SELECT id FROM diner_tab WHERE id = ? AND status = 'open'",
		i.ID,
	).Scan(&existingTableID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Não existe uma comanda aberta com este ID"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar comandas abertas", "details": err.Error()})
		return
	}

	// Step 1: Get current orders from the database
	var currentOrderData []byte
	err = config.DB.QueryRow(
		"SELECT `order` FROM diner_tab WHERE id = ?",
		i.ID,
	).Scan(&currentOrderData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter pedidos atuais", "details": err.Error()})
		return
	}

	// Step 2: Parse current and new orders
	var currentOrder, newOrder struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(currentOrderData, &currentOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao interpretar pedidos atuais", "details": err.Error()})
		return
	}

	// Step 3: Merge orders
	mergedItems := append(currentOrder.Items, newOrder.Items...)
	mergedOrder := map[string]interface{}{
		"items": mergedItems,
	}
	mergedOrderJSON, err := json.Marshal(mergedOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao preparar pedidos para atualização", "details": err.Error()})
		return
	}

	var old models.Tab
	var existingAuditJSON sql.NullString

	err = config.DB.QueryRow(
		"SELECT id, `order`, status, updated_info FROM diner_tab WHERE id = ?",
		i.ID,
	).Scan(&old.ID, &old.Order, &old.Status, &existingAuditJSON)
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

	// Step 4: Update database
	_, err = config.DB.Exec(
		"UPDATE diner_tab SET `order` = ?, updated_info = ? WHERE id = ?",
		mergedOrderJSON, finalAuditJSON, i.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar pedidos", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pedidos atualizados com sucesso"})
}

// RemoveOrder godoc
// @Summary Remove items from an open diner tab
// @Description Removes specific item_cod entries from the current order of an open tab
// @Tags Diner Tabs
// @Accept json
// @Produce json
// @Param removeRequest body models.RemoveOrderRequest true "Order item_cod list to remove"
// @Success 200 {object} map[string]string "Itens removidos com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos ou erro de parsing"
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Failure 404 {object} map[string]string "Comanda não encontrada ou já fechada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /tab/remove [put]
func RemoveOrder(c *gin.Context) {
	var removeRequest models.RemoveOrderRequest
	if err := c.ShouldBindJSON(&removeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var existingTabID int
	err := config.DB.QueryRow(
		"SELECT id FROM diner_tab WHERE id = ? AND status = 'open'",
		removeRequest.ID,
	).Scan(&existingTabID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comanda não encontrada ou já fechada"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar comanda", "details": err.Error()})
		return
	}

	var currentOrderData []byte
	err = config.DB.QueryRow(
		"SELECT `order` FROM diner_tab WHERE id = ?",
		removeRequest.ID,
	).Scan(&currentOrderData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter pedidos atuais", "details": err.Error()})
		return
	}

	var currentOrder struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(currentOrderData, &currentOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao interpretar pedidos atuais", "details": err.Error()})
		return
	}

	removeSet := make(map[int]bool)
	for _, id := range removeRequest.Order.Remove {
		removeSet[id] = true
	}

	var filteredItems []map[string]interface{}
	for _, item := range currentOrder.Items {
		if itemCod, ok := item["item_cod"].(float64); ok {
			if !removeSet[int(itemCod)] {
				filteredItems = append(filteredItems, item)
			}
		}
	}

	newOrder := map[string]interface{}{
		"items": filteredItems,
	}
	newOrderJSON, err := json.Marshal(newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao preparar nova ordem", "details": err.Error()})
		return
	}

	var old models.Tab
	var existingAuditJSON sql.NullString

	err = config.DB.QueryRow(
		"SELECT id, `order`, status, updated_info FROM diner_tab WHERE id = ?",
		removeRequest.ID,
	).Scan(&old.ID, &old.Order, &old.Status, &existingAuditJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar dados antigos", "details": err.Error()})
		return
	}

	finalAuditJSON, err := utils.AppendAuditLog(existingAuditJSON, userId, old, removeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar JSON de auditoria", "details": err.Error()})
		return
	}

	_, err = config.DB.Exec(
		"UPDATE diner_tab SET `order` = ?, updated_info = ? WHERE id = ?",
		newOrderJSON, finalAuditJSON, removeRequest.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar pedidos", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itens removidos com sucesso"})
}

// CloseTab godoc
// @Summary Fecha a comanda e registra o pagamento
// @Description Fecha uma comanda aberta, gera o recibo em arquivo e insere os dados na tabela de pagamentos
// @Tags Diner Tabs
// @Accept json
// @Produce json
// @Param tab body models.Tab true "Dados da comanda a ser fechada (ID)"
// @Success 200 {object} map[string]interface{} "Comanda fechada com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Failure 409 {object} map[string]string "Comanda não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /tab/close [put]
func CloseTab(c *gin.Context) {
	var i models.Tab

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var (
		existingTableID   int
		tableId           int
		clientName        string
		createdAt         time.Time
		createdAtRaw      []byte
		createdBy         int
		existingOrderData []byte
		tableName         string
		userFirstName     string
		userLastName      string
		itemsWithPrices   []utils.ReceiptItem
	)

	err := config.DB.QueryRow(`
		SELECT t.id, t.table_id, t.client_name, t.order, t.created_at, t.created_by, dt.table_name, u.firstname, u.lastname
		FROM diner_tab t
		JOIN diner_tables dt ON dt.id = t.table_id
		JOIN users u ON u.id = t.created_by
		WHERE t.id = ? AND t.status = 'open'
	`, i.ID).Scan(&existingTableID, &tableId, &clientName, &existingOrderData, &createdAtRaw, &createdBy, &tableName, &userFirstName, &userLastName)

	// Parse createdAtRaw into createdAt
	if err == nil {
		createdAt, err = time.Parse("2006-01-02 15:04:05", string(createdAtRaw))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao interpretar created_at", "details": err.Error()})
			return
		}
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Não existe uma comanda aberta com este ID"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados da comanda", "details": err.Error()})
		return
	}

	var parsedOrder struct {
		Items []struct {
			ItemCod int    `json:"item_cod"`
			Qtd     int    `json:"qtd"`
			Notes   string `json:"notes"`
		} `json:"items"`
	}

	if err = json.Unmarshal(existingOrderData, &parsedOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao interpretar pedidos", "details": err.Error()})
		return
	}

	var totalPrice float64

	for i, item := range parsedOrder.Items {
		var (
			name  string
			price float64
		)

		err = config.DB.QueryRow("SELECT name, price FROM diner_items WHERE cod = ?", item.ItemCod).Scan(&name, &price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar detalhes do item", "item_cod": item.ItemCod, "details": err.Error()})
			return
		}

		itemTotal := price * float64(item.Qtd)
		totalPrice += itemTotal

		itemsWithPrices = append(itemsWithPrices, utils.ReceiptItem{
			Number:   i + 1,
			Cod:      item.ItemCod,
			Name:     name,
			Quantity: item.Qtd,
			Price:    price,
		})
	}

	receiptData := utils.ReceiptData{
		ClientName:    clientName,
		CreatedAt:     createdAt,
		ClosedAt:      time.Now(),
		TableName:     tableName,
		UserFirstName: userFirstName,
		UserLastName:  userLastName,
		Items:         itemsWithPrices,
		TotalValue:    totalPrice,
	}

	now := time.Now()
	sanitizedName := strings.ReplaceAll(strings.ToLower(clientName), " ", "-")
	dir := fmt.Sprintf("public/%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	filename := fmt.Sprintf("%s/%s_%s.txt", dir, now.Format("2006-01-02"), sanitizedName)

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar diretório do recibo", "details": err.Error()})
		return
	}
	if err = utils.GenerateReceipt(filename, receiptData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar recibo", "details": err.Error()})
		return
	}

	receiptJSON, err := json.Marshal(receiptData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar JSON do recibo", "details": err.Error()})
		return
	}

	// Begin a transaction to ensure atomicity
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação", "details": err.Error()})
		return
	}
	defer tx.Rollback()

	// 1. Update diner_tab
	_, err = tx.Exec(`
	UPDATE diner_tab 
	SET status = 'closed'
	WHERE id = ?`, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fechar comanda", "details": err.Error()})
		return
	}

	// 2. Insert into diner_payments
	_, err = tx.Exec(`
INSERT INTO diner_payments 
(user_id, client_name, total_price, receipt_data, created_at, created_by)
VALUES (?, ?, ?, ?, NOW(), ?)`,
		userId, clientName, totalPrice, receiptJSON, userId)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao confirmar transação", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Comanda fechada e pagamento registrado com sucesso",
		"total_price":  totalPrice,
		"receipt_file": filename,
	})
}

// FinishPayment godoc
// @Summary Finaliza um pagamento
// @Description Atualiza o status de um pagamento para "Paid", define o tipo de pagamento e registra o fechamento
// @Tags Diner Payments
// @Accept json
// @Produce json
// @Param payment body models.Payment true "Dados do pagamento (ID e tipo de pagamento)"
// @Success 200 {object} map[string]string "Pagamento finalizado com sucesso"
// @Failure 400 {object} map[string]interface{} "Dados inválidos ou pagamento já finalizado"
// @Failure 401 {object} map[string]string "Usuário não autenticado"
// @Failure 404 {object} map[string]string "Pagamento não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /payment/finish [put]
func FinishPayment(c *gin.Context) {
	var i models.Payment

	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var existingPayment struct {
		Status string
	}
	err := config.DB.QueryRow(`
		SELECT status FROM diner_payments WHERE id = ?
	`, i.ID).Scan(&existingPayment.Status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pagamento não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pagamento", "details": err.Error()})
		return
	}

	if existingPayment.Status != "Not Paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pagamento já foi finalizado"})
		return
	}

	_, err = config.DB.Exec(`
		UPDATE diner_payments
		SET status = 'Paid', type_payment = ?, closed_at = NOW(), closed_by = ?
		WHERE id = ?
	`, i.TypePayment, userId, i.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar pagamento", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pagamento finalizado com sucesso"})
}
