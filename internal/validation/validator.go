package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	apperrors "server_kufatech/internal/errors"
)

type PasswordPolicy struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSpecial   bool
	DisallowedWords  []string
	MaxRepeatedChars int
	MinUniqueChars   int
}

var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:        8,
	RequireUppercase: true,
	RequireLowercase: true,
	RequireNumbers:   true,
	RequireSpecial:   true,
	DisallowedWords:  []string{"password", "123456", "qwerty"},
	MaxRepeatedChars: 3,
	MinUniqueChars:   5,
}

// ValidateEmail verifica se o email é válido e sanitiza
func ValidateEmail(email string) (string, error) {
	// Sanitizar
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	// Validar formato
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", apperrors.NewValidationError("email inválido")
	}

	// Extrair email limpo
	email = addr.Address

	// Validações adicionais
	if len(email) > 255 {
		return "", apperrors.NewValidationError("email muito longo")
	}

	// Verificar domínio
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", apperrors.NewValidationError("email inválido")
	}

	// Verificar caracteres especiais no domínio
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_.]+\.[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(parts[1]) {
		return "", apperrors.NewValidationError("domínio do email inválido")
	}

	return email, nil
}

// ValidatePassword verifica se a senha atende aos requisitos de segurança
func ValidatePassword(password string, policy PasswordPolicy) error {
	// Sanitizar
	password = strings.TrimSpace(password)

	// Verificar comprimento mínimo
	if len(password) < policy.MinLength {
		return apperrors.NewValidationError(
			fmt.Sprintf("senha deve ter pelo menos %d caracteres", policy.MinLength))
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
		charCount  = make(map[rune]int)
	)

	for _, char := range password {
		// Contar caracteres únicos
		charCount[char]++

		// Verificar tipos de caracteres
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

		// Verificar repetições excessivas
		if charCount[char] > policy.MaxRepeatedChars {
			return apperrors.NewValidationError(
				fmt.Sprintf("senha não pode ter mais que %d caracteres repetidos", policy.MaxRepeatedChars))
		}
	}

	// Verificar caracteres únicos
	if len(charCount) < policy.MinUniqueChars {
		return apperrors.NewValidationError(
			fmt.Sprintf("senha deve ter pelo menos %d caracteres únicos", policy.MinUniqueChars))
	}

	// Verificar requisitos de tipos de caracteres
	if policy.RequireUppercase && !hasUpper {
		return apperrors.NewValidationError("senha deve conter pelo menos uma letra maiúscula")
	}
	if policy.RequireLowercase && !hasLower {
		return apperrors.NewValidationError("senha deve conter pelo menos uma letra minúscula")
	}
	if policy.RequireNumbers && !hasNumber {
		return apperrors.NewValidationError("senha deve conter pelo menos um número")
	}
	if policy.RequireSpecial && !hasSpecial {
		return apperrors.NewValidationError("senha deve conter pelo menos um caractere especial")
	}

	// Verificar palavras proibidas
	passwordLower := strings.ToLower(password)
	for _, word := range policy.DisallowedWords {
		if strings.Contains(passwordLower, strings.ToLower(word)) {
			return apperrors.NewValidationError("senha contém uma sequência de caracteres proibida")
		}
	}

	return nil
}

// SanitizeString limpa uma string de caracteres potencialmente perigosos
func SanitizeString(input string) string {
	// Remover espaços extras
	input = strings.TrimSpace(input)

	// Remover caracteres de controle
	input = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(input, "")

	// Remover tags HTML
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")

	// Remover sequências de escape SQL comuns
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "--", "")

	return input
}
