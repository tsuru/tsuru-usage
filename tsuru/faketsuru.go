package tsuru

type FakeTsuruAPI struct {
	Apps      []App
	Instances map[string][]ServiceInstance
	Nodes     []Node
}

func (f *FakeTsuruAPI) ListApps() ([]App, error) {
	return f.Apps, nil
}

func (f *FakeTsuruAPI) ListServiceInstances(service string) ([]ServiceInstance, error) {
	return f.Instances[service], nil
}

func (f *FakeTsuruAPI) ListNodes() ([]Node, error) {
	return f.Nodes, nil
}
