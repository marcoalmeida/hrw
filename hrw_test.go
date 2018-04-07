package hrw

import (
	"math"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	dht := New(nil)
	if dht.NodesCount() != 0 {
		t.Error("Expected 0 nodes, got", dht.NodesCount())
	}

	nodes := []Node{struct {
		name   string
		weight float64
	}{name: "a", weight: 1}, struct {
		name   string
		weight float64
	}{name: "b", weight: 1}}

	dht = New(nodes)
	if dht.NodesCount() != 2 {
		t.Error("Expected 0 nodes, got", dht.NodesCount())
	}

}

func TestNodes_AddNode(t *testing.T) {
	dht := New(nil)
	if dht.NodesCount() != 0 {
		t.Error("Expected 0 nodes, got", dht.NodesCount())
	}

	// add a single node
	node := Node{
		name:   "a",
		weight: 1,
	}
	err := dht.AddNode(node)
	if err != nil {
		t.Error(err)
	}
	if dht.NodesCount() != 1 {
		t.Error("Expected 0 nodes, got", dht.NodesCount())
	}

	// fail to add repeated node
	err = dht.AddNode(node)
	if err == nil {
		t.Error("Expected to fail on duplicate AddNode")
	}
}

func TestAll(t *testing.T) {
	var err error
	dht := New(nil)
	nNodes := 5
	nItems := 1000

	// add nodes
	for i := 0; i < nNodes; i++ {
		err = dht.AddNode(Node{
			name:   strconv.Itoa(i),
			weight: 1,
		})
		if err != nil {
			t.Error(err)
		}
	}

	// count the number of times a given node is used
	counter := make(map[string]int64)
	for i := 0; i < nItems; i++ {
		counter[dht.GetNode(strconv.Itoa(i))]++
	}

	// make sure the distribution is reasonably uniform (off by 2.5% at most)
	for node, count := range counter {
		ideal := 1.0 / float64(nNodes)
		pct := float64(count) / float64(nItems)
		if math.Abs(pct-ideal) > 0.025 {
			t.Errorf("expected %.2f, got %.2f on node %s\n", ideal, pct, node)
		}
	}

	// remove one node and repeat the queries
	dht.RemoveNode("0")

	// count the number of times a given node is used
	counter = make(map[string]int64)
	for i := 0; i < nItems; i++ {
		counter[dht.GetNode(strconv.Itoa(i))]++
	}

	// make sure the distribution is reasonably uniform (off by 3% at most, given that we just removed a node)
	for node, count := range counter {
		// nNodes - 1 because we just removed one node
		ideal := 1.0 / float64(nNodes-1)
		pct := float64(count) / float64(nItems)
		if math.Abs(pct-ideal) > 0.03 {
			t.Errorf("expected %.2f, got %.2f on node %s\n", ideal, pct, node)
		}
	}
}
