package hrw

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

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
	nNodes := 7
	nItems := 10000
	itemLength := 11

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
		counter[dht.GetNode(randomString(itemLength))]++
	}
	ensureUniformDistribution(counter, nItems, dht.NodesCount(), t)

	// remove one node and repeat the queries
	dht.RemoveNode("0")
	counter = make(map[string]int64)
	for i := 0; i < nItems; i++ {
		counter[dht.GetNode(randomString(itemLength))]++
	}
	ensureUniformDistribution(counter, nItems, dht.NodesCount(), t)

	// add 2 nodes and do it again
	dht.AddNode(Node{name: "new0", weight: 1})
	dht.AddNode(Node{name: "new1", weight: 1})
	counter = make(map[string]int64)
	for i := 0; i < nItems; i++ {
		counter[dht.GetNode(randomString(itemLength))]++
	}
	ensureUniformDistribution(counter, nItems, dht.NodesCount(), t)
}

// make sure the distribution is reasonably uniform -- 3%
func ensureUniformDistribution(counter map[string]int64, nItems int, nNodes int, t *testing.T) {
	maxDifference := 0.03

	for node, count := range counter {
		ideal := 1.0 / float64(nNodes)
		pct := float64(count) / float64(nItems)
		if math.Abs(pct-ideal) > maxDifference {
			t.Errorf("expected %.2f, got %.2f on node %s\n", ideal, pct, node)
		}
	}
}
