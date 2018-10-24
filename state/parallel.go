package state

type ParallelDefinition struct {
	BaseDefinition
	TransitionDefinition
	IOPathDefinition
	ResultPathDefinition
	Branches []MachineDefinition `json:"Branches"`
}

func (p ParallelDefinition) Type() string {
	return ParallelStateType
}

func (p ParallelDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if err := p.BaseDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := p.TransitionDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := p.IOPathDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := p.ResultPathDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if len(p.Branches) == 0 {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Branches", "Is empty",
		))
	}

	if p.InputPathExp != "" {
		if err := p.InputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"InputPath", string(p.InputPathExp),
			))
		}
	}

	if p.OutputPathExp != "" {
		if err := p.OutputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"OutputPath", string(p.OutputPathExp),
			))
		}
	}

	if p.ResultPathExp != "" {
		if err := p.ResultPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"ResultPath", string(p.ResultPathExp),
			))
		}
	}

	for _, branch := range p.Branches {
		if err := branch.Validate(); err != nil {
			validationErrs = append(validationErrs, err.(ValidationErrors)...)
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}
