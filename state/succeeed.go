package state

type SucceedDefinition struct {
	BaseDefinition
}

func (SucceedDefinition) Type() string {
	return SucceedStateType
}

func (s SucceedDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if err := s.BaseDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}
