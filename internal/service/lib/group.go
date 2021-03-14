package lib

// Group
type Group struct {
	BaseItem
	Cookbooks []string `json:"cookbooks"`
	Groups    []string `json:"groups"`
}

// AddCookbook Add cookbook id to group
func (g Group) AddCookbook(id string) {
	g.Cookbooks = append(g.Cookbooks, id)
}

// AddGroup add a group to another
func (g Group) AddGroup(id string) {
	g.Groups = append(g.Groups, id)
}
