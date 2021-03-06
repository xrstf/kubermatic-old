// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	kubermaticv1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAddons implements AddonInterface
type FakeAddons struct {
	Fake *FakeKubermaticV1
	ns   string
}

var addonsResource = schema.GroupVersionResource{Group: "kubermatic.k8s.io", Version: "v1", Resource: "addons"}

var addonsKind = schema.GroupVersionKind{Group: "kubermatic.k8s.io", Version: "v1", Kind: "Addon"}

// Get takes name of the addon, and returns the corresponding addon object, and an error if there is any.
func (c *FakeAddons) Get(name string, options v1.GetOptions) (result *kubermaticv1.Addon, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(addonsResource, c.ns, name), &kubermaticv1.Addon{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubermaticv1.Addon), err
}

// List takes label and field selectors, and returns the list of Addons that match those selectors.
func (c *FakeAddons) List(opts v1.ListOptions) (result *kubermaticv1.AddonList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(addonsResource, addonsKind, c.ns, opts), &kubermaticv1.AddonList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &kubermaticv1.AddonList{ListMeta: obj.(*kubermaticv1.AddonList).ListMeta}
	for _, item := range obj.(*kubermaticv1.AddonList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested addons.
func (c *FakeAddons) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(addonsResource, c.ns, opts))

}

// Create takes the representation of a addon and creates it.  Returns the server's representation of the addon, and an error, if there is any.
func (c *FakeAddons) Create(addon *kubermaticv1.Addon) (result *kubermaticv1.Addon, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(addonsResource, c.ns, addon), &kubermaticv1.Addon{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubermaticv1.Addon), err
}

// Update takes the representation of a addon and updates it. Returns the server's representation of the addon, and an error, if there is any.
func (c *FakeAddons) Update(addon *kubermaticv1.Addon) (result *kubermaticv1.Addon, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(addonsResource, c.ns, addon), &kubermaticv1.Addon{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubermaticv1.Addon), err
}

// Delete takes name of the addon and deletes it. Returns an error if one occurs.
func (c *FakeAddons) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(addonsResource, c.ns, name), &kubermaticv1.Addon{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAddons) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(addonsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &kubermaticv1.AddonList{})
	return err
}

// Patch applies the patch and returns the patched addon.
func (c *FakeAddons) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kubermaticv1.Addon, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(addonsResource, c.ns, name, pt, data, subresources...), &kubermaticv1.Addon{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubermaticv1.Addon), err
}
