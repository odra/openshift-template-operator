package okdtemplate

import (
	"context"

	odrav1alpha1 "github.com/odra/openshift-template-operator/pkg/apis/odra/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	//"k8s.io/apimachinery/pkg/types"
	"github.com/gobuffalo/packr"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/template"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"github.com/integr8ly/operator-sdk-openshift-utils/pkg/api/kubernetes"
	kube "github.com/odra/openshift-template-operator/pkg/kube/helpers"
	"github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

var log = logf.Log.WithName("controller_okdtemplate")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new OKDTemplate Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOKDTemplate{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		config:   mgr.GetConfig(),
		tmpl:     &template.Tmpl{},
		box:      packr.NewBox("../../../res"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("okdtemplate-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource OKDTemplate
	err = c.Watch(&source.Kind{Type: &odrav1alpha1.OKDTemplate{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner OKDTemplate
	err = c.Watch(&source.Kind{Type: &v1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &odrav1alpha1.OKDTemplate{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileOKDTemplate{}

// ReconcileOKDTemplate reconciles a OKDTemplate object
type ReconcileOKDTemplate struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	config   *rest.Config
	box      packr.Box
	tmpl     *template.Tmpl
}

// Reconcile reads that state of the cluster for a OKDTemplate object and makes changes based on the state read
// and what is in the OKDTemplate.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOKDTemplate) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling OKDTemplate")

	var err error
	// Fetch the OKDTemplate instance
	instance := &odrav1alpha1.OKDTemplate{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.GetDeletionTimestamp() != nil {
		err = r.delete(instance)
		if err != nil {
			reqLogger.Error(err, "delete error")
			return reconcile.Result{}, err
		}
	}

	if !kube.HasFinalizer(instance, "org.odra.DefaultFinalizer") {
		err = r.setFinalizer(instance)
		if err != nil {
			reqLogger.Error(err, "finalizer error")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue:true}, nil
	}

	switch instance.Status.Type {
	case odrav1alpha1.OKDTemplateNone:
		err = r.bootstrap(instance)
		if err != nil {
			reqLogger.Error(err, "bootstrap error")
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	case odrav1alpha1.OKDTemplateNew:
		err = r.install(instance)
		if err != nil {
			reqLogger.Error(err, "install error")
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	case odrav1alpha1.OKDTemplateReconcile:
		err, isReady := r.isReady(instance)
		if err != nil {
			reqLogger.Error(err, "install error")
			return reconcile.Result{}, err
		}

		if !isReady {
			return reconcile.Result{Requeue:true}, nil
		}

		err = r.finish(instance)
		if err != nil {
			reqLogger.Error(err, "install finalization error")
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	case odrav1alpha1.OKDTemplateError:
		return reconcile.Result{}, nil
	case odrav1alpha1.OKDTemplateReady:
		return reconcile.Result{}, nil
	case odrav1alpha1.OKDTemplateDelete:
		err = r.removeFinalizer(instance)
		if err != nil {
			reqLogger.Error(err, "remove finalizer error")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue:true}, nil
	default:
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

//status none -> new
//sets the cr finalizer
//sets the template reader obj
//sets the status to new
func (r *ReconcileOKDTemplate) setFinalizer(cr *odrav1alpha1.OKDTemplate) error {
	finalizer := "org.odra.DefaultFinalizer"
	kube.AddFinalizer(cr, finalizer)
	return r.client.Update(context.TODO(), cr)
}

func (r *ReconcileOKDTemplate) removeFinalizer(cr *odrav1alpha1.OKDTemplate) error {
	finalizer := "org.odra.DefaultFinalizer"
	dc := &v1.DeploymentConfig{}
	key := types.NamespacedName{
		Name: "tutorial-web-app",
		Namespace: cr.Namespace,
	}

	err := r.client.Get(context.TODO(), key, dc)
	if err == nil {
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}


	kube.RemoveFinalizer(cr, finalizer)
	return r.client.Update(context.TODO(), cr)
}

func (r *ReconcileOKDTemplate) bootstrap(cr *odrav1alpha1.OKDTemplate) error {
	yamlData, err := r.box.Find(cr.Spec.Source.Local)
	if err != nil {
		return err
	}

	jsonData, err := yaml.ToJSON(yamlData)
	if err != nil {
		return err
	}

	tmpl, err := template.New(r.config, jsonData)
	if err != nil {
		return err
	}

	r.tmpl = tmpl

	cr.Status = odrav1alpha1.OKDTemplateStatus{
		Type:odrav1alpha1.OKDTemplateNew,
		Message: new(string),
		Reason: new(string),
	}
	*cr.Status.Message = "OKDTemplateNew"
	*cr.Status.Reason = "New OKDTemplate resource found"

	return r.client.Status().Update(context.TODO(), cr)
}

//status new -> reconcile
//reads runtime objects from template lib object
//creates those runtime objects in openshift
//sets status to reconcile
func (r *ReconcileOKDTemplate) install(cr *odrav1alpha1.OKDTemplate) error {
	params := map[string]string{}
	for k, v := range cr.Spec.Parameters {
		params[k] = v
	}

	err := r.tmpl.Process(params, cr.Namespace)
	if err != nil {
		return err
	}

	objects := r.tmpl.GetObjects(template.NoFilterFn)
	for _, ro := range objects {
		uo, err := kubernetes.UnstructuredFromRuntimeObject(ro)
		if err != nil {
			return err
		}

		uo.SetNamespace(cr.Namespace)

		err = controllerutil.SetControllerReference(cr, uo, r.scheme)
		if err != nil {
			return err
		}

		err = r.client.Create(context.TODO(), uo.DeepCopyObject())
		if err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	cr.Status = odrav1alpha1.OKDTemplateStatus{
		Type:odrav1alpha1.OKDTemplateReconcile,
		Message: new(string),
		Reason: new(string),
	}
	*cr.Status.Message = "OKDTemplateReconcile"
	*cr.Status.Reason = "OKDTemplate reconcile loop"

	return r.client.Status().Update(context.TODO(), cr)
}

//checks if deploymentconfig objects are ready and running
func (r *ReconcileOKDTemplate) isReady(cr *odrav1alpha1.OKDTemplate) (error, bool) {
	dc := &v1.DeploymentConfig{}
	key := types.NamespacedName{
		Name: "tutorial-web-app",
		Namespace: cr.Namespace,
	}

	err := r.client.Get(context.TODO(), key, dc)
	if err != nil {
		return err, false
	}

	for _, condition := range dc.Status.Conditions {
		if condition.Type == v1.DeploymentAvailable  && condition.Status == corev1.ConditionTrue {
			return nil, true
		}
	}

	return nil, true
}

//status reconcile -> ready|error
//sets status to ready or error (in case something goes wrong)
func (r *ReconcileOKDTemplate) finish(cr *odrav1alpha1.OKDTemplate) error {
	cr.Status = odrav1alpha1.OKDTemplateStatus{
		Type:odrav1alpha1.OKDTemplateReady,
		Message: new(string),
		Reason: new(string),
	}
	*cr.Status.Message = "OKDTemplateReady"
	*cr.Status.Reason = "OKDTemplate installation finished"

	return r.client.Status().Update(context.TODO(), cr)
}

//status ready -> delete
//removes finalizer once "foreGround" finalizer is gone
//sets status to "delete"
func (r *ReconcileOKDTemplate) delete(cr *odrav1alpha1.OKDTemplate) error {
	cr.Status = odrav1alpha1.OKDTemplateStatus{
		Type:odrav1alpha1.OKDTemplateDelete,
		Message: new(string),
		Reason: new(string),
	}
	*cr.Status.Message = "OKDTemplateDelete"
	*cr.Status.Reason = "OKDTemplate is being deleted"

	return r.client.Status().Update(context.TODO(), cr)
}
