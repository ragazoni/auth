package service

import (
	"fmt"
	"strings"
)

// func IsAdmin(email string) (bool, error) {
// 	user, err := db.FindUserByEmail("users", email)
// 	if err != nil {
// 		return false, err
// 	}
// 	if user == nil {
// 		return false, err
// 	}

// 	isAdm := (user.Type == "admin")

// 	return isAdm, err
// }

// Função para verificar se um usuário é administrador baseado no domínio do e-mail
func IsAdmin(email string) (bool, error) {
	if email == "" {
		return false, fmt.Errorf("email cannot be empty")
	}
	if strings.HasSuffix(email, "@br.experian.com") {
		return true, nil
	}
	return false, nil
}

// func isValidExperianEmail(email string) bool {
// 	// Expressão regular para validar se o email termina com @br.experian.com
// 	re := regexp.MustCompile(`@br.experian.com$`)
// 	return re.MatchString(email)
// }
