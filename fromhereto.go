package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Package struct {
	Imports []string
}

var top = map[string][]string{}
var topMx sync.Mutex
var wg sync.WaitGroup
var parallel = make(chan struct{}, 20)

var deep = map[string]map[string]bool{}

type Weights struct {
	Weight     int
	DeepWeight int
}

func main() {
	wg.Add(1)
	go explore(os.Args[1])
	wg.Wait()

	// follow flows
	for k := range top {
		fillDeep(k) // memoized
	}

	// This returns how many unique packages are imported by a child import, but not
	// how many of a given package's unique imports are imported by a child import,
	// because what would you do with duplicates?
	branchWeight := map[string]map[string]int{}
	for k, v := range top {
		branchWeight[k] = map[string]int{}
		for _, imp := range v {
			branchWeight[k][imp] = len(deep[imp])
		}
		branchWeight[k]["_unique"] = len(deep[k])
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(branchWeight)
}

func fillDeep(pkg string) map[string]bool {
	if _, ok := deep[pkg]; ok {
		return deep[pkg]
	}
	deep[pkg] = map[string]bool{}
	for _, imp := range top[pkg] {
		deep[pkg][imp] = true
		for k, v := range fillDeep(imp) {
			if strings.Contains(k, ".") {
				deep[pkg][k] = v
			}
		}
	}
	return deep[pkg]
}

func explore(start string) {
	defer wg.Done()
	parallel <- struct{}{}
	defer func() {
		<-parallel
	}()
	cmd := exec.Command("go", "list", "-json", start)
	b, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	var pkg Package
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&pkg)
	if err != nil {
		panic(err)
	}
	topMx.Lock()
	defer topMx.Unlock()
	top[start] = pkg.Imports

	for _, imp := range pkg.Imports {
		if strings.Contains(imp, ".") {
			if _, ok := top[imp]; !ok {
				wg.Add(1)
				top[imp] = nil
				go explore(imp)
			}
		}
	}
}
