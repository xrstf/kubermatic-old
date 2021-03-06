package kubernetes

import (
	"sync"

	"github.com/kubermatic/kubermatic/api/pkg/crd/client/clientset/versioned/scheme"

	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

const componentName = "kubermatic-api"

// NewEventRecorder returns a new event recorder provider object. See EventRecorder for more information.
func NewEventRecorder() *EventRecorder {
	return &EventRecorder{
		seedClusterRecorderMap: make(map[string]record.EventRecorder),
	}
}

// EventRecorder gives option to record events for objects. They can be then read from them using K8S API.
type EventRecorder struct {
	seedClusterRecorderMap map[string]record.EventRecorder
	lock                   sync.Mutex
}

// ClusterRecorderFor returns an event recorder that will be able to record events for objects in the cluster
// accessible using provided client.
func (e *EventRecorder) ClusterRecorderFor(client kubernetes.Interface) record.EventRecorder {
	return e.getRecorderForClient(client)
}

func (e *EventRecorder) getRecorderForClient(client kubernetes.Interface) record.EventRecorder {
	e.lock.Lock()
	defer e.lock.Unlock()

	coreV1Client := client.CoreV1()
	host := coreV1Client.RESTClient().Get().URL().Host
	recorder, exists := e.seedClusterRecorderMap[host]
	if exists {
		return recorder
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&corev1.EventSinkImpl{Interface: coreV1Client.Events("")})
	recorder = eventBroadcaster.NewRecorder(scheme.Scheme, apicorev1.EventSource{Component: componentName})
	e.seedClusterRecorderMap[host] = recorder

	return recorder
}
