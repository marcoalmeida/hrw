package hrw

import (
	"errors"
	"hash/fnv"
	"math"
	"sort"
)

type Node struct {
	Name   string
	Weight float64
}

type nodeInfo struct {
	hash   uint64
	weight float64
}

// mapping on nodes: Name -> Weight
type Nodes map[string]nodeInfo

func New(nodesList []Node) Nodes {
	// we need fast lookups -> transform the input list into a dictionary
	nodes := make(map[string]nodeInfo)
	for _, n := range nodesList {
		nodes[n.Name] = nodeInfo{
			weight: n.Weight,
			hash:   hash64(n.Name),
		}
	}

	return nodes
}

func (hrw Nodes) NodesCount() int {
	return len(hrw)
}

func (hrw Nodes) GetAllNodes() []string {
	nodes := make([]string, hrw.NodesCount())

	i := 0
	for name := range hrw {
		nodes[i] = name
		i++
	}

	return nodes
}

func (hrw Nodes) AddNode(node Node) error {
	_, ok := hrw[node.Name]

	if ok {
		return errors.New("node already exists")
	}

	hrw[node.Name] = nodeInfo{
		weight: node.Weight,
		hash:   hash64(node.Name),
	}

	return nil
}

func (hrw Nodes) RemoveNode(node string) error {
	_, ok := hrw[node]

	if !ok {
		return errors.New("node not found")
	}

	delete(hrw, node)

	return nil
}

type pair struct {
	k string
	v float64
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Less(i, j int) bool { return p[i].v < p[j].v }

func (hrw Nodes) GetNodesRanked(key string) []string {
	rank := make(pairList, hrw.NodesCount())
	nodes := make([]string, hrw.NodesCount())

	// compute the score for all nodes
	i := 0
	for name, info := range hrw {
		score := weightedScore(key, info.hash, info.weight)
		rank[i] = pair{k: name, v: score}
		i++
	}

	// sort the rank list (descending)
	sort.Sort(sort.Reverse(rank))

	// create and return all nodes (ranked)
	for j, p := range rank {
		nodes[j] = p.k
	}

	return nodes
}

func (hrw Nodes) GetNode(key string) string {
	return hrw.GetNodesRanked(key)[0]
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
