package sm

//SequentialServiceManager,
type SequentialServiceManager struct {
	*basicServiceManager
}

//NewSequentialServiceManager, returns a new instance of Sequential Service Mananger
func NewSequentialServiceManager() *SequentialServiceManager {
	x := new(SequentialServiceManager)
	x.basicServiceManager = newBasicServiceManager()
	return x
}

//Starts the services in Sequential order
func (s *SequentialServiceManager) Start() (err error) {
	ch := make(chan Service)
	completed := make(chan struct{})
	s.starter.Do(func() {

		go func() {
			defer func() { close(completed); s.running = true }()
			for serv := range ch {
				go serv.Executer()
				completed <- struct{}{}
			}

		}()

		for _, v := range s.list {
			ch <- v
			<-completed
		}
		close(ch)
	})
	return
}
