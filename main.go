package main

import (
	. "courseWork/assembly"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("testForMe.asm")
	if err != nil {
		panic(err)
	}

	assembly := CreateAssembly(string(content))
	err = FirstStage(assembly)
	if err != nil {
		panic(err)
	}
	PrintFirstStage(assembly)

	FirstPass(assembly)
	PrintFirstPass(assembly)
}