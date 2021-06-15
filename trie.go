package lee

import (
	"strings"
)

type Node struct {
	Pattern  string  // 待匹配的路由
	part     string  //路由节点
	children []*Node //子节点
	isWild   bool    //是否精准匹配
}

func (this *Node) MatchChild(part string) *Node {
	for _, child := range this.children {
		if child == nil {
			return nil
		}
		if child.part == part || this.isWild { //匹配到了或者模糊匹配
			return child
		}
	}
	return nil
}

func (this *Node) MatchChildren(part string) []*Node {
	childrens := make([]*Node, 0)
	for _, child := range this.children {
		if child == nil {
			return nil
		}
		if child.part == part || child.isWild { //匹配到了或者模糊匹配
			childrens = append(childrens, child)
		}
	}
	return childrens
}

func (this *Node) Insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		this.Pattern = pattern
		return
	}
	part := parts[height]
	child := this.MatchChild(part)
	if child == nil {
		child = &Node{part: part, isWild: strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*")}
		//当插入的路由为跟路由时，part 为空 这里应该会报空指针
		this.children = append(this.children, child)
	}
	child.Insert(pattern, parts, height+1)
}

func (this *Node) Search(parts []string, height int) *Node {
	if len(parts) == height {
		if this.Pattern == "" {
			return nil
		} else {
			return this
		}

	}
	if len(parts) == 0 {
		return this
	}
	children := this.MatchChildren(parts[height])

	for _, child := range children {
		result := child.Search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil

}
