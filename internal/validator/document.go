package validator

import (
	"regexp"
	"strconv"
)

type Document struct {
	CPFRegexp         *regexp.Regexp
	NumbersOnlyRegexp *regexp.Regexp
	CPFFormatRegexp   *regexp.Regexp
	CPFFormatLayout   string
}

func newDocument() *Document {
	return &Document{
		CPFRegexp:         regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`),
		NumbersOnlyRegexp: regexp.MustCompile(`[^a-zA-Z0-9]+`),
		CPFFormatRegexp:   regexp.MustCompile(`(\d{3})(\d{3})(\d{3})(\d{2})$`),
		CPFFormatLayout:   "$1.$2.$3-$4",
	}
}

// IsCPF verifica se a string dada é um documento CPF válido.
func (c *Document) IsCPF(doc string) bool {
	const (
		size = 9
		pos  = 10
	)

	return c.isCPF(doc, c.CPFRegexp, size, pos)
}

// isCPF gera os dígitos para um dado CPF e compara com os dígitos originais.
func (c *Document) isCPF(doc string, pattern *regexp.Regexp, size int, position int) bool {
	if !pattern.MatchString(doc) {
		return false
	}

	doc = c.CleanNonDigits(doc)

	// Invalida documentos com todos os dígitos iguais.
	if allEq(doc) {
		return false
	}

	d := doc[:size]
	digit := calculateDigit(d, position)

	d = d + digit
	digit = calculateDigit(d, position+1)

	return doc == d+digit
}

// CleanNonDigits remove todos os caracteres que não são dígitos.
func (c *Document) CleanNonDigits(document string) (retVal string) {
	return c.NumbersOnlyRegexp.ReplaceAllString(document, "")
}

// FormatCPF formata uma string como CPF.
func (c *Document) FormatCPF(document string) (retVal string) {
	return c.CPFFormatRegexp.ReplaceAllString(c.CleanNonDigits(document), c.CPFFormatLayout)
}

// allEq verifica se todos os caracteres de uma string são iguais.
func allEq(doc string) bool {
	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}
	return true
}

// calculateDigit calcula o próximo dígito para o documento dado.
func calculateDigit(doc string, position int) string {
	var sum int
	for _, r := range doc {
		sum += toInt(r) * position
		position--
		if position < 2 {
			position = 9
		}
	}

	sum %= 11
	if sum < 2 {
		return "0"
	}
	return strconv.Itoa(11 - sum)
}

// toInt converte um rune para int.
func toInt(r rune) int {
	return int(r - '0')
}

// func main() {
// 	// Exemplo de uso
// 	doc := NewDocument()
// 	cpf := "123.456.789-09"
// 	valid := doc.IsCPF(cpf)
// 	println("CPF válido:", valid)
// }
