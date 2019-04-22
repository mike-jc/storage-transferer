package accessSharers

type Fake struct {
	AbstractSharer
}

func (s *Fake) SharerName() string {
	return "Fake sharer"
}

func (s *Fake) Init() error {
	err := s.AbstractSharer.Init()
	if err != nil {
		return err
	}

	s.self = s

	return err
}
