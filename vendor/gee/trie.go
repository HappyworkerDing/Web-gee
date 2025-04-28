package gee

import "strings"

type node struct {
	pattern  string  //存储完整的路由路径 只有在叶子节点存储，非叶子节点都为空字符串
	part     string  //存储路由路径的部分
	children []*node //存储子节点
	isWild   bool    ///是否准确匹配，当前子节点是否为通配符节点，part含有 ：或 * 时为true “：”表示参数匹配 如：id ， “*”表示通配符匹配，如*filepath
}

// 第一个匹配成功的节点，可以用来插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点,可以用来查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 将路径pattern插入到路由树中，part表示将pattern分成的多个部分，height表示正在处理路径分段parts的第height层
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//如果没找到节点就直接插入一个新的节点
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 在路由树中搜索与parts相匹配的节点
func (n *node) search(parts []string, height int) *node {
	//通配符节点可以直接匹配剩余路径，满足递归结束条件
	//strings.HasPrefix(s , prefix) 用于检查字符串s是否以指定的前缀prefix开头
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part) //返回一个切片
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
