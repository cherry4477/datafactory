package v1

import (
	"k8s.io/kubernetes/pkg/api"
)

func init() {
	api.Scheme.AddKnownTypes("v1",
		&Application{},
		&ApplicationList{},
	)
}

func (*Application) IsAnAPIObject()     {}
func (*ApplicationList) IsAnAPIObject() {}
