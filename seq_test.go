package sm

import (
	"testing"
	"time"
)

func TestSequentialServiceManager(t *testing.T) {
	x := NewSequentialServiceManager()
	f := func(a rune) func() {
		return func() {
			t.Log("This is Service : ", string(a))
		}
	}
	x.AddService(Service{f('D'), "D"})
	x.AddService(Service{f('A'), "A"})
	x.AddService(Service{f('B'), "B"})
	x.AddService(Service{f('E'), "E"})
	x.AddService(Service{f('C'), "C"})
	x.Start()

	//Waiting for output
	time.Sleep(time.Second * 5)
}
