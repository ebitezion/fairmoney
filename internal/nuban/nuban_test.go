package nuban_test

import (
	"strconv"
	"testing"

	"github.com/ebitezion/backend-framework/internal/nuban"
)

func TestGenerateNUBAN(t *testing.T) {
	generator := nuban.NewNUBANGenerator()

	// Generate NUBAN
	nubanCode := generator.GenerateNUBAN()

	// Check length
	if len(nubanCode) != 10 {
		t.Errorf("NUBAN length is not 10 characters: %s", nubanCode)
	}

	// Check if all characters are digits
	for _, char := range nubanCode {
		if !isDigit(string(char)) {
			t.Errorf("NUBAN contains non-digit character: %s", nubanCode)
		}
	}

	// Check if bank code, branch code, and account number are separated correctly
	bankCode := nubanCode[:3]
	branchCode := nubanCode[3:6]
	accountNumber := nubanCode[6:]

	if !isValidLength(bankCode, 3) || !isValidLength(branchCode, 3) || !isValidLength(accountNumber, 4) {
		t.Errorf("NUBAN structure is incorrect: %s", nubanCode)
	}
}

func isDigit(char string) bool {
	_, err := strconv.Atoi(char)
	return err == nil
}

func isValidLength(str string, length int) bool {
	return len(str) == length && isDigit(str)
}
