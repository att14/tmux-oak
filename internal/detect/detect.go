package detect

type Agent struct {
	Name string
	Icon string
}

type Detector interface {
	Name() string
	Icon() string
	Scan(panePIDs map[int]bool) (map[int]Agent, error)
}

type Registry struct {
	detectors []Detector
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(d Detector) {
	r.detectors = append(r.detectors, d)
}

func (r *Registry) Scan(panePIDs map[int]bool) map[int]Agent {
	result := make(map[int]Agent)
	for _, d := range r.detectors {
		agents, err := d.Scan(panePIDs)
		if err != nil {
			continue
		}
		for pid, agent := range agents {
			result[pid] = agent
		}
	}
	return result
}
