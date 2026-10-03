package installer

type StepRegistry struct {
	steps []Step
}

func (r *StepRegistry) Register(step Step) {
	r.steps = append(r.steps, step)
}

func (r *StepRegistry) GetAll() []Step {
	return r.steps
}