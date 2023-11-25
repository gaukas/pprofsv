package pprofsv_test

import (
	"os"
	"testing"

	"github.com/gaukas/pprofsv"
	"github.com/google/pprof/profile"
)

func TestVerifierReachable(t *testing.T) {
	file, err := os.Open("testdata/pprof.profile")
	if err != nil {
		t.Fatal(err)
	}

	pprof, err := profile.Parse(file)
	if err != nil {
		t.Fatal(err)
	}

	profile := pprofsv.NewProfile(pprof)
	if profile == nil {
		t.Fatal("profile is nil")
	}

	verifier, err := profile.Verifier("dummy")
	if err != nil {
		t.Fatal(err)
	}

	vBranch, err := verifier.SubVerifier("")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Branch", func(t *testing.T) {
		testVerifierBranchReachable(t, vBranch)
	})

	vDeep, err := verifier.SubVerifier("")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Deep", func(t *testing.T) {
		testVerifierDeepReachable(t, vDeep)
	})

	vLoop, err := verifier.SubVerifier("")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Loop", func(t *testing.T) {
		testVerifierLoopReachable(t, vLoop)
	})

	vMulti, err := verifier.SubVerifier("")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Multi", func(t *testing.T) {
		testVerifierMultiReachable(t, vMulti)
	})

	vRecursive, err := verifier.SubVerifier("")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Recursive", func(t *testing.T) {
		testVerifierRecursive(t, vRecursive)
	})
}

func testVerifierBranchReachable(t *testing.T, v *pprofsv.Verifier) {
	callStacks := v.Callstack()
	if len(callStacks) == 0 {
		t.Fatal("no callstacks found")
	}

	v.SetFunctionPrefix("github.com/gaukas/pprofsv/dummy.(*Dummy).")

	if !v.Reachable("BranchFunc", "branchAinner") {
		t.Errorf("BranchFunc -> branchA should be reachable")
	}
	if !v.Next("BranchFunc", "branchA") {
		t.Errorf("BranchFunc -> branchA should be next")
	}

	if !v.Reachable("BranchFunc", "branchBinner") {
		t.Errorf("BranchFunc -> branchB should be reachable")
	}
	if !v.Next("BranchFunc", "branchB") {
		t.Errorf("BranchFunc -> branchB should be next")
	}

	if v.Reachable("branchA", "branchA") {
		t.Errorf("branchA -> branchA should not be reachable")
	}
	if v.Reachable("branchA", "branchB") {
		t.Errorf("branchA -> branchB should not be reachable")
	}
	if v.Reachable("branchA", "branchBinner") {
		t.Errorf("branchA -> branchBinner should not be reachable")
	}

	if v.Reachable("branchB", "branchB") {
		t.Errorf("branchB -> branchB should not be reachable")
	}
	if v.Reachable("branchB", "branchA") {
		t.Errorf("branchB -> branchA should not be reachable")
	}
	if v.Reachable("branchB", "branchAinner") {
		t.Errorf("branchB -> branchAinner should not be reachable")
	}
}

func testVerifierDeepReachable(t *testing.T, v *pprofsv.Verifier) {
	callStacks := v.Callstack()
	if len(callStacks) == 0 {
		t.Fatal("no callstacks found")
	}

	v.SetFunctionPrefix("github.com/gaukas/pprofsv/dummy.(*Dummy).")

	// top to bottom must be reachable
	if !v.Reachable("DeepFunc", "deepFuncLv5") {
		t.Errorf("DeepFunc -> deepFuncLv5 should be reachable")
	}

	// no backward or self-loop
	if v.Reachable("deepFuncLv5", "deepFuncLv5") {
		t.Errorf("deepFuncLv5 -> deepFuncLv5 should not be reachable")
	}
	if v.Reachable("deepFuncLv5", "deepFuncLv4") {
		t.Errorf("deepFuncLv5 -> deepFuncLv4 should not be reachable")
	}
	if v.Reachable("deepFuncLv5", "deepFuncLv3") {
		t.Errorf("deepFuncLv5 -> deepFuncLv3 should not be reachable")
	}
	if v.Reachable("deepFuncLv5", "deepFuncLv2") {
		t.Errorf("deepFuncLv5 -> deepFuncLv2 should not be reachable")
	}
	if v.Reachable("deepFuncLv5", "deepFuncLv1") {
		t.Errorf("deepFuncLv5 -> deepFuncLv1 should not be reachable")
	}
}

func testVerifierLoopReachable(t *testing.T, v *pprofsv.Verifier) {
	callStacks := v.Callstack()
	if len(callStacks) == 0 {
		t.Fatal("no callstacks found")
	}

	v.SetFunctionPrefix("github.com/gaukas/pprofsv/dummy.(*Dummy).")

	// outer to inner must be reachable
	if !v.Reachable("LoopFunc", "loopFuncInner") {
		t.Errorf("LoopFunc -> loopFuncInner should be reachable")
	}

	// looping doesn't really create self-loop in control flow graph
	if v.Reachable("LoopFunc", "LoopFunc") {
		t.Errorf("LoopFunc -> LoopFunc should not be reachable")
	}

	if v.Reachable("loopFuncInner", "loopFuncInner") {
		t.Errorf("loopFuncInner -> loopFuncInner should not be reachable")
	}

	// no backward
	if v.Reachable("loopFuncInner", "LoopFunc") {
		t.Errorf("loopFuncInner -> loopFunc should not be reachable")
	}
}

func testVerifierMultiReachable(t *testing.T, v *pprofsv.Verifier) {
	callStacks := v.Callstack()
	if len(callStacks) == 0 {
		t.Fatal("no callstacks found")
	}

	v.SetFunctionPrefix("github.com/gaukas/pprofsv/dummy.(*Dummy).")

	if !v.Reachable("MultiFunc", "multiFuncA") {
		t.Errorf("MultiFunc -> multiFuncA should be reachable")
	}
	if !v.Next("MultiFunc", "multiFuncA") {
		t.Errorf("MultiFunc -> multiFuncA should be next")
	}

	if !v.Reachable("MultiFunc", "multiFuncB") {
		t.Errorf("MultiFunc -> multiFuncB should be reachable")
	}
	if !v.Next("MultiFunc", "multiFuncB") {
		t.Errorf("MultiFunc -> multiFuncB should be next")
	}

	if !v.Reachable("MultiFunc", "multiFuncC") {
		t.Errorf("MultiFunc -> multiFuncC should be reachable")
	}
	if !v.Next("MultiFunc", "multiFuncC") {
		t.Errorf("MultiFunc -> multiFuncC should be next")
	}

	if !v.Reachable("MultiFunc", "multiFuncD") {
		t.Errorf("MultiFunc -> multiFuncD should be reachable")
	}
	if !v.Next("MultiFunc", "multiFuncD") {
		t.Errorf("MultiFunc -> multiFuncD should be next")
	}

	// no neighboring path
	if v.Reachable("multiFuncA", "multiFuncB") {
		t.Errorf("multiFuncA -> multiFuncB should not be reachable")
	}

	if v.Reachable("multiFuncB", "multiFuncC") {
		t.Errorf("multiFuncB -> multiFuncC should not be reachable")
	}

	if v.Reachable("multiFuncC", "multiFuncD") {
		t.Errorf("multiFuncC -> multiFuncD should not be reachable")
	}
}

func testVerifierRecursive(t *testing.T, v *pprofsv.Verifier) {
	callStacks := v.Callstack()
	if len(callStacks) == 0 {
		t.Fatal("no callstacks found")
	}

	v.SetFunctionPrefix("github.com/gaukas/pprofsv/dummy.(*Dummy).")

	if !v.Reachable("RecursiveFunc", "recursiveFuncInnerA") {
		t.Errorf("RecursiveFunc -> recursiveFuncInnerA should be reachable")
	}
	if !v.Reachable("RecursiveFunc", "recursiveFuncInnerB") {
		t.Errorf("RecursiveFunc -> recursiveFuncInnerB should be reachable")
	}

	// self-loop should be found (via A->B->A->B...)
	if !v.Reachable("recursiveFuncInnerA", "recursiveFuncInnerA") {
		t.Errorf("recursiveFuncInnerA -> recursiveFuncInnerA should be reachable")
	}
	if !v.Reachable("recursiveFuncInnerB", "recursiveFuncInnerB") {
		t.Errorf("recursiveFuncInnerB -> recursiveFuncInnerB should be reachable")
	}

	// switch between A and B
	if !v.Next("recursiveFuncInnerA", "recursiveFuncInnerB") {
		t.Errorf("recursiveFuncInnerA -> recursiveFuncInnerB should be next")
	}
	if !v.Next("recursiveFuncInnerB", "recursiveFuncInnerA") {
		t.Errorf("recursiveFuncInnerB -> recursiveFuncInnerA should be next")
	}

	// B is not the terminating condition
	if v.Next("recursiveFuncInnerB", "final") {
		t.Errorf("recursiveFuncInnerB -> final should not be next")
	}

	// A is the terminating condition
	if !v.Next("recursiveFuncInnerA", "final") {
		t.Errorf("recursiveFuncInnerA -> final should be next")
	}
}
