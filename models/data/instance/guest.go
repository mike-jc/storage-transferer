package modelsDataInstance

import "fmt"

type Guest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (g *Guest) NameAndEmail() string {
	if g.Name != "" && g.Email != "" {
		return fmt.Sprintf("%s (%s)", g.Name, g.Email)
	} else if g.Name != "" {
		return g.Name
	} else if g.Email != "" {
		return g.Email
	}
	return ""
}
