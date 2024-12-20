package controllers

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	v1 "github.com/aomaryoussef/ephemeral-operator/api/v1"
)

// EphemeralResourceReconciler reconciles a EphemeralResource object
type EphemeralResourceReconciler struct {
	client.Client
}

// +kubebuilder:rbac:groups=core.example.com,resources=ephemeralresources,verbs=get;list;create;update;patch;delete

// Reconcile reads that state of the cluster for a EphemeralResource object and makes changes based on the state read
// and what is in the EphemeralResource.Spec
func (r *EphemeralResourceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	// Fetch the EphemeralResource instance
	ephemeralResource := &v1.EphemeralResource{}
	err := r.Get(ctx, req.NamespacedName, ephemeralResource)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the TTL has expired
	if ephemeralResource.Spec.TTLSeconds < time.Now().Unix() {
		for _, res := range ephemeralResource.Spec.Resources {
			// Here, we assume the resources are Deployments or Secrets for simplicity.
			// Handle deletion logic for each resource
			if res.Kind == "Deployment" {
				deployment := &corev1.Deployment{}
				err := r.Get(ctx, client.ObjectKey{
					Namespace: req.Namespace,
					Name:      res.Name,
				}, deployment)
				if err == nil {
					err = r.Delete(ctx, deployment)
					if err != nil {
						return reconcile.Result{}, fmt.Errorf("failed to delete deployment %s: %v", res.Name, err)
					}
				}
			} else if res.Kind == "Secret" {
				secret := &corev1.Secret{}
				err := r.Get(ctx, client.ObjectKey{
					Namespace: req.Namespace,
					Name:      res.Name,
				}, secret)
				if err == nil {
					err = r.Delete(ctx, secret)
					if err != nil {
						return reconcile.Result{}, fmt.Errorf("failed to delete secret %s: %v", res.Name, err)
					}
				}
			}
		}
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EphemeralResourceReconciler) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.EphemeralResource{}).
		Owns(&corev1.Deployment{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
