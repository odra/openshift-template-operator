package controller

import (
	"github.com/odra/openshift-template-operator/pkg/controller/okdtemplate"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, okdtemplate.Add)
}
