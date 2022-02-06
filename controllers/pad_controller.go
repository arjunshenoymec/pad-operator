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
	indicatorv1alpha1 "github.com/arjunshenoymec/pad-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// PadReconciler reconciles a Pad object
type PadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=indicator.padoperator,resources=pads,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=indicator.padoperator,resources=pads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=indicator.padoperator,resources=pads/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

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
	log := ctrllog.FromContext(ctx)

	pad := &indicatorv1alpha1.Pad{}
	// attempting to get an object of kind Pad in the specified namespace
	err := r.Get(ctx, req.NamespacedName, pad)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Pad resource not found. Ignoring activity")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get the Pad object")
		return ctrl.Result{}, err
	}

	found := &appsv1.Deployment{}
	// atempting to get the deployment
	err = r.Get(ctx, types.NamespacedName{Name: pad.Name, Namespace: pad.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// if the deployment does not exist, we create a new one
		dep := r.padDeployment(pad)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		// Failed to retrieve deployment for some unknown reason
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Deployment exists, checking the individual spec parameters to see if they match
	replicas := pad.Spec.Replicas
	image := pad.Spec.Image
	newEnvVars := generateEnvVars(pad.Spec.Source, pad.Spec.Metrics, pad.Spec.Retraining_interval, pad.Spec.Training_window_size)
	update_flag := false
	pad_container := (*found).Spec.Template.Spec.Containers[0]

	// Checking if there is mismatch in the replicas specified
	if *found.Spec.Replicas != replicas {
		update_flag = true
		found.Spec.Replicas = &replicas
	}

	// Checking if there is mismatch in the image of container
	if pad_container.Image != image {
		update_flag = true
		found.Spec.Template.Spec.Containers[0].Image = image
	}

	// Checking if any of the envrionment variables have been changed
	if reflect.DeepEqual(pad_container.Env, newEnvVars) == false {
		update_flag = true
		found.Spec.Template.Spec.Containers[0].Env = newEnvVars
	}

	if update_flag == true {
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

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
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
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
						Name:  "prometheus-anomaly-detector",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "pad",
						}},
						Env: generateEnvVars(source, metrics, retraining_interval, training_window_size),
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(p, deployment, r.Scheme)
	return deployment
}

func padLabels(name string) map[string]string {
	return map[string]string{"app": "pad", "pad_cr": name}
}

func generateEnvVars(source string, metric_list string, retraining_interval string, training_window_size string) []corev1.EnvVar {
	container_vars := []corev1.EnvVar{
		{
			Name:  "FLT_PROM_URL",
			Value: source,
		},
		{
			Name:  "FLT_PROM_ACCESS_TOKEN",
			Value: "my-access-token",
		},
		{
			Name:  "FLT_METRICS_LIST",
			Value: metric_list,
		},
		{
			Name:  "FLT_RETRAINING_INTERVAL_MINUTES",
			Value: retraining_interval,
		},
		{
			Name:  "FLT_ROLLING_TRAINING_WINDOW_SIZE",
			Value: training_window_size,
		},
		{
			Name:  "FLT_DEBUG_MODE",
			Value: "True",
		},
		{
			Name:  "APP_FILE",
			Value: "app.py",
		},
	}
	return container_vars
}

// SetupWithManager sets up the controller with the Manager.
func (r *PadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&indicatorv1alpha1.Pad{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
