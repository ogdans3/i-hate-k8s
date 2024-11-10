package models

type NodeRole uint8

const (
	ControlPlane NodeRole = iota
	Worker
)

type Node struct {
	Ip       string
	Name     string
	HostName string
	Role     NodeRole
}
