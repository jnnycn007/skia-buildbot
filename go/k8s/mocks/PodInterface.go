// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mock "github.com/stretchr/testify/mock"

	policyv1 "k8s.io/api/policy/v1"

	rest "k8s.io/client-go/rest"

	types "k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/client-go/applyconfigurations/core/v1"

	v1beta1 "k8s.io/api/policy/v1beta1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// PodInterface is an autogenerated mock type for the PodInterface type
type PodInterface struct {
	mock.Mock
}

// Apply provides a mock function with given fields: ctx, pod, opts
func (_m *PodInterface) Apply(ctx context.Context, pod *v1.PodApplyConfiguration, opts metav1.ApplyOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for Apply")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.PodApplyConfiguration, metav1.ApplyOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.PodApplyConfiguration, metav1.ApplyOptions) *corev1.Pod); ok {
		r0 = rf(ctx, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.PodApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ApplyStatus provides a mock function with given fields: ctx, pod, opts
func (_m *PodInterface) ApplyStatus(ctx context.Context, pod *v1.PodApplyConfiguration, opts metav1.ApplyOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for ApplyStatus")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.PodApplyConfiguration, metav1.ApplyOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.PodApplyConfiguration, metav1.ApplyOptions) *corev1.Pod); ok {
		r0 = rf(ctx, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.PodApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Bind provides a mock function with given fields: ctx, binding, opts
func (_m *PodInterface) Bind(ctx context.Context, binding *corev1.Binding, opts metav1.CreateOptions) error {
	ret := _m.Called(ctx, binding, opts)

	if len(ret) == 0 {
		panic("no return value specified for Bind")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Binding, metav1.CreateOptions) error); ok {
		r0 = rf(ctx, binding, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, pod, opts
func (_m *PodInterface) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Pod, metav1.CreateOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Pod, metav1.CreateOptions) *corev1.Pod); ok {
		r0 = rf(ctx, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *corev1.Pod, metav1.CreateOptions) error); ok {
		r1 = rf(ctx, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *PodInterface) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *PodInterface) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCollection")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Evict provides a mock function with given fields: ctx, eviction
func (_m *PodInterface) Evict(ctx context.Context, eviction *v1beta1.Eviction) error {
	ret := _m.Called(ctx, eviction)

	if len(ret) == 0 {
		panic("no return value specified for Evict")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1beta1.Eviction) error); ok {
		r0 = rf(ctx, eviction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EvictV1 provides a mock function with given fields: ctx, eviction
func (_m *PodInterface) EvictV1(ctx context.Context, eviction *policyv1.Eviction) error {
	ret := _m.Called(ctx, eviction)

	if len(ret) == 0 {
		panic("no return value specified for EvictV1")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *policyv1.Eviction) error); ok {
		r0 = rf(ctx, eviction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EvictV1beta1 provides a mock function with given fields: ctx, eviction
func (_m *PodInterface) EvictV1beta1(ctx context.Context, eviction *v1beta1.Eviction) error {
	ret := _m.Called(ctx, eviction)

	if len(ret) == 0 {
		panic("no return value specified for EvictV1beta1")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1beta1.Eviction) error); ok {
		r0 = rf(ctx, eviction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *PodInterface) Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *corev1.Pod); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLogs provides a mock function with given fields: name, opts
func (_m *PodInterface) GetLogs(name string, opts *corev1.PodLogOptions) *rest.Request {
	ret := _m.Called(name, opts)

	if len(ret) == 0 {
		panic("no return value specified for GetLogs")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func(string, *corev1.PodLogOptions) *rest.Request); ok {
		r0 = rf(name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// List provides a mock function with given fields: ctx, opts
func (_m *PodInterface) List(ctx context.Context, opts metav1.ListOptions) (*corev1.PodList, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *corev1.PodList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (*corev1.PodList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *corev1.PodList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.PodList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *PodInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*corev1.Pod, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*corev1.Pod, error)); ok {
		return rf(ctx, name, pt, data, opts, subresources...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) *corev1.Pod); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProxyGet provides a mock function with given fields: scheme, name, port, path, params
func (_m *PodInterface) ProxyGet(scheme string, name string, port string, path string, params map[string]string) rest.ResponseWrapper {
	ret := _m.Called(scheme, name, port, path, params)

	if len(ret) == 0 {
		panic("no return value specified for ProxyGet")
	}

	var r0 rest.ResponseWrapper
	if rf, ok := ret.Get(0).(func(string, string, string, string, map[string]string) rest.ResponseWrapper); ok {
		r0 = rf(scheme, name, port, path, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rest.ResponseWrapper)
		}
	}

	return r0
}

// Update provides a mock function with given fields: ctx, pod, opts
func (_m *PodInterface) Update(ctx context.Context, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Pod, metav1.UpdateOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Pod, metav1.UpdateOptions) *corev1.Pod); ok {
		r0 = rf(ctx, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *corev1.Pod, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEphemeralContainers provides a mock function with given fields: ctx, podName, pod, opts
func (_m *PodInterface) UpdateEphemeralContainers(ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, podName, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEphemeralContainers")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *corev1.Pod, metav1.UpdateOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, podName, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *corev1.Pod, metav1.UpdateOptions) *corev1.Pod); ok {
		r0 = rf(ctx, podName, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *corev1.Pod, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, podName, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateResize provides a mock function with given fields: ctx, podName, pod, opts
func (_m *PodInterface) UpdateResize(ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, podName, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateResize")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *corev1.Pod, metav1.UpdateOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, podName, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *corev1.Pod, metav1.UpdateOptions) *corev1.Pod); ok {
		r0 = rf(ctx, podName, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *corev1.Pod, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, podName, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: ctx, pod, opts
func (_m *PodInterface) UpdateStatus(ctx context.Context, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	ret := _m.Called(ctx, pod, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 *corev1.Pod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Pod, metav1.UpdateOptions) (*corev1.Pod, error)); ok {
		return rf(ctx, pod, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *corev1.Pod, metav1.UpdateOptions) *corev1.Pod); ok {
		r0 = rf(ctx, pod, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.Pod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *corev1.Pod, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, pod, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *PodInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Watch")
	}

	var r0 watch.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (watch.Interface, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPodInterface creates a new instance of PodInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPodInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *PodInterface {
	mock := &PodInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
