package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	//generate_fa_args := []string{"new", "testfa.json"}
	//convert_fa_args := []string{"conv", "testfa.json", "dfa.json"}
	//test_fa_args := []string{"test", "dfa.json", "abababbb"}
	//args = test_fa_args

	cmd, args := popArg(args)

	switch cmd{
	case "new":
		//new empty
		var filename string
		filename, args = popArg(args)
		createEmptyFAFile(filename)
		break
	case "test":
		//testing
		var filename string
		var testString string
		filename, args = popArg(args)
		testString, args = popArg(args)
		testFiniteAutomata(filename, testString)
		break
	case "conv":
		// convert to dfa
		var filename1, filename2 string
		filename1, args = popArg(args)
		filename2, args = popArg(args)
		convertFiniteAutomata(filename1, filename2)
		break
	}

}

func popArg(xs []string) (string, []string){
	if len(xs) == 0{
		consoleError("Incorrect command usage, please look at github for more info")
	}
	return xs[0], xs[1:]
}


func consoleError(reason string){
	fmt.Println("Error: "+reason)
	os.Exit(0)
}

func createEmptyFAFile(fileName string){
	js, _ := json.MarshalIndent(nfaEmpty(), "", "\t")
	err := os.WriteFile(fileName, js, 0644)
	if err != nil{
		consoleError(err.Error())
	}
}
func nfaEmpty() FA{
	fa, _ := NewFA(
		[]string{"0"},
		[]rune{},
		[]string{},
		"0",
		[]string{"0"},
	)
	return fa
}

func testFiniteAutomata(filename string, s string){
	fa := readFAFromFile(filename)
	eval, err := fa.Evaluate([]rune(s))
	if err != nil{
		consoleError(err.Error())
	}
	if eval{
		fmt.Println("accept")
	} else{
		fmt.Println("reject")
	}
}

func convertFiniteAutomata(fn1, fn2 string){
	fa := readFAFromFile(fn1)
	dfa := fa.ToDFA()
	fmt.Println(fa)
	js, _ := json.MarshalIndent(dfa, "", "\t")
	err := os.WriteFile(fn2, js, 0644)
	if err != nil{
		consoleError(err.Error())
	}

}

func readFAFromFile(fn string) FA{
	js, err := os.ReadFile(fn)
	if err != nil{
		consoleError(err.Error())
	}
	fa := FA{}
	err = json.Unmarshal(js, &fa)
	if err != nil{
		consoleError(err.Error())
	}
	return fa
}