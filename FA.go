package main

import (
	"sort"
	"strings"
)

type FA struct {
	States             []string `json:"states"`
	Alphabet           string   `json:"alphabet"`
	TransitionFunction []string `json:"transition_function"` // []string{"a,b:c","d,e:f"}
	StartState         string   `json:"start_state"`
	EndStates          []string `json:"end_states"`
}

type basicError struct{
	s string
}

func (b basicError)Error()string{
	return b.s
}

func NewFA(states []string, alphabet []rune, transitionFunction []string, startState string, endStates []string) (FA, error){
	blacklistRunes := []rune("+: ,~")
	if anyInRuneList(blacklistRunes, alphabet){
		return FA{}, basicError{"Cannot use that Alphabet, it contains blacklisted character"}
	}

	if !inList(startState, states){
		return FA{}, basicError{"Start state not in list of States"}
	}

	allEndStatesInStates := true
	for _, e := range endStates{
		if !inList(e, states){
			allEndStatesInStates = false
		}
	}
	if !allEndStatesInStates{
		return FA{}, basicError{"One or more of the end States does not exist in the States list"}
	}

	return FA{
		states,
		string(alphabet),
		transitionFunction,
		startState,
		endStates,
	}, nil
}

func (fa FA) ToDFA() FA{
	// Find the start state and merge it with all States that can be reached with an epsilon transition
	startEpsilonState := combineStates(append([]string{fa.StartState}, findStates(fa.TransitionFunction, fa.StartState+",~")...))
	// Create a new finite automata with the Alphabet and new start state
	fa1 := FA{
		make([]string, 0),
		fa.Alphabet,
		make([]string, 0),
		startEpsilonState,
		make([]string, 0),
	}
	// This is the current state we are scanning and a queue of the States we need to check
	currentState := ""
	statesToDo := []string{startEpsilonState}

	// while there are still items in the queue
	for len(statesToDo) > 0{
		// pop a state from the queue
		currentState, statesToDo = pop(statesToDo)
		// If the state we are checking has already been checked, ignore
		if !inList(currentState, fa1.States) {
			// split the state into its substates (eg a+c -> [a, c])
			subStates := splitStates(currentState)
			// check if the state is final
			if anyInList(fa.EndStates, subStates){
				fa1.EndStates = append(fa1.EndStates, currentState)
			}
			// add this state to the new dfa's list of States
			fa1.States = append(fa1.States, currentState)
			// for every letter in the Alphabet, find and combine all the reachable States, add that new state to the queue, and create a transition to it
			for _, l := range fa1.Alphabet {
				reachableStates := make([]string, 0)
				// for every substate find all the reachable States and add them to our list of all reachable States (for this state and letter)
				for _, s := range subStates {
					fx := s + "," + string(l)
					thisReachableStates := findStates(fa.TransitionFunction, fx)
					reachableStates = append(reachableStates, thisReachableStates...)
				}
				// create a new list of reachable States, but include epsilon reachable States
				reachableStatesIncludingEpsilon := make([]string, 0)
				for _, state := range reachableStates{
					reachableStatesIncludingEpsilon = append(reachableStatesIncludingEpsilon, state)
					reachableStatesIncludingEpsilon = append(reachableStatesIncludingEpsilon, findStates(fa.TransitionFunction, state+",~")...)
				}
				// ensure that there are no duplicates and reorder the States
				reachableStatesIncludingEpsilon = removeDuplicates(reachableStatesIncludingEpsilon)
				// combine these reachable States into one new state ([a, b, c] -> a+b+c)
				newState := combineStates(reachableStatesIncludingEpsilon)
				// create a transition from this state to the new state
				newTransition := currentState + "," + string(l) + ":" + newState
				// add the new state to the checking queue if it is not the empty state
				if newState != "EMPTY" {
					statesToDo = append(statesToDo, newState)
				}
				// add the transition to the new transition function. this will never be a duplicate as it is uniquely identified by this state and letter
				fa1.TransitionFunction = append(fa1.TransitionFunction, newTransition)
			}
		}
	}

	return fa1
}

func (dfa FA) Evaluate(s []rune)(bool,error){
	currentState := dfa.StartState
	for len(s)>0{
		if currentState == "EMPTY"{
			return false, nil
		}
		var c rune
		c, s = popRune(s)
		if !inRuneList(c, []rune(dfa.Alphabet)){
			return false, basicError{"Letter was not in dfa Alphabet so cannot compute"}
		}
		fx := currentState+","+string(c)
		nextStates := findStates(dfa.TransitionFunction, fx)
		if len(nextStates) != 1{
			return false, basicError{"This is not a DFA so cannot calculate"}
		}
		currentState = nextStates[0]
	}
	return inList(currentState, dfa.EndStates), nil
}

func (dfa FA) String() string{
	s := "Finite Automata:\n"
	s += "	States: "+listToString(dfa.States)+"\n"
	s += "	Alphabet: "+runeListToString([]rune(dfa.Alphabet))+"\n"
	s += "	Transition Function: "+ listToString(dfa.TransitionFunction) + "\n"
	s += "	Start State: " + dfa.StartState + "\n"
	s += "	End States: " + listToString(dfa.EndStates)
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