package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	sidecarInjectLabel = "sidecar/inject"
	sidecarName        = "mySidecar"
	sidecarImage       = "busybox"
)

type SidecarReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *SidecarReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func (r *SidecarReconciler) Reconcile(
	ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// STEP 1: get the deployment object
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		return ctrl.Result{}, err
	}

	// STEP 2: reconcile (THIS IS BAD CODE)

	labels := deployment.GetLabels()

	sidecar := &corev1.Container{
		Name:    sidecarName,
		Image:   sidecarImage,
		Command: []string{"sleep"},
		Args:    []string{"36000"},
	}

	desiredDeployment := deployment.DeepCopy()
	containers := deployment.Spec.Template.Spec.Containers
	// check if inject label present
	if labels[sidecarInjectLabel] == "true" {

		// add sidecar
		containers = append(containers, *sidecar)
		desiredDeployment.Spec.Template.Spec.Containers = containers

		// update
		err := r.Update(ctx, desiredDeployment)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	for i, container := range containers {
		if container.Name == sidecarName &&
			container.Image == sidecarImage {

			// remove sidecar
			desiredContainers := append(containers[:i], containers[i+1:]...)
			desiredDeployment.Spec.Template.Spec.Containers = desiredContainers

			// update
			err := r.Update(ctx, desiredDeployment)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	// STEP 2: reconcile (Good, testable code)
	if err := r.handleDeploymentReconciliation(ctx, deployment); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *SidecarReconciler) handleDeploymentReconciliation(ctx context.Context,
	deployment *appsv1.Deployment) error {
	labels := deployment.GetLabels()

	sidecar := &corev1.Container{
		Name:    sidecarName,
		Image:   sidecarImage,
		Command: []string{"sleep"},
		Args:    []string{"36000"},
	}

	desiredDeployment := deployment.DeepCopy()
	containers := deployment.Spec.Template.Spec.Containers
	// check if inject label present
	if labels[sidecarInjectLabel] == "true" {

		// add sidecar
		containers = append(containers, *sidecar)
		desiredDeployment.Spec.Template.Spec.Containers = containers

		// update
		err := r.Update(ctx, desiredDeployment)
		if err != nil {
			return err
		}
		return nil
	}
	for i, container := range containers {
		if container.Name == sidecarName &&
			container.Image == sidecarImage {

			// remove sidecar
			desiredContainers := append(containers[:i], containers[i+1:]...)
			desiredDeployment.Spec.Template.Spec.Containers = desiredContainers

			// update
			err := r.Update(ctx, desiredDeployment)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
