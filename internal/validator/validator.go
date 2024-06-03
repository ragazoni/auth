package validator

type Validator struct {
	BirthDate *BirthDate
	Document  *Document
}

func New() *Validator {

	validator := new(Validator)

	validator.BirthDate = newBirthDate()
	validator.Document = newDocument()

	return validator
}
