package member

import (
	"time"

	"github.com/MathewBravo/gohashring/internal/node"
)

type MemberRecord struct {
	Status      node.NodeStatus
	Incarnation int64
	LastUpdate  time.Time
	Address     string
}

type MembershipState struct {
	SelfID          string
	SelfIncarnation int64
	Nodes           map[string]MemberRecord
	AliveNodeIDs    []string
	Version         int64
}
