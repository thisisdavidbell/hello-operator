package controller

import (
	"github.com/thisisdavidbell/hello-operator/operator-sdk-0.18/hello-operator/pkg/controller/hello"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, hello.Add)
}
