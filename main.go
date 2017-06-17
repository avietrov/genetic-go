package main

import (
	"fmt"

	"github.com/avietrov/genetic-go/genstr"
)

func main() {
	var config = genstr.Config{"Goisanopensourceprogramminglanguagethatmakesiteasytobuildsimplereliableandefficientsoftware", 30, 5, 1, 5}
	var result = genstr.RunExperiment(&config)
	fmt.Println("Done: ", result)

}
