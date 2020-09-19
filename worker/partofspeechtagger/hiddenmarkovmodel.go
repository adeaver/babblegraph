package partofspeechtagger

import (
	"fmt"

	"github.com/adeaver/babblegraph/lib/math/decimal"
)

const startToken string = "<START>"

type HiddenMarkovModel interface {
	GetAllSymbols() []string
	GetSymbolCountsForToken(token string) (map[string]int64, error)
	GetTrigramCountsForSymbol(symbol string) (map[Trigram]int64, error)
}

type Trigram struct {
	FirstSymbol  string
	SecondSymbol string
	ThirdSymbol  string
}

func (t Trigram) makeBigramKey() string {
	makeBigramKey(t.FirstSymbol, t.SecondSymbol)
}

func makeBigramKey(firstSymbol, secondSymbol string) {
	return fmt.Sprintf("%s:%s", firstSymbol, secondSymbol)
}

type taggedNode struct {
	BackReferenceIdx int
	Probability      decimal.Number
}

func GetPartOfSpeechSequence(model HiddenMarkovModel, tokenSequence []string) ([]string, *decimal.Number, error) {
	allSymbols := model.GetAllSymbols()
	nodes := make([][]taggedNode, len(tokenSequence))
	for tokenIdx, token := range tokenSequence {
		symbolCountsForToken, err := model.GetSymbolCountsForToken(token)
		if err != nil {
			return nil, nil, err
		}
		symbolProbabilitiesForToken := makeSymbolProbabilitiesForToken(symbolCountsForToken)
		var nodesForToken []taggedNode
		switch tokenIdx {
		case 0:
			nodesForToken, err = getNodesForToken(nodesForTokenInput{
				Model:                      model,
				AllSymbols:                 allSymbols,
				SymbolProbabilitesForToken: symbolProbabilitiesForToken,
			})
			if err != nil {
				return nil, nil, err
			}
		case 1:
			nodesForToken, err = getNodesForToken(nodesForTokenInput{
				Model:                      model,
				AllSymbols:                 allSymbols,
				SymbolProbabilitesForToken: symbolProbabilitiesForToken,
				PreviousNodeList:           nodes[idx-1],
			})
			if err != nil {
				return nil, nil, err
			}
		default:
			nodesForToken, err = getNodesForToken(nodesForTokenInput{
				Model:                      model,
				AllSymbols:                 allSymbols,
				SymbolProbabilitesForToken: symbolProbabilitiesForToken,
				PreviousNodeList:           nodes[idx-1],
				ShouldBackRef:              true,
			})
			if err != nil {
				return nil, nil, err
			}
		}
		nodes = append(nodes, nodesForToken)
	}
	out, probability := getSymbolsFromBackref(allSymbols, nodes)
	return out, &probability, nil
}

func getSymbolsFromBackref(allSymbols []string, nodes [][]taggedNode) ([]string, decimal.Number) {
	out := make([]string, len(nodes))
	mostLikelySequenceEndNode, symbolIdx = getHighestProbabilityNode(nodes[len(nodes)-1])
	out[len(nodes)-1] = allSymbols[symbolIdx]
	symbolIdx := mostLikelySequenceEndNode.BackReferenceIdx
	for i := len(nodes) - 2; i >= 0; i-- {
		currentNode := nodes[i][symbolIdx]
		out[i] = allSymbols[symbolIdx]
		symbolIdx := currentNode.BackReferenceIdx
	}
	return out, mostLikelySequenceEndNode.Probability
}

func getHighestProbabilityNode(nodes []taggedNode) (taggedNode, int) {
	var highestProbabilityIdx *int
	var probabilityOfMostLikelyIdx *decimal.Number
	for nodeIdx, node := range nodes {
		if probabilityOfMostLikelyIdx == nil || node.Probability.GreaterThan(*probabilityOfMostLikelyIdx) {
			highestProbabilityIdx = &nodeIdx
			probabilityOfMostLikelyIdx = &node.Probability
		}
	}
	return nodes[*highestProbabilityIdx], *highestProbabilityIdx
}

type nodesForTokenInput struct {
	Model                      HiddenMarkovModel
	AllSymbols                 []string
	SymbolProbabilitesForToken map[string]decimal.Number
	PreviousNodeList           *[]taggedNode
	ShouldBackRef              bool
}

func getNodesForToken(input nodesForTokenInput) ([]taggedNodes, error) {
	var out []taggedNode

	for symbolIdx, symbol := range input.AllSymbols {
		probabilityOfSymbol := input.SymbolProbabilitesForToken[symbol]
		trigramCountsForSymbol, err := input.Model.GetTrigramCountsForSymbol(symbol)
		if err != nil {
			return nil, err
		}
		trigramProbabilities := makeTrigramProbabilitiesForSymbol(trigramCountsForSymbol)
		var taggedNodeForSymbol taggedNode
		switch {
		case input.PreviousNodeList == nil:
			// There are no previous nodes
			bigramKey := makeBigramKey(startToken, startToken)
			trigramProbabilityForSymbol := trigramProbabilities[bigramKey]
			probabilityOfNode := probabilityOfSymbol.Multiply(trigramProbabilityForSymbol)
			taggedNodeForSymbol = taggedNode{
				Probability: probabilityOfNode,
			}

		case input.PreviousNodeList != nil:
			var highestProbabilityIdx *int
			var probabilityOfMostLikelyIdx *decimal.Number
			for nodeIdx, node := range *input.PreviousNodeList {
				backRefSymbol := startToken
				if input.ShouldBackRef {
					backRefSymbol = input.AllSymbols[node.BackReferenceIdx]
				}
				bigramKeyForNode := makeBigramKey(backRefSymbol, input.AllSymbols[nodeIdx])
				trigramProbabilityOfNode := trigramProbabilities[bigramKeyForNode]
				probabilityOfNode := trigramProbabilityOfNode.Multiply(node.Probability).Multiply(probabilityOfSymbol)
				if probabilityOfMostLikelyIdx == nil || probabilityOfNode.GreaterThan(*probabilityOfMostLikelyIdx) {
					probabilityOfMostLikelyIdx = &probabilityOfNode
					highestProbabilityIdx = highestProbabilityIdx
				}
			}
			taggedNodeForSymbol = taggedNode{
				Probability:      *probabilityOfMostLikelyIdx,
				BackReferenceIdx: *highestProbabilityIdx,
			}
		}
		out = append(out, taggedNodeForSymbol)
	}
	return out
}

func makeSymbolProbabilitiesForToken(symbolCountsForToken map[string]int64) map[string]decimal.Number {
	var totalCount decimal.Number
	for _, count := range symbolCountsForToken {
		totalCount.Add(decimal.FromInt64(count))
	}
	out := make(map[string]decimal.Number)
	for symbol, count := range symbolCountsForToken {
		out[symbol] = decimal.FromInt64(count).Divide(totalCount)
	}
	return out
}

func makeTrigramProbabilitiesForSymbol(trigramCountsForSymbol map[Trigram]int64) map[string]decimal.Number {
	var totalCount decimal.Number
	for _, count := range trigramCountsForSymbol {
		totalCount.Add(decimal.FromInt64(count))
	}
	out := make(map[string]decimal.Number)
	for trigram, count := range trigramCountsForSymbol {
		out[trigram.makeBigramKey()] = decimal.FromInt64(count).Divide(totalCount)
	}
	return out
}
