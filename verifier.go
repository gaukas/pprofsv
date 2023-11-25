package pprofsv

import (
	"log"
	"regexp"
)

type Verifier struct {
	// callStacks is a subset of masterProfile.callStacks
	//
	// Each slice in callStacks is a call stack, represented by a
	// series of function IDs that are in masterProfile.functionIdMap.
	// It could be reused to create sub-verifiers.
	//
	// Note: it does not use the pseudoID.
	callStacks [][]uint64

	// path describes the reachability between functions.
	//
	// It uses pseudoID to represent functions in order to save memory.
	path *Path

	// functionIdPseudoMap is a map from real function ID to pseudoID.
	// It is used to convert real function ID to pseudoID, so that
	// the path can be built with a minimal memory footprint while
	// including only the functions that are interesting.
	functionIdPseudoMap map[uint64]uint64

	// masterProfile links back to the root profile where all functions
	// are included.
	masterProfile *Profile

	// prefix of the function name helps to reduce the length of
	// input per each verification request.
	functionPrefix string
}

// NewVerifier returns a new Verifier for functions matching a
// given regular expression. The call stack will be reduced to
// include only functions that match the name pattern.
//
// If baseCallStacks is nil, then the Verifier will build based on
// masterProfile.callStacks.
//
// If namePattern is empty, then the Verifier will use all functions
// in masterProfile. This may result in a very slow verification or
// even a memory overflow.
func NewVerifier(masterProfile *Profile, baseCallStacks [][]uint64, namePattern string) (*Verifier, error) {
	// filter call stacks
	var finalCallStacks [][]uint64
	var originalCallStacks [][]uint64
	if baseCallStacks == nil {
		originalCallStacks = masterProfile.callStacks
	} else {
		originalCallStacks = baseCallStacks
	}

	var interestingFunctionIds []uint64 // function IDs that match the name pattern
	if namePattern == "" {
		finalCallStacks = originalCallStacks
		interestingFunctionIds = make([]uint64, 0, len(masterProfile.functionNameMap))
		for f := range masterProfile.functionIdMap {
			interestingFunctionIds = append(interestingFunctionIds, f)
		}
	} else {
		finalCallStacks = make([][]uint64, 0, len(originalCallStacks))
		for name, function := range masterProfile.functionNameMap {
			// regex match
			if match, err := regexp.Match(namePattern, []byte(name)); match {
				// fmt.Printf("Matched: %s\n", name)
				interestingFunctionIds = append(interestingFunctionIds, function)
			} else if err != nil {
				return nil, err
			}
		}

		for _, callStack := range originalCallStacks {
			reducedCallStack := make([]uint64, 0, len(callStack))
		LOOP_FUNC_IN_CALLSTACK:
			for _, function := range callStack {
				for _, interestingFunction := range interestingFunctionIds {
					if function == interestingFunction {
						// fmt.Printf("Function %d is interesting\n", function)
						reducedCallStack = append(reducedCallStack, function)
						continue LOOP_FUNC_IN_CALLSTACK
					}
				}
			}
			if len(reducedCallStack) > 0 {
				finalCallStacks = append(finalCallStacks, reducedCallStack)
			}
		}
	}

	if len(finalCallStacks) == 0 {
		return nil, nil
	}

	// build pseudoID <-> inProfileID relationship
	// pseudoFunctionIdMap := make(map[uint64]uint64, len(interestingFunctionIds)) // pseudoID -> realID
	functionIdPseudoMap := make(map[uint64]uint64, len(interestingFunctionIds)) // realID -> pseudoID
	for i, function := range interestingFunctionIds {
		// pseudoFunctionIdMap[uint64(i)] = function
		functionIdPseudoMap[function] = uint64(i)
	}

	// build path
	path := NewPath(len(interestingFunctionIds))
	for _, callStack := range finalCallStacks {
		for i := 0; i < len(callStack)-1; i++ {
			// convert realID to pseudoID
			to := functionIdPseudoMap[callStack[i]]
			from := functionIdPseudoMap[callStack[i+1]]
			path.Set(int(from), int(to))
		}
	}

	return &Verifier{
		callStacks: finalCallStacks,
		path:       path,
		// pseudoFunctionIdMap: pseudoFunctionIdMap,
		functionIdPseudoMap: functionIdPseudoMap,
		masterProfile:       masterProfile,
	}, nil
}

// Reachable checks if there's a path from function `from` to function `to`.
func (v *Verifier) Reachable(from, to string) bool {
	fromName := v.functionPrefix + from
	toName := v.functionPrefix + to

	fromId, ok := v.masterProfile.functionNameMap[fromName]
	if !ok {
		log.Printf("function %s not found", fromName)
		return false
	}

	toId, ok := v.masterProfile.functionNameMap[toName]
	if !ok {
		log.Printf("function %s not found", toName)
		return false
	}

	return v.path.HasPath(int(v.functionIdPseudoMap[fromId]), int(v.functionIdPseudoMap[toId]))
}

// Next checks if there's a direct path from function `from` to function `to`.
func (v *Verifier) Next(from, to string) bool {
	fromName := v.functionPrefix + from
	toName := v.functionPrefix + to

	fromId, ok := v.masterProfile.functionNameMap[fromName]
	if !ok {
		log.Printf("function %s not found", fromName)
		return false
	}

	toId, ok := v.masterProfile.functionNameMap[toName]
	if !ok {
		log.Printf("function %s not found", toName)
		return false
	}

	return v.path.HasDirectPath(int(v.functionIdPseudoMap[fromId]), int(v.functionIdPseudoMap[toId]))
}

func (v *Verifier) Callstack() [][]uint64 {
	return v.callStacks
}

func (v *Verifier) DumpCallstack() [][]string {
	result := make([][]string, 0, len(v.callStacks))
	for _, callStack := range v.callStacks {
		convertedCallStack := make([]string, 0, len(callStack))
		for _, function := range callStack {
			convertedCallStack = append(convertedCallStack, v.masterProfile.functionIdMap[function])
		}
		result = append(result, convertedCallStack)
	}
	return result
}

func (v *Verifier) SetFunctionPrefix(prefix string) {
	v.functionPrefix = prefix
}

// SubVerifier returns a new Verifier that is a subset of the current Verifier.
//
// If the name pattern is empty, then the new Verifier will be identical to
// the current Verifier.
//
// If the name pattern contradicts with the current Verifier (no match when
// combined), then the new Verifier will be nil.
func (v *Verifier) SubVerifier(namePattern string) (*Verifier, error) {
	// filter call stacks
	var finalCallStacks [][]uint64
	var originalCallStacks [][]uint64 = v.callStacks

	var interestingFunctionIds []uint64 // function IDs that match the name pattern
	if namePattern == "" {
		finalCallStacks = originalCallStacks
		interestingFunctionIds = make([]uint64, 0, len(v.functionIdPseudoMap))
		for f := range v.functionIdPseudoMap {
			interestingFunctionIds = append(interestingFunctionIds, f)
		}
	} else {
		finalCallStacks = make([][]uint64, 0, len(originalCallStacks))
		for fid := range v.functionIdPseudoMap {
			// regex match
			if match, err := regexp.Match(namePattern, []byte(v.masterProfile.functionIdMap[fid])); match {
				// fmt.Printf("Matched: %s\n", name)
				interestingFunctionIds = append(interestingFunctionIds, fid)
			} else if err != nil {
				return nil, err
			}
		}

		for _, callStack := range originalCallStacks {
			reducedCallStack := make([]uint64, 0, len(callStack))
		LOOP_FUNC_IN_CALLSTACK:
			for _, function := range callStack {
				for _, interestingFunction := range interestingFunctionIds {
					if function == interestingFunction {
						// fmt.Printf("Function %d is interesting\n", function)
						reducedCallStack = append(reducedCallStack, function)
						continue LOOP_FUNC_IN_CALLSTACK
					}
				}
			}
			if len(reducedCallStack) > 0 {
				finalCallStacks = append(finalCallStacks, reducedCallStack)
			}
		}
	}

	if len(finalCallStacks) == 0 {
		return nil, nil
	}

	// build pseudoID <-> inProfileID relationship
	// pseudoFunctionIdMap := make(map[uint64]uint64, len(interestingFunctionIds)) // pseudoID -> realID
	functionIdPseudoMap := make(map[uint64]uint64, len(interestingFunctionIds)) // realID -> pseudoID
	for i, function := range interestingFunctionIds {
		// pseudoFunctionIdMap[uint64(i)] = function
		functionIdPseudoMap[function] = uint64(i)
	}

	// build path
	path := NewPath(len(interestingFunctionIds))
	for _, callStack := range finalCallStacks {
		for i := 0; i < len(callStack)-1; i++ {
			// convert realID to pseudoID
			to := functionIdPseudoMap[callStack[i]]
			from := functionIdPseudoMap[callStack[i+1]]
			path.Set(int(from), int(to))
		}
	}

	return &Verifier{
		callStacks: finalCallStacks,
		path:       path,
		// pseudoFunctionIdMap: pseudoFunctionIdMap,
		functionIdPseudoMap: functionIdPseudoMap,
		masterProfile:       v.masterProfile,
	}, nil
}
