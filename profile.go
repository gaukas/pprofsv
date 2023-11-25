package pprofsv

import "github.com/google/pprof/profile"

type Profile struct {
	functionNameMap map[string]uint64
	functionIdMap   map[uint64]string

	callStacks [][]uint64 // callStacks[i] is the call stack of sample i, created from chaining all locations in sample i.
}

func NewProfile(pprof *profile.Profile) *Profile {
	p := &Profile{
		functionNameMap: make(map[string]uint64),
		functionIdMap:   make(map[uint64]string),
		callStacks:      make([][]uint64, len(pprof.Sample)),
	}

	for _, function := range pprof.Function {
		p.functionNameMap[function.Name] = function.ID
		p.functionIdMap[function.ID] = function.Name
	}

	for i, sample := range pprof.Sample {
		callStack := make([]uint64, 0, len(sample.Location))
		for _, location := range sample.Location {
			for _, line := range location.Line {
				callStack = append(callStack, line.Function.ID)
			}
		}
		p.callStacks[i] = callStack
	}

	return p
}

// Verifier returns a new Verifier for functions matching a
// given regular expression.
func (p *Profile) Verifier(namePattern string) (*Verifier, error) {
	return NewVerifier(p, nil, namePattern)
}
