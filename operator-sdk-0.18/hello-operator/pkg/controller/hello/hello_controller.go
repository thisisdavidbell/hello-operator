package hello

import (
	"context"
	"reflect"
	"strconv"

	"github.com/go-logr/logr"
	thisisdavidbellv1alpha1 "github.com/thisisdavidbell/hello-operator/operator-sdk-0.18/hello-operator/pkg/apis/thisisdavidbell/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_hello")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Hello Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileHello{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("hello-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Hello
	err = c.Watch(&source.Kind{Type: &thisisdavidbellv1alpha1.Hello{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Hello
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &thisisdavidbellv1alpha1.Hello{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileHello implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileHello{}

// ReconcileHello reconciles a Hello object
type ReconcileHello struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Hello object and makes changes based on the state read
// and what is in the Hello.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHello) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Hello")

	// Fetch the Hello instance
	helloInstance := &thisisdavidbellv1alpha1.Hello{}
	err := r.client.Get(context.TODO(), request.NamespacedName, helloInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Hello cr not found. Presume deleted. Reconcile complete.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	result, err := reconcileDeployment(helloInstance, r, reqLogger)
	if err != nil || result.Requeue == true {
		// reconcile requested requeue or errored, so requeue
		return result, err
	}

	// got to end. Reconcile completed and was successful.
	reqLogger.Info("End of successful Reconcile.")
	return reconcile.Result{}, nil
}

// reconcileDeployment creates or updates the k8s deployment based on the cr.
func reconcileDeployment(helloInstance *thisisdavidbellv1alpha1.Hello, r *ReconcileHello, reqLogger logr.Logger) (reconcile.Result, error) {
	// Define a new Deployment Object
	deployment := newDeploymentForCR(helloInstance)

	// Set Hello instance as the owner and controller
	if err := controllerutil.SetControllerReference(helloInstance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists
	foundDeployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Deployment created successfully - requeue
		reqLogger.Info("Successfully created deployment. Reconcile complete.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		// get deployment failed, and not with NotFound
		return reconcile.Result{}, err
	}

	// Deployment already exists
	reqLogger.Info("Deployment already exists. Check state matches desired...", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)

	// Ensure hello app version is correct
	desiredImage := getHelloImage(helloInstance)
	foundImage := foundDeployment.Spec.Template.Spec.Containers[0].Image
	if desiredImage != foundImage {
		reqLogger.Info("Hello image version mismatch found.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		foundDeployment.Spec.Template.Spec.Containers[0].Image = desiredImage
		err = r.client.Update(context.TODO(), foundDeployment)
		if err != nil {
			log.Error(err, "Failed to update hello image version in deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return reconcile.Result{}, err
		}
		// Deployment spec updated - return and requeue
		reqLogger.Info("Successfully updated hello image version. Requeue to continue.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		return reconcile.Result{Requeue: true}, nil
	}

	// Ensure repeat and verbose env vars are correct
	desiredEnvVars := getPodEnvVars(helloInstance)
	foundEnvVars := foundDeployment.Spec.Template.Spec.Containers[0].Env
	if !reflect.DeepEqual(desiredEnvVars, foundEnvVars) {
		reqLogger.Info("Env vars mismatch found.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		foundDeployment.Spec.Template.Spec.Containers[0].Env = desiredEnvVars
		err = r.client.Update(context.TODO(), foundDeployment)
		if err != nil {
			log.Error(err, "Failed to update env vars in deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return reconcile.Result{}, err
		}
		// Deployment spec updated - return and requeue
		reqLogger.Info("Successfully updated env vars. Requeue to continue.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)

		return reconcile.Result{Requeue: true}, nil
	}

	// deployment successfully reconciled. Continue Reconcile() loop
	return reconcile.Result{}, nil
}

// newDeploymentForCR returns a deployment with the same name/namespace as the cr
// code came from memcache example deployment inline here: https://docs.openshift.com/container-platform/4.6/operators/operator_sdk/osdk-getting-started.html
func newDeploymentForCR(cr *thisisdavidbellv1alpha1.Hello) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}

	envs := getPodEnvVars(cr)

	// Note: currently spec.version is a required field, so will have value. Its now an enum which only accepts valid values, so we can be confident its always valid.
	image := getHelloImage(cr)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "hello",
							Image:   image,
							Command: []string{"./hello"},
							Env:     envs,
						},
					},
				},
			},
		},
	}
}

func getHelloImage(cr *thisisdavidbellv1alpha1.Hello) string {
	return "SET_TO_IRHOSTNAME/SET_TO_IRNAMESPACE/hello:" + cr.Spec.Version
}

func getPodEnvVars(cr *thisisdavidbellv1alpha1.Hello) []corev1.EnvVar {

	// Note: currently these are required fields, so will exist and will have valid int and bool values
	repeat := strconv.Itoa(cr.Spec.Repeat)
	verbose := strconv.FormatBool(cr.Spec.Verbose)

	return []corev1.EnvVar{
		{
			Name:  "REPEAT",
			Value: repeat,
		},
		{
			Name:  "VERBOSE",
			Value: verbose,
		},
	}
}
