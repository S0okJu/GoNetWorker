package core

// Validator 유효성 검사를 위한 구조체
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

// Port 포트 범위 검사
func (v *Validator) Port(port int) {
	if port == 0 || port < 0 || port > 65535 {
		v.result = false
	} else {
		v.result = true
	}
}
