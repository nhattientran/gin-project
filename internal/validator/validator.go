package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError add new message if not have in map error
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds error message to map if the validation is not ok
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In check if value in the list
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Unique check if value is unique in list
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		if uniqueValues[value] {
			return false
		}
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
