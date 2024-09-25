package core

type Validator struct {
	result bool
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) IsError() bool {
	if v.result == false {
		return true
	}
	return false
}

// Port
// Check if the port number is valid
func (v *Validator) Port(port int) {
	if port == 0 || port < 0 || port > 65535 {
		v.result = false
	} else {
		v.result = true
	}
}
