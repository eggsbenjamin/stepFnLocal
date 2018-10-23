package state

import "encoding/json"

type PassDefinition struct {
	BaseDefinition
	TransitionDefinition
	IOPathDefinition
	ResultPathDefinition
	Result json.RawMessage `json:"Result"`
}

func (PassDefinition) Type() string {
	return PassStateType
}

func (p PassDefinition) Validate() error {
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

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}
