package ring

import (
	"fmt"
	"sort"

	"github.com/MathewBravo/gohashring/internal/hash"
	"github.com/MathewBravo/gohashring/internal/node"
)

type RingView struct {
	VirtualNodes       []node.VirtualNode
	NodeToVNodeIndexes map[string][]int
	ReplicationFactor  int
	VNodePerNode       int
	Version            int64
}

func BuildRing(nodes []string, vnodePerNode int, rf int) RingView {
	sort.Strings(nodes)
	var vNodes []node.VirtualNode

	for _, n := range nodes {
		for i := range vnodePerNode {
			label := fmt.Sprintf("%s#%d", n, i)
			hl := hash.Hash64([]byte(label))

			v := node.VirtualNode{
				Hash:    hl,
				OwnerID: n,
			}

			vNodes = append(vNodes, v)
		}
	}

	sort.Slice(vNodes, func(i, j int) bool {
		return vNodes[i].Hash < vNodes[j].Hash
	})

	nToVnIndex := make(map[string][]int)
	for i, vn := range vNodes {
		nToVnIndex[vn.OwnerID] = append(nToVnIndex[vn.OwnerID], i)
	}

	return RingView{
		VirtualNodes:       vNodes,
		NodeToVNodeIndexes: nToVnIndex,
		ReplicationFactor:  rf,
		VNodePerNode:       vnodePerNode,
		Version:            1,
	}
}

func (r *RingView) LookUpPrimaryOwner(key []byte) string {
	h := hash.Hash64(key)

	i := sort.Search(len(r.VirtualNodes), func(i int) bool {
		return r.VirtualNodes[i].Hash >= h
	})

	i %= len(r.VirtualNodes)

	return r.VirtualNodes[i].OwnerID
}

// TODO:
// Edge cases: empty ring, n bigger than number of nodes, n <= 0
func (r *RingView) LookUpNOwners(key []byte, n int) []string {
	h := hash.Hash64(key)

	vn := sort.Search(len(r.VirtualNodes), func(i int) bool {
		return r.VirtualNodes[i].Hash >= h
	})

	vn %= len(r.VirtualNodes)

	seen := make(map[string]bool)
	owners := make([]string, 0, n)

	for i := 0; len(owners) < n && i < len(r.VirtualNodes); i++ {
		idx := (vn + i) % len(r.VirtualNodes)
		ownerID := r.VirtualNodes[idx].OwnerID

		if !seen[ownerID] {
			seen[ownerID] = true
			owners = append(owners, ownerID)
		}
	}

	return owners
}
