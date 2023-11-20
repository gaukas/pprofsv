package dummy

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"math"
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

func (d *Dummy) sleep() {
	time.Sleep(d.sleepDuration)
}

func (d *Dummy) BranchFunc() {
	var randBuf []byte = make([]byte, 1)
	rand.Reader.Read(randBuf)

	if randBuf[0] > 128 {
		d.branchA() // 50%
	} else {
		d.branchB() // 50%
	}
}

func (d *Dummy) DeepFunc() {
	d.deepFuncLv1()
}

func (d *Dummy) LoopFunc(n int) {
	for i := 0; i < n; i++ {
		d.loopFuncInner()
	}
}

func (d *Dummy) RecursiveFunc(n int) {
	if n == 0 {
		return
	}
	d.recursiveFuncInner(n)
}

func (d *Dummy) branchA() {
	d.sleep()
	slog.Debug("hit branchA")
	var buf []byte = make([]byte, 16)
	rand.Reader.Read(buf)
	d.branchAinner()
}

func (d *Dummy) branchAinner() {
	d.sleep()
	slog.Debug("hit branchA inner")
	var bufInner []byte = make([]byte, 32)
	rand.Reader.Read(bufInner)
}

func (d *Dummy) branchB() {
	d.sleep()
	slog.Debug("hit branchB")
	var buf []byte = make([]byte, 32)
	rand.Reader.Read(buf)
	d.branchBinner()
}

func (d *Dummy) branchBinner() {
	d.sleep()
	slog.Debug("hit branchB inner")
	var bufInner []byte = make([]byte, 64)
	rand.Reader.Read(bufInner)
}

func (d *Dummy) deepFuncLv1() {
	d.sleep()
	slog.Debug("deepFuncLv1")
	var buf []byte = make([]byte, 2)
	rand.Reader.Read(buf)
	d.deepFuncLv2()
}

func (d *Dummy) deepFuncLv2() {
	d.sleep()
	slog.Debug("deepFuncLv2")
	var buf []byte = make([]byte, 4)
	rand.Reader.Read(buf)
	d.deepFuncLv3()
}

func (d *Dummy) deepFuncLv3() {
	d.sleep()
	slog.Debug("deepFuncLv3")
	var buf []byte = make([]byte, 8)
	rand.Reader.Read(buf)
	d.deepFuncLv4()
}

func (d *Dummy) deepFuncLv4() {
	d.sleep()
	slog.Debug("deepFuncLv4")
	var buf []byte = make([]byte, 16)
	rand.Reader.Read(buf)
	d.deepFuncLv5()
}

func (d *Dummy) deepFuncLv5() {
	d.sleep()
	slog.Debug("deepFuncLv5")
	var buf []byte = make([]byte, 32)
	rand.Reader.Read(buf)
	return
}

func (d *Dummy) loopFuncInner() {
	d.sleep()
	slog.Debug("hit loopFunc Inner")
	var buf []byte = make([]byte, 8)
	rand.Reader.Read(buf)
}

func (d *Dummy) recursiveFuncInner(n int) {
	if n == 0 {
		return
	}

	d.sleep()
	slog.Debug(fmt.Sprintf("hit recursiveFunc Inner n=%d", n))
	var buf []byte = make([]byte, int(math.Pow(2, float64(n))))
	rand.Reader.Read(buf)
	d.recursiveFuncInner(n - 1)
}
