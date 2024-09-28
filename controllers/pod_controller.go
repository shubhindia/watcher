package controllers

import (
	"context"

	"github.com/shubhindia/watcher/config"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
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
		Complete(r)
}

func (r *PodWatcherReconciler) checkAndRestart(pod corev1.Pod) error {
	podList := &corev1.PodList{}
	for _, resource := range r.Config.Newest {
		switch resource.Kind {

		case "Deployment":
			// get pods for the deployment and compare with the pod
			deployment := &appsv1.Deployment{}
			err := r.Get(context.Background(), client.ObjectKey{Name: resource.Name, Namespace: r.Config.Namespace}, deployment)
			if err != nil {
				return err
			}

			// Get the label selector for the deployment
			labelSelector := labels.Set(deployment.Spec.Selector.MatchLabels)
			listOpts := []client.ListOption{
				client.InNamespace(r.Config.Namespace),
				client.MatchingLabels(labelSelector),
			}
			if err := r.List(context.Background(), podList, listOpts...); err != nil {
				return err
			}
			var deploymentPods []corev1.Pod
			for _, pod := range podList.Items {
				for _, ownerRef := range pod.OwnerReferences {
					if ownerRef.Kind == "ReplicaSet" && ownerRef.UID == deployment.UID {
						deploymentPods = append(deploymentPods, pod)
					}
				}
			}

			// Compare the pod with the deployment pods
			for _, deploymentPod := range deploymentPods {
				if pod.CreationTimestamp.After(deploymentPod.CreationTimestamp.Time) {
					r.Delete(context.Background(), &deploymentPod)
				}
			}

		case "StatefulSet":
			// get pods for the statefulset and compare with the pod
		case "DaemonSet":
			// get pods for the daemonset and compare with the pod

		}

	}

	return nil
}
