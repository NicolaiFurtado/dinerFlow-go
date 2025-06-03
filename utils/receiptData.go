package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type ReceiptItem struct {
	Number   int
	Cod      int
	Name     string
	Quantity int
	Price    float64
}

type ReceiptData struct {
	ClientName    string
	CreatedAt     time.Time
	ClosedAt      time.Time
	TableName     string
	UserFirstName string
	UserLastName  string
	Items         []ReceiptItem
	TotalValue    float64
}

// GenerateReceipt writes a formatted receipt to a file
func GenerateReceipt(filename string, data ReceiptData) error {
	var sb strings.Builder

	// Header
	sb.WriteString("Nota Fiscal\n")
	sb.WriteString("DinerFlow\n\n")

	sb.WriteString(fmt.Sprintf("Cliente: %s\n", data.ClientName))
	sb.WriteString(fmt.Sprintf("Entrada: %s\n", data.CreatedAt.Format("02/01/2006 15:04:05")))
	sb.WriteString(fmt.Sprintf("Fechamento: %s\n", data.ClosedAt.Format("02/01/2006 15:04:05")))

	elapsed := data.ClosedAt.Sub(data.CreatedAt).Round(time.Second)
	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60
	sb.WriteString(fmt.Sprintf("Tempo Total: %02d:%02d:%02d\n\n", hours, minutes, seconds))

	sb.WriteString(fmt.Sprintf("Mesa: %s\n", data.TableName))
	sb.WriteString(fmt.Sprintf("Atendente: %s %s\n", data.UserFirstName, data.UserLastName))

	sb.WriteString("\n----------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("%-2s | %-5s | %-20s | %3s | %6s\n", "#", "Cod", "Item", "Qtd", "Valor"))
	sb.WriteString("\n----------------------------------------------------\n")
	for _, item := range data.Items {
		sb.WriteString(fmt.Sprintf("%-2d | %-5d | %-20s | %3d | %6.2f\n", item.Number, item.Cod, item.Name, item.Quantity, item.Price))
	}
	sb.WriteString("----------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("%-41s %6.2f\n", "Valor a Pagar:", data.TotalValue))
	sb.WriteString("\nCPF na Nota:____________________________________\n")

	// Write to file
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
