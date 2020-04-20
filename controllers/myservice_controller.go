/*
Copyright 2020 Hirotaka Yamamoto.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	ssav1 "github.com/ymmt2005/kubebuilder-ssa/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MyServiceReconciler reconciles a MyService object
type MyServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ssa.ymmt2005.github.io,resources=myservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ssa.ymmt2005.github.io,resources=myservices/status,verbs=get;update;patch

// Reconcile ...
func (r *MyServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("myservice", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager ...
func (r *MyServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ssav1.MyService{}).
		Complete(r)
}
