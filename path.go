package pprofsv

import (
	"fmt"
	"sync"

	"github.com/crillab/gophersat/bf"
)

// Path holds a n-by-n matrix representing the path between nodes.
//
// Axioms:
//  1. Path.directPaths[i][i] being true means there is a self-loop at node i.
//  2. Path.directPaths[i][j] being true means there is a path from i to j.
//  3. For any i,j,k, if there is a path from i to j ([i][j] is true)
//     and a path from j to k ([j][k] is true), then there is a path
//     from i to k no matter whether [i][k] is true or not.
type Path struct {
	directPaths [][]bool
	allPaths    [][]bool

	rw *sync.RWMutex
}

// NewPath returns a new Path with size n.
func NewPath(n int) *Path {
	p := make([][]bool, n)
	for i := range p {
		p[i] = make([]bool, n)
	}

	allPaths := make([][]bool, n)
	for i := range allPaths {
		allPaths[i] = make([]bool, n)
	}

	return &Path{
		directPaths: p,
		allPaths:    allPaths,
		rw:          &sync.RWMutex{},
	}
}

// Set sets Path[i][j] to true, which means there is a DIRECT path from i to j.
func (p *Path) Set(i, j int) {
	p.rw.Lock()
	defer p.rw.Unlock()
	p.directPaths[i][j] = true
	p.allPaths[i][j] = true
}

// HasPath returns true if there is a
func (p *Path) HasPath(i, j int) bool {
	p.rw.RLock()
	defer p.rw.RUnlock()
	if p.allPaths[i][j] {
		return true
	}

	p.allPaths[i][j] = p.satCheckPath(i, j)
	return p.allPaths[i][j]
}

func (p *Path) HasDirectPath(i, j int) bool {
	p.rw.RLock()
	defer p.rw.RUnlock()
	return p.directPaths[i][j]
}

func (p *Path) satCheckPath(i, j int) bool {
	// SAT problem:
	//  1) For all a, b, if p.allPaths[a][b] is true, add constraint: R(a,b) == true.
	//  2) add constraint: (R(a,b) && R(b,c)) => R(a,c)
	//  3) add hypothesis: R(i,j) == false
	// if the hypothesis is not satisfiable, then there IS a path from i to j.
	const varFmt = "R(%d,%d)==%t"
	constraints := bf.True

	// either R(i,j)==true or R(i,j)==false
	constraints = bf.And(constraints, bf.Unique(fmt.Sprintf(varFmt, i, j, true), fmt.Sprintf(varFmt, i, j, false)))

	// 1) For all a, b, if p.allPaths[a][b] is true, add constraint: R(a,b) == true.
	for a := range p.allPaths {
		for b := range p.allPaths[a] {
			if p.allPaths[a][b] {
				constraints = bf.And(constraints, bf.Var(fmt.Sprintf(varFmt, a, b, true)))
			}
		}
	}

	// 2) add constraint: (R(a,b) && R(b,c)) => R(a,c)
	for a := range p.allPaths {
		for b := range p.allPaths {
			for c := range p.allPaths {
				constraints = bf.And(constraints, bf.Implies(bf.And(bf.Var(fmt.Sprintf(varFmt, a, b, true)), bf.Var(fmt.Sprintf(varFmt, b, c, true))), bf.Var(fmt.Sprintf(varFmt, a, c, true))))
			}
		}
	}

	// 3) add (negated) hypothesis: R(i,j) == false
	constraints = bf.And(constraints, bf.Not(bf.Var(fmt.Sprintf(varFmt, i, j, true))))

	// solve the SAT problem
	model := bf.Solve(constraints)
	return model == nil // nil -> unsatisfiable -> there is a path from i to j
}
