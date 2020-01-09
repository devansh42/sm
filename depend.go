package sm

type DependencyInjecter interface {
	Sequence(string) []string
}

//DependentServiceManager is the manager dependent services
type DependentServiceManager struct {
	*basicServiceManager
	deps   DependencyInjecter
	source string
}

func (d *DependentServiceManager) SetSource(s string) {
	d.source = s
}

//DependentServiceManager, is the service manager for dependent services
func (d *DependentServiceManager) Start() {
	d.starter.Do(func() {

		var servs map[string]func()
		for _, v := range d.list {
			servs[v.Name] = v.Executer
		}
		seq := d.deps.Sequence(d.source)
		s := new(SequentialServiceManager)
		for _, v := range seq {
			f, ok := servs[v]
			if !ok { //Checking for occurence
				continue
			}
			s.AddService(Service{Executer: f})
		}
		s.Start()
		for !s.running { //Waiting for Sequential start of services
		}
		d.running = true
	})
}

func (d *DependentServiceManager) SetDependencyInjecter(deps DependencyInjecter) {
	d.deps = deps
}

//TopologicalDependencyIntjecter, is the dependency injecter for services
type TopologicalDependencyIntjecter struct {
	adj map[string][]string
}

//NewTopologicalDependencyInjecter,returns a new dependency injecter
func NewTopologicalDependencyInjecter() *TopologicalDependencyIntjecter {
	x := new(TopologicalDependencyIntjecter)
	x.adj = make(map[string][]string)
	x.adj[""] = []string{} //For terminal value
	return x
}

//AddDependency, adds a dependency of a dependent
func (d *TopologicalDependencyIntjecter) AddDependency(dependent, dependency string) {

	d.adj[dependent] = append(d.adj[dependent], dependency)
}

//Sequence, is the topological sequence of our dependency graph
func (d TopologicalDependencyIntjecter) Sequence(source string) []string {
	count := 0
	for range d.adj {
		count++
	}
	head := 0
	stack := make([]string, count)
	visited := make(map[string]bool)
	var dfs func(string)
	dfs = func(node string) {
		if node == "" {
			return
		}

		for _, v := range d.adj[node] {
			if isv, ok := visited[v]; ok && isv {
				continue
			}
			dfs(v)
		}
		stack[head] = node
		head++
	}
	dfs(source)

	return stack[:head]
}
