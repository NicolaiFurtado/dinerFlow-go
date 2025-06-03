package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type PaymentSummary struct {
	User  string
	Value float64
}

var GenerateClosingReportMock func(filename string, summaries []PaymentSummary, closedBy string, now time.Time) error

func GenerateClosingReport(filename string, summaries []PaymentSummary, closedBy string, closedAt time.Time) error {
	if GenerateClosingReportMock != nil {
		return GenerateClosingReportMock(filename, summaries, closedBy, closedAt)
	}

	var sb strings.Builder

	sb.WriteString("DinerFlow\n\n")
	sb.WriteString(fmt.Sprintf("Data de Fechamento: %s\n", closedAt.Format("02/01/2006")))
	sb.WriteString(fmt.Sprintf("Horário de Fechamento: %s\n", closedAt.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("\nUsuário de Fechamento: %s\n\n", closedBy))

	sb.WriteString("----------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("%-4s| %-30s| %12s\n", "#", "Funcionário - Preço Unitário", "Preço Total"))
	sb.WriteString("----------------------------------------------------\n")

	userTotals := make(map[string]float64)
	line := 1
	var overallTotal float64

	for _, summary := range summaries {
		userTotals[summary.User] += summary.Value
		sb.WriteString(fmt.Sprintf("%-4d| %-30s| %12.2f\n", line, summary.User, summary.Value))
		line++
	}

	sb.WriteString("----------------------------------------------------\n")

	for user, total := range userTotals {
		sb.WriteString(fmt.Sprintf("Total %-30s| %12.2f\n", user, total))
		overallTotal += total
	}

	sb.WriteString("----------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Fechamento do Dia               | %12.2f\n\n", overallTotal))
	sb.WriteString("Assinatura do Fechamento: _________________________________\n")

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
