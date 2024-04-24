package nuban

import (
	"crypto/rand"
	"math/big"
)

type NUBANGenerator struct {
}

func NewNUBANGenerator() *NUBANGenerator {
	return &NUBANGenerator{}
}

func (n *NUBANGenerator) generateRandomNumber(length int) string {
	var result string

	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		result += num.String()
	}

	return result
}

func (n *NUBANGenerator) GenerateNUBAN() string {
	bankCode := n.generateRandomNumber(3)
	branchCode := n.generateRandomNumber(3)
	accountNumber := n.generateRandomNumber(4)

	return bankCode + branchCode + accountNumber
}

//Test
// func main() {
// 	nubanGenerator := NewNUBANGenerator()
// 	nuban := nubanGenerator.GenerateNUBAN()
// 	fmt.Println("Generated NUBAN:", nuban)
// }
