package apis

import (
	"github.com/thisisdavidbell/hello-operator/operator-sdk-0.18/hello-operator/pkg/apis/thisisdavidbell/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
