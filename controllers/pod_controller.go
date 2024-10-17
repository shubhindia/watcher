package controllers

import (
	"context"

	"github.com/shubhindia/watcher/config"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// PodWatcherReconciler watches for pods and make sure the intended pod/s are always newer than the existing pods
type PodWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Config *config.Config
}

func (r *PodWatcherReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	if req.Namespace != r.Config.Namespace {
		return reconcile.Result{Requeue: false}, nil
	}

	pod := &corev1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	err := r.checkAndRestart(*pod)
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *PodWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
		}).
		Complete(r)
}

func (r *PodWatcherReconciler) checkAndRestart(pod corev1.Pod) error {
	for _, resource := range r.Config.Newest {
		switch resource.Kind {

		case "Deployment":
			return r.ensureDeployment(pod, resource)

		case "StatefulSet":
			return r.ensureStatefulSet(pod, resource)

		case "DaemonSet":
			// get pods for the daemonset and compare with the pod

		}

	}

	return nil
}
