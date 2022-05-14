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
	"github.com/daviderli614/vm-controller-operator/resource"
	"github.com/daviderli614/vm-controller-operator/utils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vmcontrollerv1 "github.com/daviderli614/vm-controller-operator/api/v1"
)

// VMClaimReconciler reconciles a VMClaim object
type VMClaimReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
}

//+kubebuilder:rbac:groups=vmcontroller.uk8s.ucloud.cn,resources=vmclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vmcontroller.uk8s.ucloud.cn,resources=vmclaims/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vmcontroller.uk8s.ucloud.cn,resources=vmclaims/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VMClaim object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *VMClaimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("vmclaim", req.NamespacedName)

	klog.V(4).Info("Reconciling vm controller")

	instance := &vmcontrollerv1.VMClaim{}
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("VMClaim %s not found, will not reconcile", req.Name)
			return ctrl.Result{}, nil
		}
		klog.Errorf("Error reading VMClaim: %v", err)
		return ctrl.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	deploy := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deploy); err != nil {
		if errors.IsNotFound(err) {
			if err := r.CreateVMController(instance); err != nil {
				errMsg := fmt.Sprintf("Error creating vm controller deployment: %v", err)
				klog.Error(errMsg)
				return reconcile.Result{}, err
			}
			klog.Info("Creating vm controller deployment successful")
		} else {
			return ctrl.Result{}, err
		}
	} else {
		if err := r.UpdateVMController(instance); err != nil {
			errMsg := fmt.Sprintf("Error updating vm controller deployment: %v", err)
			klog.Error(errMsg)
			return reconcile.Result{}, err
		}
		klog.Info("Updated vm controller deployment successful")
	}

	return ctrl.Result{}, nil
}

func (r *VMClaimReconciler) CreateVMController(vmc *vmcontrollerv1.VMClaim) error {
	klog.Info("Creating vm controller deployment")

	deployment := resource.NewDeployment(vmc)

	// Set ClusterAutoscaler instance as the owner and controller.
	if err := controllerutil.SetControllerReference(vmc, deployment, r.Scheme); err != nil {
		return err
	}

	return r.Client.Create(context.TODO(), deployment)
}

func (r *VMClaimReconciler) UpdateVMController(vmc *vmcontrollerv1.VMClaim) error {
	klog.Info("Update vm controller deployment")

	existingDeployment, err := r.GetVMController()
	if err != nil {
		return err
	}

	existingSpec := existingDeployment.Spec.Template.Spec
	expectedSpec := resource.VMControllerPodSpec(vmc)

	// Only comparing podSpec and release version for now.
	if equality.Semantic.DeepEqual(existingSpec, expectedSpec) {
		return nil
	}

	existingDeployment.Spec.Template.Spec = *expectedSpec

	return r.Client.Update(context.TODO(), existingDeployment)
}

func (r *VMClaimReconciler) GetVMController() (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	n := r.VMControllerName()
	if err := r.Client.Get(context.TODO(), n, deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}

func (r *VMClaimReconciler) VMControllerName() types.NamespacedName {
	return types.NamespacedName{
		Name:      utils.VMControllerName,
		Namespace: utils.KubeSystemNamespace,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *VMClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vmcontrollerv1.VMClaim{}).
		Complete(r)
}
