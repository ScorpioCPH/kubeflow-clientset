/*
Copyright 2017 The Caicloud KubeFlow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/package fake

import (
	clientset "github.com/caicloud/kubeflow-clientset/clientset/versioned"
	kubeflowv1alpha1 "github.com/caicloud/kubeflow-clientset/clientset/versioned/typed/kubeflow/v1alpha1"
	fakekubeflowv1alpha1 "github.com/caicloud/kubeflow-clientset/clientset/versioned/typed/kubeflow/v1alpha1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o))
	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &Clientset{fakePtr, &fakediscovery.FakeDiscovery{Fake: &fakePtr}}
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

// KubeflowV1alpha1 retrieves the KubeflowV1alpha1Client
func (c *Clientset) KubeflowV1alpha1() kubeflowv1alpha1.KubeflowV1alpha1Interface {
	return &fakekubeflowv1alpha1.FakeKubeflowV1alpha1{Fake: &c.Fake}
}

// Kubeflow retrieves the KubeflowV1alpha1Client
func (c *Clientset) Kubeflow() kubeflowv1alpha1.KubeflowV1alpha1Interface {
	return &fakekubeflowv1alpha1.FakeKubeflowV1alpha1{Fake: &c.Fake}
}