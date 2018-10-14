package state

type State interface {
	Run(input []byte) ([]byte, error)
	Next() string
}

type TaskState struct {
	def TaskStateDefinition
}

func (t TaskState) Run(input []byte) ([]byte, error) {
	return nil, nil
}

func (t TaskState) Next() string {
	return def.Next()
}
