package modelsDataInstance

import "fmt"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) NameAndEmail() string {
	if u.Name != "" && u.Email != "" {
		return fmt.Sprintf("%s (%s)", u.Name, u.Email)
	} else if u.Name != "" {
		return u.Name
	} else if u.Email != "" {
		return u.Email
	}
	return ""
}
