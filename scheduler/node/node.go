package node

import (
	"errors"

	"github.com/docker/swarm/cluster"
)

// Node is an abstract type used by the scheduler.
type Node struct {
	ID         string
	IP         string
	Addr       string
	Name       string
	Labels     map[string]string
	Containers cluster.Containers
	Images     []*cluster.Image

	UsedMemory  int64
	UsedCpus    int64
	UsedBlkiio  int64
	TotalMemory int64
	TotalCpus   int64

	IsHealthy bool
}

// NewNode creates a node from an engine.
func NewNode(e *cluster.Engine) *Node {
	return &Node{
		ID:          e.ID,
		IP:          e.IP,
		Addr:        e.Addr,
		Name:        e.Name,
		Labels:      e.Labels,
		Containers:  e.Containers(),
		Images:      e.Images(),
		UsedMemory:  e.UsedMemory(),
		UsedCpus:    e.UsedCpus(),
		UsedBlkio:   e.UsedBlkio(),
		TotalMemory: e.TotalMemory(),
		TotalCpus:   e.TotalCpus(),
		IsHealthy:   e.IsHealthy(),
	}
}

// Container returns the container with IDOrName in the engine.
func (n *Node) Container(IDOrName string) *cluster.Container {
	return n.Containers.Get(IDOrName)
}

// AddContainer injects a container into the internal state.
func (n *Node) AddContainer(container *cluster.Container) error {
	if container.Config != nil {
		memory := container.Config.Memory
		cpus := container.Config.CpuShares
		blkio := container.Config.HostConfig.BlkioWeight
		if n.TotalMemory-memory < 0 || n.TotalCpus-cpus < 0 {
			return errors.New("not enough resources")
		}
		n.UsedMemory = n.UsedMemory + memory
		n.UsedCpus = n.UsedCpus + cpus
		n.UsedBlkio = n.UsedBlkio + blkio
	}
	n.Containers = append(n.Containers, container)
	return nil
}
