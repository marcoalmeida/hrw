package hrw

import (
	"errors"
	"hash/fnv"
	"math"
)

type Node struct {
	name   string
	weight float64
}

type nodeInfo struct {
	hash   uint64
	weight float64
}

// mapping on nodes: name -> weight
type nodes map[string]nodeInfo

func New(nodesList []Node) nodes {
	// we need fast lookups -> transform the input list into a dictionary
	nodes := make(map[string]nodeInfo)
	for _, n := range nodesList {
		nodes[n.name] = nodeInfo{
			weight: n.weight,
			hash:   hash64(n.name),
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
		hash:   hash64(node.name),
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

	for name, info := range hrw {
		score := weightedScore(key, info.hash, info.weight)
		if score > highestScore {
			champion = name
			highestScore = score
		}
	}

	return champion
}

func weightedScore(key string, nodeHash uint64, nodeWeight float64) float64 {
	hash := mergeHashes(hash64(key), nodeHash)
	score := 1.0 / -math.Log(int2float(hash))

	return nodeWeight * score
}

func hash64(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

func mergeHashes(x uint64, y uint64) uint64 {
	// 18446744073709551616 == 2**64-1
	const twoToSixtyFour = 18446744073709551615

	acc := x ^ y
	acc ^= acc >> 33
	acc = (acc * 0xff51afd7ed558ccd) % twoToSixtyFour
	acc ^= acc >> 33
	acc = (acc * 0xc4ceb9fe1a85ec53) % twoToSixtyFour
	acc ^= acc >> 33

	return acc
}

// converts a uniformly random 64-bit integer to uniformly random floating point number on interval [0, 1)
func int2float(v uint64) float64 {
	fiftyThreeOnes := uint64(0xFFFFFFFFFFFFFFFF >> (64 - 53))
	fiftyThreeZeros := float64(1 << 53)
	return float64(v&fiftyThreeOnes) / fiftyThreeZeros
}
