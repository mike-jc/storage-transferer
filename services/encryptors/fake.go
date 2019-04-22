package encryptors

type Fake struct {
	AbstractEncryptor
}

func (e *Fake) SharerName() string {
	return "Fake encryptor"
}

func (e *Fake) Init(options interface{}) error {
	err := e.AbstractEncryptor.Init(options)
	if err != nil {
		return err
	}

	e.self = e

	return err
}
