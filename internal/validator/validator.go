package validator

type Validator struct {
	Error map[string]string
}

func New() *Validator {
	return &Validator{
		Error: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Error) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Error[key]; !exists {
		v.Error[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
