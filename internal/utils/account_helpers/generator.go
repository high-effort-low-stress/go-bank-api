// Package accounthelpers provides helper functions for account-related operations, such as generating check digits.
package accounthelpers

import (
	"fmt"
	"strconv"
)

func GenerateDigit(numero string) (int, error) {
	soma := 0
	peso := 2

	// Itera sobre a string da direita para a esquerda
	for i := len(numero) - 1; i >= 0; i-- {
		digito, err := strconv.Atoi(string(numero[i]))
		if err != nil {
			return 0, fmt.Errorf("provided number contains non-numeric characters")
		}

		soma += digito * peso

		// O peso cicla de 2 a 9
		peso++
		if peso > 9 {
			peso = 2
		}
	}

	resto := soma % 11
	dv := 11 - resto

	// Aplica a regra especial
	if dv == 10 || dv == 11 || dv == 0 { // Algumas implementações consideram o resto 0 (dv=11) como 0 também
		return 0, nil
	}

	return dv, nil
}
