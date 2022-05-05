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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	toolsv1alpha1 "github.com/ministryofjustice/analytical-platform-tools-operator/api/v1alpha1"
)

// ToolReconciler reconciles a Tool object
type ToolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=tools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=tools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice,resources=tools/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ToolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Check the status of the tool
	tool := &toolsv1alpha1.Tool{}
	if err := r.Get(ctx, req.NamespacedName, tool); err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("tool resource not found")
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Failed to get Tool resource")
		return ctrl.Result{}, nil
	}

	var resource interface{}
	switch strings.ToLower(tool.Name) {
	case "jupyterlab":
		log.Log.Info("Initiating JupyterLab instance for reconcile")
		resource = &toolsv1alpha1.JupyterLab{}
	case "rstudio":
		log.Log.Info("Initiating Rstudio instance for reconcile")
		resource = &toolsv1alpha1.Rstudio{}
	case "airflow":
		log.Log.Info("Initiating Airflow instance for reconcile")
		resource = &toolsv1alpha1.Airflow{}
	default:
		log.Log.Info("Reconciling unknown tool, FAIL")
		return ctrl.Result{}, fmt.Errorf("unknown tool %s", tool.Name)
	}

	obj, dep := r.getResource(resource, tool)
	err := r.Get(ctx, types.NamespacedName{Name: tool.Name, Namespace: tool.Namespace}, obj)
	if err != nil && errors.IsNotFound(err) {
		err := r.Create(ctx, dep)
		if err != nil {
			log.Log.Error(err, "Failed to create tool resource")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// TODO: Update status of the tool

	return ctrl.Result{}, nil
}

// getResource extracts the resource from an interface
func (r *ToolReconciler) getResource(obj interface{}, tool *toolsv1alpha1.Tool) (client.Object, client.Object) {
	switch obj.(type) {
	case *toolsv1alpha1.JupyterLab:
		return &toolsv1alpha1.JupyterLab{}, r.createJupyterLabDeployment(tool)
	case *toolsv1alpha1.Rstudio:
		return &toolsv1alpha1.Rstudio{}, r.createRstudioDeployment(tool)
	case *toolsv1alpha1.Airflow:
		return &toolsv1alpha1.Airflow{}, r.createAirflowDeployment(tool)
	default:
		return nil, nil
	}
}

func (r *ToolReconciler) createJupyterLabDeployment(tool *toolsv1alpha1.Tool) *toolsv1alpha1.JupyterLab {
	jupyterlab := &toolsv1alpha1.JupyterLab{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tool.Name,
			Namespace: tool.Namespace,
		},
		Spec: toolsv1alpha1.JupyterLabSpec{
			Image:   tool.Spec.Image,
			Version: tool.Spec.ImageVersion,
		},
	}

	ctrl.SetControllerReference(tool, jupyterlab, r.Scheme)
	return jupyterlab
}

func (r *ToolReconciler) createRstudioDeployment(tool *toolsv1alpha1.Tool) *toolsv1alpha1.Rstudio {
	rstudio := &toolsv1alpha1.Rstudio{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tool.Name,
			Namespace: tool.Namespace,
		},
		Spec: toolsv1alpha1.RstudioSpec{
			Image:   tool.Spec.Image,
			Version: tool.Spec.ImageVersion,
		},
	}
	ctrl.SetControllerReference(tool, rstudio, r.Scheme)
	return rstudio
}

func (r *ToolReconciler) createAirflowDeployment(tool *toolsv1alpha1.Tool) *toolsv1alpha1.Airflow {
	airflow := &toolsv1alpha1.Airflow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tool.Name,
			Namespace: tool.Namespace,
		},
		Spec: toolsv1alpha1.AirflowSpec{
			Image:   tool.Spec.Image,
			Version: tool.Spec.ImageVersion,
		},
	}

	ctrl.SetControllerReference(tool, airflow, r.Scheme)
	return airflow
}

// SetupWithManager sets up the controller with the Manager.
func (r *ToolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.Tool{}).
		Owns(&toolsv1alpha1.JupyterLab{}).
		Owns(&toolsv1alpha1.Rstudio{}).
		Owns(&toolsv1alpha1.Airflow{}).
		Complete(r)
}
