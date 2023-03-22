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

	enbuildv1alpha1 "vivsoftorg/enbuild/api/v1alpha1"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"github.com/fluxcd/pkg/apis/meta"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
func generateGitRepositorySpec(name string, namespace string, url string, secret_ref string, branch string) (*sourcev1.GitRepository, error) {

	return &sourcev1.GitRepository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sourcev1.GitRepositorySpec{
			URL:       url,
			SecretRef: &meta.LocalObjectReference{Name: secret_ref},
			Reference: &sourcev1.GitRepositoryRef{Branch: branch},
		},
	}, nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
func CreateGitRepositoryCR(ctx context.Context, r client.Client, gitrepository *sourcev1.GitRepository, log logr.Logger) error {
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

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------

func generateMicroServiceSpec(name string, namespace string, template enbuildv1alpha1.TemplateSpec, secret_ref string, owner string) (*enbuildv1alpha1.MicroService, error) {

	return &enbuildv1alpha1.MicroService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: enbuildv1alpha1.MicroServiceSpec{
			Owner:     owner,
			Template:  template,
			SecretRef: &meta.LocalObjectReference{Name: secret_ref},
		},
	}, nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
func CreateMicroService(ctx context.Context, r client.Client, microservice *enbuildv1alpha1.MicroService, log logr.Logger) error {
	foundMicroService := &enbuildv1alpha1.MicroService{}
	if err := r.Get(ctx, types.NamespacedName{Name: microservice.Name, Namespace: microservice.Namespace}, foundMicroService); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Creating MicroService", "namespace", microservice.Namespace, "name", microservice.Name)
			if err := r.Create(ctx, microservice); err != nil {
				log.Error(err, "unable to create microservice")
				return err
			}
		} else {
			log.Error(err, "error getting microservice")
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
func CreateGitRepository(ctx context.Context, name string, secret_ref string, visibility string, owner string, templateOwner string, templateRepo string, log logr.Logger) *github.Repository {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: secret_ref},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	private := false

	templateRepoReq := &github.TemplateRepoRequest{
		Name:               &name,
		IncludeAllBranches: &private,
		Owner:              &owner,
		Private:            &private,
	}
	repo, _, err := client.Repositories.CreateFromTemplate(ctx, templateOwner, templateRepo, templateRepoReq)
	if err != nil {
		log.Error(err, "unable to create Repository in Github")
		return nil
	}
	return repo
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------
