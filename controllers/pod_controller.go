package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// PodReconciler watches for pods and make sure the intended pod/s are always newer than the existing pods
type PodWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PodWatcherReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	pod := &corev1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Printf("Pod %s/%s has been updated\n", pod.Namespace, pod.Name)
	return reconcile.Result{}, nil
}

func (r *PodWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}). // Watch for Pod updates
		Complete(r)
}
