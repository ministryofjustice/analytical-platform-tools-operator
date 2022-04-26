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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

// JupyterLabReconciler reconciles a JupyterLab object
type JupyterLabReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=jupyterlabs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=jupyterlabs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=jupyterlabs/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// +kubebuilder:rbac:groups=networking,resources=ingress,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the JupyterLab object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *JupyterLabReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the JupyterLab instance
	jupyterlab := &toolsv1alpha1.JupyterLab{}
	err := r.Get(ctx, req.NamespacedName, jupyterlab)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("JupyterLab resource not found", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	// Check if the deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Namespace: jupyterlab.Namespace, Name: jupyterlab.Name}, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("Creating a new Deployment", "namespace", jupyterlab.Namespace, "name", jupyterlab.Name)
			dep := r.deploymentForJupyterLab(ctx, jupyterlab)
			err = r.Create(ctx, dep)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		} else {
			log.Log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *JupyterLabReconciler) deploymentForJupyterLab(ctx context.Context, jupyterlab *toolsv1alpha1.JupyterLab) *appsv1.Deployment {
	image := jupyterlab.Spec.Image
	version := jupyterlab.Spec.Version

	if image == "" {
		image = "jupyter/minimal-notebookest"
	}

	if version == "" {
		version = "latest"
	}

	lab := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jupyterlab.Name,
			Namespace: jupyterlab.Namespace,
			Labels: map[string]string{
				"app":                          "jupyter",
				"chart":                        jupyterlab.Name,
				"app.kubernetes.io/managed-by": "jupyterlab-operator",
				"app.kubernetes.io/component":  "jupyterlab",
				"app.kubernetes.io/part-of":    "jupyterlab",
				"app.kubernetes.io/version":    jupyterlab.Spec.Version,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "jupyterlab"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "jupyterlab"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "jupyterlab",
							Image: image + ":" + version,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8888,
								},
							},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(jupyterlab, lab, r.Scheme)
	return lab
}

// SetupWithManager sets up the controller with the Manager.
func (r *JupyterLabReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.JupyterLab{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
