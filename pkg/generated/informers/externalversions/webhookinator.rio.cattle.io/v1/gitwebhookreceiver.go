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
	time "time"

	webhookinatorriocattleiov1 "github.com/rancher/rio/pkg/apis/webhookinator.rio.cattle.io/v1"
	versioned "github.com/rancher/rio/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/rancher/rio/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/rancher/rio/pkg/generated/listers/webhookinator.rio.cattle.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GitWebHookReceiverInformer provides access to a shared informer and lister for
// GitWebHookReceivers.
type GitWebHookReceiverInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.GitWebHookReceiverLister
}

type gitWebHookReceiverInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGitWebHookReceiverInformer constructs a new informer for GitWebHookReceiver type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGitWebHookReceiverInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGitWebHookReceiverInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGitWebHookReceiverInformer constructs a new informer for GitWebHookReceiver type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGitWebHookReceiverInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WebhookinatorV1().GitWebHookReceivers(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WebhookinatorV1().GitWebHookReceivers(namespace).Watch(options)
			},
		},
		&webhookinatorriocattleiov1.GitWebHookReceiver{},
		resyncPeriod,
		indexers,
	)
}

func (f *gitWebHookReceiverInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGitWebHookReceiverInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gitWebHookReceiverInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&webhookinatorriocattleiov1.GitWebHookReceiver{}, f.defaultInformer)
}

func (f *gitWebHookReceiverInformer) Lister() v1.GitWebHookReceiverLister {
	return v1.NewGitWebHookReceiverLister(f.Informer().GetIndexer())
}
