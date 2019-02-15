package helpers

import (
	"github.com/odra/openshift-template-operator/pkg/kube"
)

func HasFinalizer(obj kube.FinalizerSpec, name string) bool {
	finalizers := obj.GetFinalizers()

	for _, finalizer := range finalizers {
		if finalizer == name {
			return true
		}
	}

	return false
}

func AddFinalizer(obj kube.FinalizerSpec, name string) {
	finalizers := append(obj.GetFinalizers(), name)
	obj.SetFinalizers(finalizers)
}

func RemoveFinalizer(obj kube.FinalizerSpec, name string) {
	finalizers := make([]string, 0)
	for _, finalizer := range obj.GetFinalizers() {
		if finalizer != name {
			finalizers = append(finalizers, finalizer)
		}
	}

	obj.SetFinalizers(finalizers)
}