package sm

import (
	"errors"
)

type DependencyInjecter interface {
	Sequence(string) ([]string, error)
}

//DependentServiceManager is the manager dependent services
type DependentServiceManager struct {
	*basicServiceManager
	deps   DependencyInjecter
	target string
}

//NewDependentServiceManager, returns a new instance of DependentServiceManager
func NewDependentServiceManager() *DependentServiceManager {
	x := new(DependentServiceManager)
	x.basicServiceManager = newBasicServiceManager()
	return x
}

//SetTarget, sets the name of target Service to launch
func (d *DependentServiceManager) SetTarget(sourceService string) {
	d.target = sourceService
}

//Start, starts Dependent service are per dependency injector
func (d *DependentServiceManager) Start() (err error) {
	defer func() { //Recover from panic
		if errr := recover(); errr != nil {
			err = errr.(error)
		}
	}()

	d.starter.Do(func() {

		var servs = make(map[string]func())
		for _, v := range d.list {
			servs[v.Name] = v.Executer
		}
		seq, err := d.deps.Sequence(d.target)
		if err != nil {
			panic(err)
		}

		s := NewSequentialServiceManager()
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
	return
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
func (d TopologicalDependencyIntjecter) Sequence(source string) (sstr []string, err error) {

	//Checking for cycles

	defer func() {
		if errr := recover(); errr != nil {
			err = errr.(error) //This will be a cycle error most probably
			if err != CyclicDependencyError {
				panic(err)
			}
		}
	}()

	count := 0
	for range d.adj {
		count++
	}
	stack := make([]string, count)
	si := -1
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var dfs func(string)
	dfs = func(node string) {

		if node == "" {
			return
		}

		if isv, ok := recStack[node]; ok && isv {
			panic(CyclicDependencyError)
		}

		if isv, ok := visited[node]; ok && isv {
			return
		}

		visited[node] = true
		recStack[node] = true
		for _, v := range d.adj[node] {

			dfs(v)
		}
		si++
		stack[si] = node
		recStack[node] = false
	}
	dfs(source)

	sstr = stack[:si+1]
	return
}

//CyclicDependencyError, is the error represting cyclic dependency
var CyclicDependencyError = errors.New("Cycle detected")
