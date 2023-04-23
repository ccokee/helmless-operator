package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	helmv3action "helm.sh/helm/v3/pkg/action"
	helmv3chart "helm.sh/helm/v3/pkg/chart"
	helmv3loader "helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	helmv3cli "helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/ccokee/helmless-operator/api/v1alpha1"
)

const (
	// Reconcile interval to requeue chart CRD.
	reconcileInterval = 30 * time.Second
)

// HelmlessReconciler reconciles a HelmLess object
type HelmlessReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=helmless.redrvm.cloud,resources=helmlesss,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=helmless.redrvm.cloud,resources=helmlesss/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=helmless.redrvm.cloud,resources=helmlesss/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps;secrets;services;pods;serviceaccounts,verbs=get;list;watch;create;update;patch;delete

func (r *HelmlessReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Get the HelmLess object
	helmless := &v1alpha1.HelmLess{}
	if err := r.Client.Get(ctx, req.NamespacedName, helmless); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	settings := helmv3cli.New()
	actionConfig := new(helmv3action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), helmless.Namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v)
	}); err != nil {
		return ctrl.Result{}, err
	}

	chartRequested, err := helmv3loader.Load(helmless.Spec.ChartRepo)
	if err != nil {
		return ctrl.Result{}, err
	}

	values, err := fetchValues(helmless.Spec.ValuesUrl)
	if err != nil {
		return ctrl.Result{}, err
	}

	install := helmv3action.NewInstall(actionConfig)
	install.ReleaseName = helmless.Spec.ChartReleaseName
	install.Namespace = helmless.Namespace
	_, err = install.Run(chartRequested, values)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: reconcileInterval}, nil
}

func fetchValues(valuesURL string) (map[string]interface{}, error) {
	resp, err := http.Get(valuesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch values from %s, status code: %d", valuesURL, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rawValues := make(map[string]interface{})
	if err := yaml.Unmarshal(body, &rawValues); err != nil {
		return nil, err
	}

	chart := &helmv3chart.Chart{}
	coalescedValues, err := chartutil.CoalesceValues(chart, rawValues)
	if err != nil {
		return nil, err
	}

	return coalescedValues, nil
}

func (r *HelmlessReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.HelmLess{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ServiceAccount{}).
		Complete(r)
}
