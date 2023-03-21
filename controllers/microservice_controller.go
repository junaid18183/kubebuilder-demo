/*
Copyright 2023 Juned Memon.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by svclicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	enbuildv1alpha1 "vivsoftorg/enbuild/api/v1alpha1"
)

// MicroServiceReconciler reconciles a MicroService object
type MicroServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//---------------------------------------------------------------------------------------------------------------------------------------------------------------

//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=microservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=microservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=microservices/finalizers,verbs=update

// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=GitRepository,verbs=list;watch;get;patch;create;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MicroService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MicroServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here
	microservice := &enbuildv1alpha1.MicroService{}
	if err := r.Get(ctx, req.NamespacedName, microservice); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "failed to get Cluster resource", "MicroService", req.NamespacedName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// logger.Info("got microservice", "Owner", microservice.Spec.Owner)

	// Reconcile k8s gitrepository.
	ReconcileGitRepository1(ctx, r, microservice, logger)

	return ctrl.Result{}, nil
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
// SetupWithManager sets up the controller with the Manager.
func (r *MicroServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&enbuildv1alpha1.MicroService{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{}, predicate.AnnotationChangedPredicate{})).
		Owns(&sourcev1.GitRepository{}).
		Complete(r)
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
func CreateGitRepository1(ctx context.Context, r *MicroServiceReconciler, gitrepository *sourcev1.GitRepository, log logr.Logger) error {
	foundGitRepository := &sourcev1.GitRepository{}
	if err := r.Get(ctx, types.NamespacedName{Name: gitrepository.Name, Namespace: gitrepository.Namespace}, foundGitRepository); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating GitRepository", "namespace", gitrepository.Namespace, "name", gitrepository.Name)
			if err := r.Create(ctx, gitrepository); err != nil {
				log.Error(err, "unable to create gitrepository")
				return err
			}
		} else {
			log.Error(err, "error getting gitrepository")
			return err
		}
	}

	return nil
}

//---------------------------------------------------------------------------------------------------------------------------------------------------------------

func ReconcileGitRepository1(ctx context.Context, r *MicroServiceReconciler, microservice *enbuildv1alpha1.MicroService, logger logr.Logger) (ctrl.Result, error) {
	gitrepository, err := generateGitRepositorySpec(microservice.Name, microservice.Namespace, "https://gitlab.com/enbuild-staging/iac-templates/bigbang", microservice.Spec.SecretRef.Name, "main")

	if err != nil {
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(microservice, gitrepository, r.Scheme); err != nil {

		return ctrl.Result{}, err
	}
	if err := CreateGitRepository1(ctx, r, gitrepository, logger); err != nil {
		return ctrl.Result{}, err
	}

	microservice.Status.Repository = gitrepository.Spec.URL

	_err := r.Status().Update(ctx, microservice)
	if _err != nil {
		return ctrl.Result{}, _err
	}
	return ctrl.Result{}, _err
}

//---------------------------------------------------------------------------------------------------------------------------------------------------------------
