package apis

import (
	"github.com/odra/openshift-template-operator/pkg/apis/odra/v1alpha1"
	okdSchemes "github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/schemes"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		v1alpha1.SchemeBuilder.AddToScheme,
			okdSchemes.AddToScheme)
}
