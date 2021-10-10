package main

import (
	"fmt"
	"strings"
)

func main() {
	/*
	demoNFA(nfaEndingIn3Bs(), []string{
		"ababab",
		"bbba",
		"bbbbbbbb",
		"aaaaaa",
		"abababbb",
		"abababaaa",
	})
	demoNFA(nfaWithEpsilon(), []string{
		"00",
		"10",
		"001",
		"101",
		"1",
		"11",
		"01",
		"011",
		"100",
	})*/

	demoNFA(nfaNoBStartNo4BSeq(), []string{
		"ababab",
		"bbba",
		"bbbbbbbb",
		"aaaaaa",
		"abababbb",
		"abababaaa",
		"baaaa",
		"abbbbaaa",
	})
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