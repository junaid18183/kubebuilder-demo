/*
Copyright 2023 Juned Memon.

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

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	enbuildv1alpha1 "vivsoftorg/enbuild/api/v1alpha1"

	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=enbuild.vivsoft.io,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	application := &enbuildv1alpha1.Application{}
	if err := r.Get(ctx, req.NamespacedName, application); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "failed to get Cluster resource", "application", req.NamespacedName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	logger.Info("got application", "Owner", application.Spec.Owner)

	// Reconcile k8s gitrepository.
	gitrepository, err := generateGitRepository(application, logger, r)
	if err != nil {
		return ctrl.Result{}, err
	}
	logger.Info(gitrepository.Spec.URL)
	if err := ctrl.SetControllerReference(application, gitrepository, r.Scheme); err != nil {

		return ctrl.Result{}, err
	}
	// if err := GitRepository(ctx, r, application, logger); err != nil {
	// 	return ctrl.Result{}, err
	// }

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&enbuildv1alpha1.Application{}).
		Owns(&sourcev1.GitRepository{}).
		Complete(r)
}

func generateGitRepository(tb *enbuildv1alpha1.Application, log logr.Logger, r *ApplicationReconciler) (*sourcev1.GitRepository, error) {
	return &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tb.Name,
			Namespace: tb.Namespace,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL: "https://gitlab.com/enbuild-staging/iac-templates/bigbang",
		},
	}, nil
}

func GitRepository(ctx context.Context, r client.Client, application *enbuildv1alpha1.Application, log logr.Logger) error {
	foundApplication := &enbuildv1alpha1.Application{}
	// justCreated := false
	if err := r.Get(ctx, types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, foundApplication); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating Application", "namespace", application.Namespace, "name", application.Name)
			if err := r.Create(ctx, application); err != nil {
				log.Error(err, "unable to create application")
				return err
			}
			// justCreated = true
		} else {
			log.Error(err, "error getting application")
			return err
		}
	}
	// if !justCreated && CopyApplicationSetFields(application, foundApplication) {
	// 	log.Info("Updating Application", "namespace", application.Namespace, "name", application.Name)
	// 	if err := r.Update(ctx, foundApplication); err != nil {
	// 		log.Error(err, "unable to update application")
	// 		return err
	// 	}
	// }

	return nil
}
