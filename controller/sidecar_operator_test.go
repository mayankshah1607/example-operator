package controller

import (
	"context"
	"testing"

	"github.com/mayankshah1607/sidecar-operator/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestHandleDeploymentReconciler(t *testing.T) {
	client := utils.NewClient()

	// setup expectations
	client.On("Update",
		mock.IsType(context.Background()),
		mock.IsType(&appsv1.Deployment{}),
		mock.Anything,
	).Return(nil)

	ctx := context.Background()
	reconciler := &SidecarReconciler{
		Client: client,
		Scheme: newTestScheme(),
	}

	err := reconciler.handleDeploymentReconciliation(ctx, newTestDeployment())
	require.NoError(t, err)
	client.AssertExpectations(t)

}

func newTestDeployment() *appsv1.Deployment {
	labels := make(map[string]string)
	labels[sidecarInjectLabel] = "true"
	return &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:   "test-deployment",
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
}

func newTestScheme() *runtime.Scheme {
	testScheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(testScheme)
	return testScheme
}
