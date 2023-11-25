package pprofsv_test

// func TestProfile(t *testing.T) {
// 	file, err := os.Open("testdata/pprof.profile")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	profile, err := profile.Parse(file)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// for _, sample := range profile.Sample {
// 	// 	fmt.Printf("Sample: \n")
// 	// 	for _, location := range sample.Location {
// 	// 		var functionChain []string
// 	// 		for _, line := range location.Line {
// 	// 			functionChain = append(functionChain, fmt.Sprintf("%d:%s", line.Function.ID, line.Function.Name))
// 	// 		}
// 	// 		fmt.Printf("Location#%d: %s\n", location.ID, functionChain)
// 	// 	}
// 	// }

// 	for _, function := range profile.Function {
// 		t.Logf("%d:%s", function.ID, function.Name)
// 	}
// }
