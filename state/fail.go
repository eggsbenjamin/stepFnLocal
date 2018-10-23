package state

type FailDefinition struct {
	BaseDefinition
	Error string `json:"Error"`
	Cause string `json:"Cause"`
}

func (FailDefinition) Type() string {
	return FailStateType
}

func (s FailDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if err := s.BaseDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if s.Error == "" {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Error", "",
		))
	}

	if s.Cause == "" {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Cause", "",
		))
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}
