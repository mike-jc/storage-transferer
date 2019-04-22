package modelsDataUniverse

type Instance struct {
	Id       int
	DomainId string
}

func (o *Instance) GetDomainId() string {
	return o.DomainId
}

func (o *Instance) GetId() int {
	return o.Id
}

func (o *Instance) SetId(id int) {
	o.Id = id
}

func (o *Instance) SetDomainId(domainId string) {
	o.DomainId = domainId
}
