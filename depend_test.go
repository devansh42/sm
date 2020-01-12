package sm

import "testing"

import "time"

func TestTopologicalDependencieInjector(t *testing.T) {

	x := injecterWithOutCycle()
	ss, err := x.Sequence("a")
	t.Log(ss, err)
}

func TestTopologicalDependencieInjectorWithCycle(t *testing.T) {
	x := injecterWithCycle()
	ss, err := x.Sequence("a")
	t.Log(ss, err)

}

func TestDependentServiceManager(t *testing.T) {
	x := NewDependentServiceManager()
	x.SetDependencyInjecter(injecterWithOutCycle())
	x.SetTarget("a")
	f := func(i rune) Service {
		return Service{func() {
			t.Log("This is service : ", string(i))

		}, string(i)}
	}
	x.AddService(f('a'))
	x.AddService(f('b'))
	x.AddService(f('c'))
	x.AddService(f('d'))
	x.AddService(f('e'))
	x.Start()
	time.Sleep(time.Second * 5)

}

func TestDependentServiceManagerWithCycle(t *testing.T) {
	x := NewDependentServiceManager()
	x.SetDependencyInjecter(injecterWithCycle())
	x.SetTarget("a")
	f := func(i rune) Service {
		return Service{func() {
			t.Log("This is service : ", string(i))

		}, string(i)}
	}
	x.AddService(f('a'))
	x.AddService(f('b'))
	x.AddService(f('c'))

	err := x.Start()
	if err != nil {
		t.Log(err)
	}
	time.Sleep(time.Second * 5)

}

func injecterWithOutCycle() *TopologicalDependencyIntjecter {
	x := NewTopologicalDependencyInjecter()
	a, b, c, d, e := "a", "b", "c", "d", "e"
	x.AddDependency(a, b)
	x.AddDependency(a, c)
	x.AddDependency(c, d)
	x.AddDependency(b, d)
	x.AddDependency(d, e)
	return x
}

func injecterWithCycle() *TopologicalDependencyIntjecter {
	x := NewTopologicalDependencyInjecter()
	a, b, c, _, _ := "a", "b", "c", "d", "e"
	x.AddDependency(a, b)
	x.AddDependency(b, c)
	x.AddDependency(c, a)
	return x
}
