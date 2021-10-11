package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	//generate_fa_args := []string{"new", "testfa.json"}
	//convert_fa_args := []string{"conv", "testfa.json", "dfa.json"}
	test_fa_args := []string{"test", "dfa.json", "abababbb"}
	args = test_fa_args

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
	return NewFA(
		[]string{"0"},
		[]rune{},
		[]string{},
		"0",
		[]string{"0"},
	)
}

func testFiniteAutomata(filename string, s string){
	fa := readFAFromFile(filename)
	if fa.Evaluate([]rune(s)){
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




func demoNFA(nfa FA, testInputs []string){
	dfa := nfa.ToDFA()
	fmt.Println("------------------Original NFA------------------")
	fmt.Println(nfa)
	fmt.Println("------------------Generated DFA-----------------")
	fmt.Println(dfa)
	fmt.Println("------------------------------------------------")
	acceptions, rejections := evaluateAll(dfa, testInputs)
	fmt.Print("The DFA Accepts: ")
	fmt.Println(acceptions)
	fmt.Print("The DFA Rejects: ")
	fmt.Println(rejections)
}



func evaluateAll(dfa FA, inputs []string)([]string,[]string){
	acceptions := make([]string, 0)
	rejections := make([]string, 0)

	for _, inp := range inputs{
		if dfa.Evaluate([]rune(inp)){
			acceptions = append(acceptions, inp)
		} else{
			rejections = append(rejections, inp)
		}
	}

	return acceptions, rejections
}

// This accepts 00,10,001,101,1 and nothing else
func nfaWithEpsilon() FA{
	return NewFA(
		strings.Split("a,b,c,d", ","),
		[]rune{'0', '1'},
		[]string{
			"a,0:b",
			"a,1:b",
			"a,~:d",
			"b,0:c",
			"b,0:d",
			"d,1:c",
		},
		"a",
		[]string{"c"},
	)
}
func nfaEndingIn3Bs() FA{
	return NewFA(
		strings.Split("0,1,2,3", ","),
		[]rune{'a', 'b'},
		[]string{
			"0,a:0",
			"0,b:0",
			"0,b:1",
			"1,b:2",
			"2,b:3",
		},
		"0",
		[]string{"3"},
	)
}

// this cannot start with a b or have more than 3 bs in a row
func nfaNoBStartNo4BSeq()FA{
	return NewFA(
		strings.Split("1,2,3,4", ","),
		[]rune{'a', 'b'},
		[]string{
			"1,a:1",
			"1,a:2",
			"1,a:3",
			"1,a:4",
			"2,b:1",
			"3,b:2",
			"4,b:3",
		},
		"1",
		[]string{"1"},
		)
}

