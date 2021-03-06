package api

import "k8s.io/kubernetes/pkg/fields"

// BackingServiceInstanceToSelectableFields returns a label set that represents the object
func BackingServiceInstanceToSelectableFields(backingServiceInstance *BackingServiceInstance) fields.Set {
	return fields.Set{
		"metadata.name":      backingServiceInstance.Name,
		"metadata.namespace": backingServiceInstance.Namespace,
		
		"spec.provisioning.backingservice_name": backingServiceInstance.Spec.BackingServiceName,
	}
}
