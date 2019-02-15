package kube


type FinalizerSpec interface {
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
}