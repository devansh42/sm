package sm

import (
	"testing"
	"time"
)

func TestServiceManager(t *testing.T) {
	x := NewServiceManager()
	f := func(i rune) func() {
		return func() {
			t.Log("This is Service : ", string(i))
		}
	}
	x.AddService(f('D'))
	x.AddService(f('A'))
	x.AddService(f('B'))
	x.AddService(f('C'))
	x.AddService(f('E'))
	x.Start()
	time.Sleep(time.Second * 5)
}
