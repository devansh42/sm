package sm

import "sync"

type basicServiceManager struct {
	starter *sync.Once
	//running, is a state variable, which says whether service in online or not
	running bool
	list    []Service
}

func newBasicServiceManager() *basicServiceManager {
	x := new(basicServiceManager)
	x.starter = new(sync.Once)
	return x
}

func (b basicServiceManager) Count() int {
	return len(b.list)
}

func (b *basicServiceManager) AddService(serv Service) {
	b.list = append(b.list, serv)
}

type GenericServiceManager interface {
	//Start, starts a service listed in service manager
	Start() error
	//Count, returns the no of services being handled by Service Manager
	Count() int
	//AddService, adds the service object to the list of service
	AddService(Service)
}

//Service, is the service object of the service being executed
type Service struct {
	//Executer is the underlying executing function
	Executer func()
	//Name is the name of service being used
	Name string
}
