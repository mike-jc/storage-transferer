package accessSharers

import "service-recordingStorage/models/data/accessSharer"

type AbstractSharer struct {
	keyPair  modelsDataAccessSharer.KeyPair
	password string

	self AbstractSharerContract

	AbstractSharerContract
}

type AbstractSharerContract interface {
	SharerName() string
	Status() string
	Init() error
	Reset()
	SetKeyPair(pair modelsDataAccessSharer.KeyPair)
	SetPassword(password string)
	Share(options interface{}) error
}

func (s *AbstractSharer) SharerName() string {
	return "Abstract sharer"
}

func (s *AbstractSharer) Status() string {
	return "access sharing by " + s.self.SharerName()
}

func (s *AbstractSharer) Init() error {
	s.self = s
	return nil
}

func (s *AbstractSharer) Reset() {

}

func (s *AbstractSharer) SetKeyPair(pair modelsDataAccessSharer.KeyPair) {
	s.keyPair = pair
}

func (s *AbstractSharer) SetPassword(password string) {
	s.password = password
}

// @Description share access to all room/folder users
func (s *AbstractSharer) Share(options interface{}) error {
	return nil
}
