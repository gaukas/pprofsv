package pprofsv_test

import (
	"testing"

	"github.com/gaukas/pprofsv"
)

func TestPath(t *testing.T) {
	t.Run("Direct", testPathDirect)
	t.Run("Indirect", testPathIndirect)
}

func testPathDirect(t *testing.T) {
	p := pprofsv.NewPath(5)

	p.Set(0, 1)
	p.Set(1, 2)
	p.Set(2, 3)
	p.Set(3, 4)

	// All direct routes must be found.
	if !(p.HasDirectPath(0, 1) && p.HasPath(0, 1)) {
		t.Errorf("direct route 0->1 not found")
	}

	if !(p.HasDirectPath(1, 2) && p.HasPath(1, 2)) {
		t.Errorf("direct route 1->2 not found")
	}

	if !(p.HasDirectPath(2, 3) && p.HasPath(2, 3)) {
		t.Errorf("direct route 2->3 not found")
	}

	if !(p.HasDirectPath(3, 4) && p.HasPath(3, 4)) {
		t.Errorf("direct route 3->4 not found")
	}

	// no indirect routes can be found with HasDirectPath.
	if p.HasDirectPath(0, 2) {
		t.Errorf("direct route 0->2 reported unexpectedly")
	}

	if p.HasDirectPath(0, 3) {
		t.Errorf("direct route 0->3 reported unexpectedly")
	}

	if p.HasDirectPath(0, 4) {
		t.Errorf("direct route 0->4 reported unexpectedly")
	}

	if p.HasDirectPath(1, 3) {
		t.Errorf("direct route 1->3 reported unexpectedly")
	}

	if p.HasDirectPath(1, 4) {
		t.Errorf("direct route 1->4 reported unexpectedly")
	}

	if p.HasDirectPath(2, 4) {
		t.Errorf("direct route 2->4 reported unexpectedly")
	}

	// routes have direction. 0->1 is not 1->0.
	if p.HasDirectPath(1, 0) || p.HasPath(1, 0) {
		t.Errorf("direct route 1->0 found unexpectedly")
	}

	if p.HasDirectPath(2, 1) || p.HasPath(2, 1) {
		t.Errorf("direct route 2->1 found unexpectedly")
	}

	if p.HasDirectPath(3, 2) || p.HasPath(3, 2) {
		t.Errorf("direct route 3->2 found unexpectedly")
	}

	if p.HasDirectPath(4, 3) || p.HasPath(4, 3) {
		t.Errorf("direct route 4->3 found unexpectedly")
	}
}

func testPathIndirect(t *testing.T) {
	p := pprofsv.NewPath(5)

	p.Set(0, 1)
	p.Set(1, 2)
	p.Set(2, 3)
	p.Set(3, 4)

	// 0->2 is not a direct route, but it is an indirect route.
	if !p.HasPath(0, 2) {
		t.Errorf("indirect route 0->2 not found")
	}

	// 0->3 is not a direct route, but it is an indirect route.
	if !p.HasPath(0, 3) {
		t.Errorf("indirect route 0->3 not found")
	}

	// 0->4 is not a direct route, but it is an indirect route.
	if !p.HasPath(0, 4) {
		t.Errorf("indirect route 0->4 not found")
	}

	// 1->3 is not a direct route, but it is an indirect route.
	if !p.HasPath(1, 3) {
		t.Errorf("indirect route 1->3 not found")
	}

	// 1->4 is not a direct route, but it is an indirect route.
	if !p.HasPath(1, 4) {
		t.Errorf("indirect route 1->4 not found")
	}

	// 2->4 is not a direct route, but it is an indirect route.
	if !p.HasPath(2, 4) {
		t.Errorf("indirect route 2->4 not found")
	}
}
