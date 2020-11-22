package model

import "github.com/nilorg/sdk/convert"

// ResultSelect 返回select
type ResultSelect struct {
	Value interface{} `json:"value"`
	Label string      `json:"label"`
}

// ResultTreeSelect 返回TreeSelect
type ResultTreeSelect struct {
	Title    string              `json:"title"`
	Value    interface{}         `json:"value"`
	Key      string              `json:"key"`
	Children []*ResultTreeSelect `json:"children"`
}

// RecursiveRoleToTreeSelect 递归role转tree select
func RecursiveRoleToTreeSelect(roles []*Role) (treeSelects []*ResultTreeSelect) {
	for _, role := range roles {
		treeSelect := new(ResultTreeSelect)
		treeSelect.Title = role.Name
		treeSelect.Key = string(role.Code)
		treeSelect.Value = role.Code
		if len(role.ChildRoles) > 0 {
			treeSelect.Children = RecursiveRoleToTreeSelect(role.ChildRoles)
		} else {
			treeSelect.Children = make([]*ResultTreeSelect, 0)
		}
		treeSelects = append(treeSelects, treeSelect)
	}
	return
}

// RecursiveResourceWebMenuToTreeSelect 递归web菜单转tree select
func RecursiveResourceWebMenuToTreeSelect(menus []*ResourceWebMenu) (treeSelects []*ResultTreeSelect) {
	for _, menu := range menus {
		treeSelect := new(ResultTreeSelect)
		treeSelect.Title = menu.Name
		treeSelect.Key = convert.ToString(menu.ID)
		treeSelect.Value = menu.ID
		if len(menu.ChildResourceWebMenus) > 0 {
			treeSelect.Children = RecursiveResourceWebMenuToTreeSelect(menu.ChildResourceWebMenus)
		} else {
			treeSelect.Children = make([]*ResultTreeSelect, 0)
		}
		treeSelects = append(treeSelects, treeSelect)
	}
	return
}

// ResultTreeNode 返回 tree node
type ResultTreeNode struct {
	ID     string      `json:"id"`
	PID    string      `json:"pId"`
	Value  interface{} `json:"value"`
	Title  string      `json:"title"`
	IsLeaf bool        `json:"isLeaf"` // 是否是叶子节点，叶子节点没有子节点数据
}

// RecursiveRoleToTreeNode 递归role转tree node
func RecursiveRoleToTreeNode(roles []*Role) (treeNodes []*ResultTreeNode) {
	for _, role := range roles {
		treeNode := new(ResultTreeNode)
		treeNode.ID = string(role.Code)
		treeNode.PID = string(role.ParentCode)
		treeNode.Value = role.Code
		treeNode.Title = role.Name
		treeNode.IsLeaf = len(role.ChildRoles) == 0
		treeNodes = append(treeNodes, treeNode)

		if len(role.ChildRoles) > 0 {
			nodes := RecursiveRoleToTreeNode(role.ChildRoles)
			treeNodes = append(treeNodes, nodes...)
		}
	}
	return
}

// ResultTree 返回Tree
type ResultTree struct {
	Title    string        `json:"title"`
	Value    interface{}   `json:"value"`
	Key      interface{}   `json:"key"`
	Children []*ResultTree `json:"children"`
}

// RecursiveOrganizationToTree 递归organization转tree
func RecursiveOrganizationToTree(orgs []*Organization) (result []*ResultTree) {
	for _, org := range orgs {
		treeNode := new(ResultTree)
		treeNode.Title = org.Name
		treeNode.Key = org.ID
		treeNode.Value = org.ID
		if len(org.ChildOrganizations) > 0 {
			treeNode.Children = RecursiveOrganizationToTree(org.ChildOrganizations)
		} else {
			treeNode.Children = make([]*ResultTree, 0)
		}
		result = append(result, treeNode)
	}
	return
}
