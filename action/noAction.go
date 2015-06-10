package action

type NoAction struct{}

func (p *NoAction) Name() string {
	return "noAction"
}

func (p *NoAction) Initialize() error {
	return nil
}

func (p *NoAction) Run(config *GruActionConfig) error {
	return nil
}
