package dummy

import (
	"crypto/rand"
	"time"
)

var defaultDummySleep time.Duration = 100 * time.Microsecond

type Dummy struct {
	sleepDuration time.Duration
}

func NewDummy() *Dummy {
	return &Dummy{sleepDuration: defaultDummySleep}
}

func (d *Dummy) SetSleep(sleep time.Duration) {
	d.sleepDuration = sleep
}

//go:noinline
func (d *Dummy) BranchFunc(A bool) {
	if A {
		d.branchA() // 50%
	} else {
		d.branchB() // 50%
	}
}

//go:noinline
func (d *Dummy) DeepFunc() {
	d.deepFuncLv1()
}

//go:noinline
func (d *Dummy) LoopFunc(n int) {
	for i := 0; i < n; i++ {
		d.loopFuncInner()
	}
}

//go:noinline
func (d *Dummy) MultiFunc() {
	d.multiFuncA()
	d.multiFuncB()
	d.multiFuncC()
	d.multiFuncD()
}

//go:noinline
func (d *Dummy) RecursiveFunc(n int) {
	if n == 0 {
		return
	}
	d.recursiveFuncInnerA(n)
}

//go:noinline
func (d *Dummy) branchA() {
	d.sleep()
	d.branchAinner()
}

//go:noinline
func (d *Dummy) branchAinner() {
	d.final()
}

//go:noinline
func (d *Dummy) branchB() {
	d.sleep()
	d.branchBinner()
}

//go:noinline
func (d *Dummy) branchBinner() {
	d.final()
}

//go:noinline
func (d *Dummy) deepFuncLv1() {
	d.sleep()
	d.deepFuncLv2()
}

//go:noinline
func (d *Dummy) deepFuncLv2() {
	d.sleep()
	d.deepFuncLv3()
}

//go:noinline
func (d *Dummy) deepFuncLv3() {
	d.sleep()
	d.deepFuncLv4()
}

//go:noinline
func (d *Dummy) deepFuncLv4() {
	d.sleep()
	d.deepFuncLv5()
}

//go:noinline
func (d *Dummy) deepFuncLv5() {
	d.final()
}

//go:noinline
func (d *Dummy) multiFuncA() {
	d.final()
}

//go:noinline
func (d *Dummy) multiFuncB() {
	d.final()
}

//go:noinline
func (d *Dummy) multiFuncC() {
	d.final()
}

//go:noinline
func (d *Dummy) multiFuncD() {
	d.final()
}

//go:noinline
func (d *Dummy) loopFuncInner() {
	d.final()
}

//go:noinline
func (d *Dummy) recursiveFuncInnerA(n int) {
	d.sleep()
	if n == 0 {
		d.final()
		return
	}
	d.recursiveFuncInnerB(n)
}

//go:noinline
func (d *Dummy) recursiveFuncInnerB(n int) {
	d.sleep()
	if n == 0 {
		d.final()
		return
	}
	d.recursiveFuncInnerA(n - 1)
}

//go:noinline
func (d *Dummy) final() {
	d.sleep()
	d.alloc()
}

//go:inline
func (d *Dummy) sleep() {
	// time.Sleep(d.sleepDuration)
}

//go:noinline
func (d *Dummy) alloc() {
	var buf []byte = make([]byte, 8)
	rand.Reader.Read(buf)
}
