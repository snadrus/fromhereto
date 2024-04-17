package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

type Package struct {
	Imports []string
}

func main() {
	tree := explore(os.Args[1])
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(">", ",,")
	enc.Encode(tree)
}
func explore(start string) map[string]any {
	cmd := exec.Command("go list -json " + start)
	b, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var pkg Package
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&pkg)
	if err != nil {
		panic(err)
	}
	res := make(map[string]any)
	for _, imp := range pkg.Imports {
		res[imp] = explore(imp)
	}
	return res
}
