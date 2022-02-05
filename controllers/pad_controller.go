/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	indicatorv1alpha1 "github.com/arjunshenoymec/pad-operator/api/v1alpha1"
)

// PadReconciler reconciles a Pad object
type PadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=indicator.padoperator,resources=pads,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=indicator.padoperator,resources=pads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=indicator.padoperator,resources=pads/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pad object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here

	return ctrl.Result{}, nil
}

// Function to generate a PAD deployment object
func (r *PadReconciler) padDeployment(p *indicatorv1alpha1.Pad) *appsv1.Deployment {
	labels := padLabels(p.Name)
	replicas := p.Spec.Replicas
	image := p.Spec.Image
	source := p.Spec.Source
	metrics := p.Spec.Metrics
	retraining_interval := p.Spec.Retraining_interval
	training_window_size := p.Spec.Training_window_size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.Name,
			Namespace: p.Namespace,

		},
		spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name: "Prometheus-Anomaly-Detector",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "pad",
						}},
						Env: []corev1.EnvVar{
							{
								Name: "FLT_PROM_URL",
								Value: source
							},
							{
								Name: "FLT_PROM_ACCESS_TOKEN",
								Value: "my-access-token"
							},
							{
								Name: "FLT_METRICS_LIST",
								Value: metrics
							},
							{
								Name: "FLT_RETRAINING_INTERVAL_MINUTES",
								Value: retraining_interval
							},
							{
								Name: "FLT_ROLLING_TRAINING_WINDOW_SIZE",
								Value: training_window_size
							},
							{
								Name: "FLT_DEBUG_MODE",
								Value: true
							},
							{
								Name: "APP_FILE",
								Value: "app.py"
							},
						},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(p, deployment, r.Scheme)
	return deployment
}

func padLabels(name string) map[string]string {
        return map[string]string{"app":"pad", "pad_cr": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&indicatorv1alpha1.Pad{}).
		Complete(r)
}
