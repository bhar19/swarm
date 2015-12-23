package strategy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/scheduler/node"
)

// WeightedNode represents a node in the cluster with a given weight, typically used for sorting
// purposes.
type weightedNode struct {
	Node *node.Node
	// Weight is the inherent value of this node.
	Weight int64
}

type weightedNodeList []*weightedNode

func (n weightedNodeList) Len() int {
	return len(n)
}

func (n weightedNodeList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n weightedNodeList) Less(i, j int) bool {
	var (
		ip = n[i]
		jp = n[j]
	)

	// If the nodes have the same weight sort them out by number of containers.
	if ip.Weight == jp.Weight {
		return len(ip.Node.Containers) < len(jp.Node.Containers)
	}
	return ip.Weight < jp.Weight
}

func weighNodes(config *cluster.ContainerConfig, nodes []*node.Node) (weightedNodeList, error) {
	weightedNodes := weightedNodeList{}

	for _, node := range nodes {
		nodeMemory := node.TotalMemory
		nodeCpus := node.TotalCpus

		//debugging information for io schedule part
                log.WithFields(log.Fields{"Config memory": config.Memory, "config.share": config.CpuShares}).Debugf("Printing Environment values to console")
                log.WithFields(log.Fields{"nodeMemory": nodeMemory, "nodeCpus": nodeCpus}).Debugf("Printing Environment values to console")
                log.WithFields(log.Fields{"UsedMemory": node.UsedMemory, "UsedCpus": node.UsedCpus}).Debugf("Printing Environment values to console")
                log.WithFields(log.Fields{"Config Blkio": config.HostConfig.BlkioWeight, "UsedBlkio": node.UsedBlkio}).Debugf("Printing Environment values to console for blkio used")
                //log.WithFields(log.Fields{"UsedBlkio2": config.BlkioWeight}).Debugf("Printing Environment values to console for blkio used")

		// Skip nodes that are smaller than the requested resources.
		if nodeMemory < int64(config.Memory) || nodeCpus < config.CpuShares {
			continue
		}

		var (
			cpuScore	int64 = 100
			memoryScore	int64 = 100
			blkioScore	int64 = 100
			leafnodeblkio	int64 = 500
		)

		if config.CpuShares > 0 {
			cpuScore = (node.UsedCpus + config.CpuShares) * 100 / nodeCpus
			log.WithFields(log.Fields{"CpuScore": cpuScore}).Debugf("Printing cpu score")
		}
		if config.Memory > 0 {
			memoryScore = (node.UsedMemory + config.Memory) * 100 / nodeMemory
			log.WithFields(log.Fields{"MemoryScore": memoryScore}).Debugf("Printing memory score")
		}

		if config.HostConfig.BlkioWeight > 0 {
			blkioScore = (config.BlkioWeight) * 100 / (node.UsedBlkio + config.BlkioWeight + leafnodeblkio)
			log.WithFields(log.Fields{"blkioScore": blkioScore}).Debugf("Printing blkio score")
		}

		if cpuScore <= 100 && memoryScore <= 100 && blkioScore <= 100 {
			weightedNodes = append(weightedNodes, &weightedNode{Node: node, Weight: cpuScore + memoryScore + blkioScore})
		}
	}

	if len(weightedNodes) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	return weightedNodes, nil
}
