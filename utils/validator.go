package utils

import (
	"regexp"
	"strconv"
)

// IsValidCPF verifica se uma string de CPF é válida de acordo com o algoritmo oficial.
// A função lida com CPFs formatados (ex: "123.456.789-00") ou apenas com números.
func IsValidCPF(cpf string) bool {
	// 1. Limpa o CPF, removendo todos os caracteres que não são dígitos.
	re := regexp.MustCompile("[^0-9]+")
	cpf = re.ReplaceAllString(cpf, "")

	// 2. Verifica se o CPF tem exatamente 11 dígitos.
	if len(cpf) != 11 {
		return false
	}

	// 3. Verifica se todos os dígitos são iguais (ex: "111.111.111-11").
	// Se forem, o CPF é inválido.
	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	// 4. Converte a string de dígitos em um slice de inteiros.
	digitos := make([]int, 11)
	for i, r := range cpf {
		digito, err := strconv.Atoi(string(r))
		if err != nil {
			return false // Deve ser impossível chegar aqui devido à limpeza com regex
		}
		digitos[i] = digito
	}

	// 5. Calcula o primeiro dígito verificador.
	var sum int
	for i := 0; i < 9; i++ {
		sum += digitos[i] * (10 - i)
	}
	resto := sum % 11
	digitoVerificador1 := 0
	if resto >= 2 {
		digitoVerificador1 = 11 - resto
	}

	// 6. Compara o primeiro dígito calculado com o dígito real do CPF.
	if digitos[9] != digitoVerificador1 {
		return false
	}

	// 7. Calcula o segundo dígito verificador.
	sum = 0
	for i := 0; i < 10; i++ {
		sum += digitos[i] * (11 - i)
	}
	resto = sum % 11
	digitoVerificador2 := 0
	if resto >= 2 {
		digitoVerificador2 = 11 - resto
	}

	// 8. Compara o segundo dígito calculado com o dígito real do CPF.
	if digitos[10] != digitoVerificador2 {
		return false
	}

	// Se todas as verificações passaram, o CPF é válido.
	return true
}
