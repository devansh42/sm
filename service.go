//@author Devansh Gupta
//This file contains code for service manager

package sm

import "sync"

//AddService, adds a new service in the service queue
//There will be no effect if process is already started
func (s *ServiceManager) AddService(f func()) bool {
	if !s.running {
		s.list = append(s.list, f)
	}
	return !s.running //signals if process is online or not
}

//start, starts every service, atmost one time
func (s *ServiceManager) Start() {

	s.starter.Do(func() { //start assures that every process/service start atmost one time
		for _, v := range s.list {
			go v() //Running every service in its seperate go routine
		}
		s.running = true //setting state to running state
	})
}

func (s ServiceManager) Count() int {
	return len(s.list)
}

func NewServiceManager() *ServiceManager {
	sm := new(ServiceManager)
	sm.starter = new(sync.Once)
	return sm
}

//s, is the default service manager
var s = NewServiceManager()

//ServiceManager, it manages services accross the notification service
type ServiceManager struct {
	list    []func()
	starter *sync.Once
	//running, is a state variable, which says whether service in online or not
	running bool
}
