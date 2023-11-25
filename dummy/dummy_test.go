package dummy_test

import (
	"testing"

	"github.com/gaukas/pprofsv/dummy"
)

func TestDummyBranchFunc(t *testing.T) {
	d := dummy.NewDummy()

	for i := 0; i < 10; i++ {
		d.BranchFunc(i%2 == 0)
	}
}

func TestDummyDeepFunc(t *testing.T) {
	d := dummy.NewDummy()

	d.DeepFunc()
}

func TestDummyInlineFunc(t *testing.T) {
	d := dummy.NewDummy()

	d.MultiFunc()
}

func TestDummyLoopFunc(t *testing.T) {
	d := dummy.NewDummy()

	d.LoopFunc(10)
}

func TestDummyRecursiveFunc(t *testing.T) {
	d := dummy.NewDummy()

	d.RecursiveFunc(10)
}

func BenchmarkDummyBranchFunc(b *testing.B) {
	d := dummy.NewDummy()

	for i := 0; i < b.N; i++ {
		d.BranchFunc(i%2 == 0)
	}
}

func BenchmarkDummyDeepFunc(b *testing.B) {
	d := dummy.NewDummy()

	for i := 0; i < b.N; i++ {
		d.DeepFunc()
	}
}

func BenchmarkDummyInlineFunc(b *testing.B) {
	d := dummy.NewDummy()

	for i := 0; i < b.N; i++ {
		d.MultiFunc()
	}
}

func BenchmarkDummyLoopFunc(b *testing.B) {
	d := dummy.NewDummy()

	for i := 0; i < b.N; i++ {
		d.LoopFunc(10)
	}
}

func BenchmarkDummyRecursiveFunc(b *testing.B) {
	d := dummy.NewDummy()

	for i := 0; i < b.N; i++ {
		d.RecursiveFunc(10)
	}
}
