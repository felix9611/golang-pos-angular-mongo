package tools

import (
	"sort"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TreeNode struct {
	ID        primitive.ObjectID `json:"_id"`
	MainID    string             `json:"mainId"`
	Name      string             `json:"name"`
	Icon      string             `json:"icon"`
	Sort      int                `json:"sort"`
	Type      int                `json:"type"`
	Path      string             `json:"path"`
	Expand    bool               `json:"expand,omitempty"`
	Childrens []*TreeNode        `json:"childrens,omitempty"`
}

func BuildSortedTree(data []TreeNode) []*TreeNode {
	nodeMap := make(map[string]*TreeNode)
	var tree []*TreeNode

	// 初始化每個節點
	for _, item := range data {
		copy := item // 避免 reference 問題
		copy.Expand = false
		copy.Childrens = []*TreeNode{}
		nodeMap[item.ID.Hex()] = &copy
	}

	// 建立樹狀關係
	for _, item := range data {
		parentKey := item.MainID
		if parent, ok := nodeMap[parentKey]; ok && parentKey != "" {
			parent.Childrens = append(parent.Childrens, nodeMap[item.ID.Hex()])
		} else {
			tree = append(tree, nodeMap[item.ID.Hex()])
		}
	}

	// 遞迴排序
	var sortTree func(nodes []*TreeNode)
	sortTree = func(nodes []*TreeNode) {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Sort < nodes[j].Sort
		})
		for _, node := range nodes {
			if len(node.Childrens) > 0 {
				sortTree(node.Childrens)
			}
		}
	}

	sortTree(tree)

	return tree
}



