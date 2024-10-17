package controllers

import (
	"context"

	"github.com/shubhindia/watcher/config"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (r *PodWatcherReconciler) ensureStatefulSet(pod corev1.Pod, resource config.ResourceConfig) error {

	podList := &corev1.PodList{}
	// ToDo: Optimise this code, it is not efficient
	// avoid infinite loop by checking if the pod isn't part of the statefulset which we are comparing
	if pod.OwnerReferences != nil && pod.OwnerReferences[0].Kind == "StatefulSet" && pod.OwnerReferences[0].Name == resource.Name {
		return nil
	}

	sts := &appsv1.StatefulSet{}
	err := r.Get(context.Background(), client.ObjectKey{Name: resource.Name, Namespace: r.Config.Namespace}, sts)
	if err != nil {
		return err
	}

	labelSelector := labels.Set(sts.Spec.Selector.MatchLabels)
	listOpts := []client.ListOption{
		client.InNamespace(r.Config.Namespace),
		client.MatchingLabels(labelSelector),
	}
	if err := r.List(context.Background(), podList, listOpts...); err != nil {
		return err
	}

	var stsPods []corev1.Pod
	for _, pod := range podList.Items {
		for _, ownerRef := range pod.OwnerReferences {
			if ownerRef.Kind == "StatefulSet" && ownerRef.UID == sts.UID {
				stsPods = append(stsPods, pod)
			}
		}
	}

	// Compare the pod with the statefulset pods
	for _, stsPod := range stsPods {
		if pod.CreationTimestamp.After(stsPod.CreationTimestamp.Time) {
			r.Delete(context.Background(), &stsPod)
		}
	}

	return nil

}
