package state

// Typer defines the typer interface which all state definitions must implement
type Typer interface {
	Type() string
}

// Validator defines the validator interface which all state definitions must implement
type Validator interface {
	Validate() error
}

type Nexter interface {
	Next() string
}

type Transitioner interface {
	Nexter
	End() bool
}

type InputPather interface {
	InputPath() JSONPathExp
}

type OutputPather interface {
	OutputPath() JSONPathExp
}

type IOPather interface {
	InputPather
	OutputPather
}

type ResultPather interface {
	ResultPath() JSONPathExp
}

// Definition defines the definition interface which all state definitions must implement
type Definition interface {
	Typer
	Validator
}

// BaseDefinition represents an AWS states language state. It contains fields that can appear in all state types.
type BaseDefinition struct {
	StateType    string `json:"Type"`
	StateComment string `json:"Comment"`
}

func (s BaseDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if s.StateType == "" {
		validationErrs = append(validationErrs, NewValidationError(MissingRequiredFieldErrType, "Type", ""))
	} else {
		if _, ok := validTypes[s.StateType]; !ok {
			validationErrs = append(validationErrs, NewValidationError(InvalidValueErrType, "Type", s.StateType))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

type TransitionDefinition struct {
	NextState string `json:"Next"`
	EndState  bool   `json:"End"`
}

func (t TransitionDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if t.NextState == "" && t.EndState != true {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Next/End:true", "",
		))
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (t TransitionDefinition) Next() string {
	return t.NextState
}

func (t TransitionDefinition) End() bool {
	return t.EndState
}

type IODefinition struct {
	InputPathExp  JSONPathExp `json:"InputPath"`
	OutputPathExp JSONPathExp `json:"OutputPath"`
}

func (i IODefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if i.InputPathExp != "" {
		if err := i.InputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"InputPath", string(i.InputPathExp),
			))
		}
	}

	if i.OutputPathExp != "" {
		if err := i.OutputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"OutputPath", string(i.OutputPathExp),
			))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (i IODefinition) InputPath() JSONPathExp {
	return i.InputPathExp
}

func (i IODefinition) OutputPath() JSONPathExp {
	return i.OutputPathExp
}

type ResultDefinition struct {
	ResultPathExp JSONPathExp `json:"ResultPath"`
}

func (r ResultDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if r.ResultPathExp != "" {
		if err := r.ResultPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"ResultPath", string(r.ResultPathExp),
			))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (r ResultDefinition) ResultPath() JSONPathExp {
	return r.ResultPathExp
}
