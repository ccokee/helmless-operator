package controllers

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	helmlessv1alpha1 "github.com/ccokee/helmless-operator/api/v1alpha1"
)

type HelmLessReconciler struct {
	client.Client
	Log    logr.Logger
	DynCli dynamic.Interface
}

//+kubebuilder:rbac:groups=helmless.redrvm.cloud,resources=helmless,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=helmless.redrvm.cloud,resources=helmless/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=helmless.redrvm.cloud,resources=helmless/finalizers,verbs=update

func (r *HelmLessReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("helmless", req.NamespacedName)

	instance := &helmlessv1alpha1.HelmLess{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		// handle error
	}

	valuesFileContent, err := getValuesFromGist(instance.Spec.PublicGist)
	if err != nil {
		log.Error(err, "Failed to get values from Gist")
		return ctrl.Result{}, err
	}

	err = r.applyResources(instance.Spec.Namespace, valuesFileContent)
	if err != nil {
		log.Error(err, "Failed to apply resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *HelmLessReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.DynCli = dynamic.NewForConfigOrDie(mgr.GetConfig())

	return ctrl.NewControllerManagedBy(mgr).
		For(&helmlessv1alpha1.HelmLess{}).
		Complete(r)
}

func (r *HelmLessReconciler) applyResources(namespace string, resourcesYAML string) error {
	resources := strings.Split(resourcesYAML, "---")

	for _, resource := range resources {
		if len(strings.TrimSpace(resource)) == 0 {
			continue
		}

		unstructuredObj := &unstructured.Unstructured{}
		err := yaml.Unmarshal([]byte(resource), unstructuredObj)
		if err != nil {
			return err
		}

		gvk := unstructuredObj.GroupVersionKind()
		gvr := schema.GroupVersionResource{
			Group:    gvk.Group,
			Version:  gvk.Version,
			Resource: gvk.Kind,
		}

		unstructuredObj.SetNamespace(namespace)

		_, err = r.DynCli.Resource(gvr).Namespace(namespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func getValuesFromGist(gistURL string) (string, error) {
	resp, err := http.Get(gistURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
