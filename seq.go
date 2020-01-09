package sm

//SequentialServiceManager,
type SequentialServiceManager struct {
	*basicServiceManager
}

//Starts the services in Sequential order
func (s *SequentialServiceManager) Start() {
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
}
