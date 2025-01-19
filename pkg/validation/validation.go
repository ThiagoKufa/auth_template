package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// PasswordPolicy define os requisitos para senhas
type PasswordPolicy struct {
	MinLength           int
	RequireUppercase    bool
	RequireLowercase    bool
	RequireNumbers      bool
	RequireSpecialChars bool
	DisallowedWords     []string
}

var defaultPolicy = PasswordPolicy{
	MinLength:           8,
	RequireUppercase:    true,
	RequireLowercase:    true,
	RequireNumbers:      true,
	RequireSpecialChars: true,
	DisallowedWords:     []string{"password", "123456", "qwerty"},
}

// ValidateEmail sanitiza e valida um endereço de email
func ValidateEmail(email string) (string, error) {
	// Sanitizar
	email = SanitizeString(email)

	// Verificar comprimento
	if len(email) < 3 || len(email) > 254 {
		return "", fmt.Errorf("email deve ter entre 3 e 254 caracteres")
	}

	// Validar formato
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("formato de email inválido")
	}

	// Validar domínio
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("email deve conter exatamente um @")
	}

	domain := parts[1]
	if len(domain) < 3 || !strings.Contains(domain, ".") {
		return "", fmt.Errorf("domínio de email inválido")
	}

	return email, nil
}

// ValidatePassword verifica se uma senha atende à política definida
func ValidatePassword(password string, policy PasswordPolicy) error {
	if policy.MinLength == 0 {
		policy = defaultPolicy
	}

	// Sanitizar
	password = SanitizeString(password)

	// Verificar comprimento
	if len(password) < policy.MinLength {
		return fmt.Errorf("senha deve ter pelo menos %d caracteres", policy.MinLength)
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return fmt.Errorf("senha deve conter pelo menos uma letra maiúscula")
	}
	if policy.RequireLowercase && !hasLower {
		return fmt.Errorf("senha deve conter pelo menos uma letra minúscula")
	}
	if policy.RequireNumbers && !hasNumber {
		return fmt.Errorf("senha deve conter pelo menos um número")
	}
	if policy.RequireSpecialChars && !hasSpecial {
		return fmt.Errorf("senha deve conter pelo menos um caractere especial")
	}

	// Verificar palavras não permitidas
	passwordLower := strings.ToLower(password)
	for _, word := range policy.DisallowedWords {
		if strings.Contains(passwordLower, strings.ToLower(word)) {
			return fmt.Errorf("senha contém palavra não permitida: %s", word)
		}
	}

	return nil
}

// SanitizeString limpa uma string de caracteres potencialmente perigosos
func SanitizeString(input string) string {
	// Remover espaços em branco extras
	input = strings.TrimSpace(input)

	// Remover caracteres de controle
	input = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(input, "")

	// Remover tags HTML
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")

	// Escapar caracteres especiais
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")

	return input
}
