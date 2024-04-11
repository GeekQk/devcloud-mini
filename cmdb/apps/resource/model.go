package resource

import (
	"github.com/rs/xid"
)

func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

func NewResource() *Resource {
	return &Resource{
		Meta: &Meta{},
		Spec: &Spec{
			Extra: map[string]string{},
		},
		Status: &Status{},
	}
}

func (r *Resource) AutoFill() {
	if r.Meta.Id == "" {
		r.Meta.Id = xid.New().String()
	}
}
