package main

import (
	"sort"
	"strings"
)

type FA struct{
	states []string
	alphabet []rune
	transitionFunction []string // []string{"a,b:c","d,e:f"}
	startState string
	endStates []string
}

func NewFA(states []string, alphabet []rune, transitionFunction []string, startState string, endStates []string) FA{
	blacklistRunes := []rune("+: ,~")
	if anyInRuneList(blacklistRunes, alphabet){
		panic("Cannot use that alphabet, it contains blacklisted character")
	}

	if !inList(startState, states){
		panic("Start state not in list of states")
	}

	allEndStatesInStates := true
	for _, e := range endStates{
		if !inList(e, states){
			allEndStatesInStates = false
		}
	}
	if !allEndStatesInStates{
		panic("One or more of the end states does not exist in the states list")
	}

	return FA{
		states,
		alphabet,
		transitionFunction,
		startState,
		endStates,
	}
}

func (fa FA) ToDFA() FA{
	// Find the start state and merge it with all states that can be reached with an epsilon transition
	startEpsilonState := combineStates(append([]string{fa.startState}, findStates(fa.transitionFunction, fa.startState+",~")...))
	// Create a new finite automata with the alphabet and new start state
	fa1 := FA{
		make([]string, 0),
		fa.alphabet,
		make([]string, 0),
		startEpsilonState,
		make([]string, 0),
	}
	// This is the current state we are scanning and a queue of the states we need to check
	currentState := ""
	statesToDo := []string{startEpsilonState}

	// while there are still items in the queue
	for len(statesToDo) > 0{
		// pop a state from the queue
		currentState, statesToDo = pop(statesToDo)
		// If the state we are checking has already been checked, ignore
		if !inList(currentState, fa1.states) {
			// split the state into its substates (eg a+c -> [a, c])
			subStates := splitStates(currentState)
			// check if the state is final
			if anyInList(fa.endStates, subStates){
				fa1.endStates = append(fa1.endStates, currentState)
			}
			// add this state to the new dfa's list of states
			fa1.states = append(fa1.states, currentState)
			// for every letter in the alphabet, find and combine all the reachable states, add that new state to the queue, and create a transition to it
			for _, l := range fa1.alphabet {
				reachableStates := make([]string, 0)
				// for every substate find all the reachable states and add them to our list of all reachable states (for this state and letter)
				for _, s := range subStates {
					fx := s + "," + string(l)
					thisReachableStates := findStates(fa.transitionFunction, fx)
					reachableStates = append(reachableStates, thisReachableStates...)
				}
				// create a new list of reachable states, but include epsilon reachable states
				reachableStatesIncludingEpsilon := make([]string, 0)
				for _, state := range reachableStates{
					reachableStatesIncludingEpsilon = append(reachableStatesIncludingEpsilon, state)
					reachableStatesIncludingEpsilon = append(reachableStatesIncludingEpsilon, findStates(fa.transitionFunction, state+",~")...)
				}
				// ensure that there are no duplicates and reorder the states
				reachableStatesIncludingEpsilon = removeDuplicates(reachableStatesIncludingEpsilon)
				// combine these reachable states into one new state ([a, b, c] -> a+b+c)
				newState := combineStates(reachableStatesIncludingEpsilon)
				// create a transition from this state to the new state
				newTransition := currentState + "," + string(l) + ":" + newState
				// add the new state to the checking queue if it is not the empty state
				if newState != "EMPTY" {
					statesToDo = append(statesToDo, newState)
				}
				// add the transition to the new transition function. this will never be a duplicate as it is uniquely identified by this state and letter
				fa1.transitionFunction = append(fa1.transitionFunction, newTransition)
			}
		}
	}

	return fa1
}

func (dfa FA) Evaluate(s []rune)bool{
	currentState := dfa.startState
	for len(s)>0{
		if currentState == "EMPTY"{
			return false
		}
		var c rune
		c, s = popRune(s)
		if !inRuneList(c, dfa.alphabet){
			panic("Letter was not in dfa alphabet so cannot compute")
		}
		fx := currentState+","+string(c)
		nextStates := findStates(dfa.transitionFunction, fx)
		if len(nextStates) != 1{
			panic("This is not a DFA so cannot calculate")
		}
		currentState = nextStates[0]
	}
	return inList(currentState, dfa.endStates)
}

func (dfa FA) String() string{
	s := "Finite Automata:\n"
	s += "	States: "+listToString(dfa.states)+"\n"
	s += "	Alphabet: "+runeListToString(dfa.alphabet)+"\n"
	s += "	Transition Function: "+ listToString(dfa.transitionFunction) + "\n"
	s += "	Start State: " + dfa.startState + "\n"
	s += "	End States: " + listToString(dfa.endStates)
	return s
}

// ----------------------------------------Util functions----------------------------------------

func combineStates(ss []string)string{
	if len(ss) == 0{
		return "EMPTY"
	}
	s := ""
	for _, s1 := range ss{
		s1 += "+"
		s += s1
	}
	return s[:len(s)-1]
}

func splitStates(s string)[]string{
	return strings.Split(s, "+")
}

func removeDuplicates(xs []string) []string{
	xs1 := make([]string, 0)
	for _, x := range xs{
		if !inList(x, xs1){
			xs1 = append(xs1, x)
		}
	}
	// This is nescesary because order of the input set xs should not affect the state name (eg [a,c,b] should be a+b+c not a+c+b)
	sort.Strings(xs1)
	return xs1
}

func inList(x string, xs []string) bool{
	for _, x1 := range xs{
		if x == x1{
			return true
		}
	}
	return false
}
func inRuneList(x rune, xs []rune) bool{
	for _, x1 := range xs{
		if x == x1{
			return true
		}
	}
	return false
}


func anyInList(xs1 []string, xs2 []string) bool{
	for _, x1 := range xs1{
		for _, x2 := range xs2{
			if x1 == x2{
				return true
			}
		}
	}
	return false
}

func anyInRuneList(xs1 []rune, xs2 []rune) bool{
	for _, x1 := range xs1{
		for _, x2 := range xs2{
			if x1 == x2{
				return true
			}
		}
	}
	return false
}

func pop(xs []string) (string, []string){
	return xs[0], xs[1:]
}

func popRune(xs []rune) (rune, []rune){
	return xs[0], xs[1:]
}

func findStates(transitionFunction []string, fx string) []string{
	states := make([]string, 0)
	for _, y := range transitionFunction{
		parts := strings.Split(y, ":")
		fx1, s1 := parts[0], parts[1]
		if fx1 == fx{
			states = append(states, s1)
		}
	}
	return states
}



func listToString(xs []string)string{
	if len(xs) == 0{
		return "{}"
	}
	x := "{"
	for _, x1 := range xs{
		x += "\""+string(x1) + "\", "
	}
	x = x[:len(x)-2] + "}"
	return x
}


func runeListToString(xs []rune)string{
	if len(xs) == 0{
		return "{}"
	}
	x := "{"
	for _, x1 := range xs{
		x += "\""+string(x1) + "\", "
	}
	x = x[:len(x)-2] + "}"
	return x
}