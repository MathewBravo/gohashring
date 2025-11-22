package node

import "time"

type NodeStatus int

const (
	StatusAlive NodeStatus = iota
	StatusSuspect
	StatusDead
	StatusLeft
)

type Node struct {
	ID          string
	Address     string
	Status      NodeStatus
	Incarnation int64
	LastUpdate  time.Time
}

type VirtualNode struct {
	Hash    uint64
	OwnerID string
}
