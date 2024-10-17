package controllers

import (
	"context"
	"strings"

	"github.com/shubhindia/watcher/config"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// EnsureDeployment checks if the given pod is newer than the deployments and returns true if it is newer
func (r *PodWatcherReconciler) EnsureDeployment(pod corev1.Pod, resource config.ResourceConfig) error {

	podList := &corev1.PodList{}
	// ToDo: Optimise this code, it is not efficient
	// avoid infinite loop by checking if the pod isn't part of the statefulset which we are comparing
	if pod.OwnerReferences != nil && pod.OwnerReferences[0].Kind == "ReplicaSet" && strings.Contains(pod.OwnerReferences[0].Name, resource.Name) {
		return nil
	}

	// get pods for the deployment and compare with the pod
	deployment := &appsv1.Deployment{}
	err := r.Get(context.Background(), client.ObjectKey{Name: resource.Name, Namespace: r.Config.Namespace}, deployment)
	if err != nil {
		return err
	}

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

	for _, deploymentPod := range deploymentPods {
		if pod.CreationTimestamp.After(deploymentPod.CreationTimestamp.Time) {
			r.Delete(context.Background(), &deploymentPod)
		}
	}

	return nil
}
