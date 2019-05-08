/*
Copyright 2019 Rancher Labs.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/apps/v1"
	clientset "k8s.io/client-go/kubernetes/typed/apps/v1"
	listers "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

type DeploymentHandler func(string, *v1.Deployment) (*v1.Deployment, error)

type DeploymentController interface {
	DeploymentClient

	OnChange(ctx context.Context, name string, sync DeploymentHandler)
	OnRemove(ctx context.Context, name string, sync DeploymentHandler)
	Enqueue(namespace, name string)

	Cache() DeploymentCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type DeploymentClient interface {
	Create(*v1.Deployment) (*v1.Deployment, error)
	Update(*v1.Deployment) (*v1.Deployment, error)
	UpdateStatus(*v1.Deployment) (*v1.Deployment, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Deployment, error)
	List(namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Deployment, err error)
}

type DeploymentCache interface {
	Get(namespace, name string) (*v1.Deployment, error)
	List(namespace string, selector labels.Selector) ([]*v1.Deployment, error)

	AddIndexer(indexName string, indexer DeploymentIndexer)
	GetByIndex(indexName, key string) ([]*v1.Deployment, error)
}

type DeploymentIndexer func(obj *v1.Deployment) ([]string, error)

type deploymentController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.DeploymentsGetter
	informer          informers.DeploymentInformer
	gvk               schema.GroupVersionKind
}

func NewDeploymentController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.DeploymentsGetter, informer informers.DeploymentInformer) DeploymentController {
	return &deploymentController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromDeploymentHandlerToHandler(sync DeploymentHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Deployment
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Deployment))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *deploymentController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Deployment))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateDeploymentOnChange(updater generic.Updater, handler DeploymentHandler) DeploymentHandler {
	return func(key string, obj *v1.Deployment) (*v1.Deployment, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.Deployment)
			}
		}

		return copyObj, err
	}
}

func (c *deploymentController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *deploymentController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *deploymentController) OnChange(ctx context.Context, name string, sync DeploymentHandler) {
	c.AddGenericHandler(ctx, name, FromDeploymentHandlerToHandler(sync))
}

func (c *deploymentController) OnRemove(ctx context.Context, name string, sync DeploymentHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromDeploymentHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *deploymentController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *deploymentController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *deploymentController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *deploymentController) Cache() DeploymentCache {
	return &deploymentCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *deploymentController) Create(obj *v1.Deployment) (*v1.Deployment, error) {
	return c.clientGetter.Deployments(obj.Namespace).Create(obj)
}

func (c *deploymentController) Update(obj *v1.Deployment) (*v1.Deployment, error) {
	return c.clientGetter.Deployments(obj.Namespace).Update(obj)
}

func (c *deploymentController) UpdateStatus(obj *v1.Deployment) (*v1.Deployment, error) {
	return c.clientGetter.Deployments(obj.Namespace).UpdateStatus(obj)
}

func (c *deploymentController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Deployments(namespace).Delete(name, options)
}

func (c *deploymentController) Get(namespace, name string, options metav1.GetOptions) (*v1.Deployment, error) {
	return c.clientGetter.Deployments(namespace).Get(name, options)
}

func (c *deploymentController) List(namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error) {
	return c.clientGetter.Deployments(namespace).List(opts)
}

func (c *deploymentController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Deployments(namespace).Watch(opts)
}

func (c *deploymentController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Deployment, err error) {
	return c.clientGetter.Deployments(namespace).Patch(name, pt, data, subresources...)
}

type deploymentCache struct {
	lister  listers.DeploymentLister
	indexer cache.Indexer
}

func (c *deploymentCache) Get(namespace, name string) (*v1.Deployment, error) {
	return c.lister.Deployments(namespace).Get(name)
}

func (c *deploymentCache) List(namespace string, selector labels.Selector) ([]*v1.Deployment, error) {
	return c.lister.Deployments(namespace).List(selector)
}

func (c *deploymentCache) AddIndexer(indexName string, indexer DeploymentIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Deployment))
		},
	}))
}

func (c *deploymentCache) GetByIndex(indexName, key string) (result []*v1.Deployment, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.Deployment))
	}
	return result, nil
}
