package main

import "github.com/avietrov/genetic-go/genstr"

func main() {
	var config = genstr.Config{"HelloWorld", 7}
	genstr.RunExperiment(&config)

}
