package errors

type ValidationError struct {
	Message string                 `json:"message,omitempty"`
	Fields  []ValidationErrorField `json:"field_name,omitempty"`
}

// Error implements error.
func (v *ValidationError) Error() string {
	panic("unimplemented")
}

type ValidationErrorField struct {
	Message string `json:"message,omitempty"`
	Field   string `json:"field_name,omitempty"`
}

func NewRequiredFieldsErrorList(fields []string) *ValidationError {
	message := "Campo Obrigatório"
	fieldErr := []ValidationErrorField{}
	for i := 0; i < len(fields); i++ {
		fieldErr = append(fieldErr, ValidationErrorField{
			Message: message,
			Field:   fields[i],
		})
	}
	return &ValidationError{
		Fields:  fieldErr,
		Message: message,
	}
}
