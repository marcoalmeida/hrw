package hrw

import (
	"errors"
	"hash/fnv"
	"math"
	"math/rand"
	"strconv"
)

type Node struct {
	name   string
	weight float64
}

type nodeInfo struct {
	// use an array of bytes to avoid conversions on all calls to the hash function
	seed   []byte
	weight float64
}

type nodes map[string]nodeInfo

func New(nodesList []Node) nodes {
	// we need fast lookups -> transform the input list into a dictionary
	nodes := make(map[string]nodeInfo)
	for _, n := range nodesList {
		nodes[n.name] = nodeInfo{
			weight: n.weight,
			seed:   []byte(strconv.Itoa(int(rand.Int63()))),
		}
	}

	return nodes
}

func (hrw nodes) NodesCount() int {
	return len(hrw)
}

func (hrw nodes) AddNode(node Node) error {
	_, ok := hrw[node.name]

	if ok {
		return errors.New("node already exists")
	}

	hrw[node.name] = nodeInfo{
		weight: node.weight,
		seed:   []byte(strconv.Itoa(int(rand.Int63()))),
	}

	return nil
}

func (hrw nodes) RemoveNode(node string) error {
	_, ok := hrw[node]

	if !ok {
		return errors.New("node not found")
	}

	delete(hrw, node)

	return nil
}

func (hrw nodes) GetNode(key string) string {
	highestScore := -1.0
	champion := ""

	for node, info := range hrw {
		score := weightedScore(key, info)
		if score > highestScore {
			champion = node
			highestScore = score
		}
	}

	return champion
}

func weightedScore(key string, node nodeInfo) float64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	// this hash function does not support seeding, but we can
	// obtain a similar result by concatenating the seed
	h.Write(node.seed)

	score := 1.0 / -math.Log(int2float(h.Sum64()))

	return node.weight * score
}

// converts a uniformly random 64-bit integer to uniformly random floating point number on interval [0, 1)
func int2float(v uint64) float64 {
	fiftyThreeOnes := uint64(0xFFFFFFFFFFFFFFFF >> (64 - 53))
	fiftyThreeZeros := float64(1 << 53)
	return float64(v&fiftyThreeOnes) / fiftyThreeZeros
}
