package validation

import (
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/util/validation/field"

	"fmt"
	oapi "github.com/openshift/origin/pkg/api"
	applicationapi "github.com/openshift/origin/pkg/application/api"
	applicationutil "github.com/openshift/origin/pkg/application/util"
	oclient "github.com/openshift/origin/pkg/client"
	kerrors "k8s.io/kubernetes/pkg/api/errors"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"strings"
)

const MinApplicationLength = 32

func ValidateApplicationName(name string, prefix bool) []string {
	if reasons := oapi.MinimalNameRequirements(name, prefix); len(reasons) != 0 {
		return reasons
	}

	if len(name) < MinApplicationLength {
		return []string{fmt.Sprintf("must be at least %d characters long", MinApplicationLength)}
	}
	return nil
}

func ValidationApplicationItemKind(items applicationapi.ItemList) []string {
	for _, item := range items {
		if !applicationutil.Contains(applicationapi.ApplicationItemSupportKinds, item.Kind) {
			return []string{fmt.Sprintf("item unsupport selected kind %s", item.Kind)}
		}

		if len(item.Name) < 2 {
			return []string{"item name must be at least 2 characters long"}
		}

		if reasons := oapi.MinimalNameRequirements(item.Name, false); len(reasons) != 0 {
			return reasons
		}
	}
	return nil
}

func ValidationApplicationItemName(namespace string, items applicationapi.ItemList, oClient *oclient.Client, kClient *kclient.Client) (bool, string) {
	for _, item := range items {
		switch item.Kind {
		case "ServiceBroker":
			if _, err := oClient.ServiceBrokers().Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}
		case "BackingServiceInstance":
			if _, err := oClient.BackingServiceInstances(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "Build":
			if _, err := oClient.Builds(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "BuildConfig":
			if _, err := oClient.BuildConfigs(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "DeploymentConfig":
			if _, err := oClient.DeploymentConfigs(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "ImageStream":
			if _, err := oClient.ImageStreams(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "ReplicationController":
			if _, err := kClient.ReplicationControllers(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "Node":
			if _, err := kClient.Nodes().Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "Pod":
			if _, err := kClient.Pods(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}

		case "Service":
			if _, err := kClient.Services(namespace).Get(item.Name); err != nil {
				if kerrors.IsNotFound(err) {
					return false, fmt.Sprintf("resource %s=%s no found.", item.Kind, item.Name)
				}
			}
		}
	}
	return true, ""
}

// ValidateApplication tests required fields for a Application.
// This should only be called when creating a application (not on update),
// since its name validation is more restrictive than default namespace name validation
func ValidateApplication(application *applicationapi.Application, oClient *oclient.Client, kClient *kclient.Client) field.ErrorList {
	result := validation.ValidateObjectMeta(&application.ObjectMeta, true, oapi.MinimalNameRequirements, field.NewPath("metadata"))

	if reasons := ValidationApplicationItemKind(application.Spec.Items); len(reasons) != 0 {
		result = append(result, field.Invalid(field.NewPath("items"), application.Spec.Items, strings.Join(reasons, ", ")))
	}

	if ok, err := ValidationApplicationItemName(application.Namespace, application.Spec.Items, oClient, kClient); !ok {
		result = append(result, field.Invalid(field.NewPath("items"), application.Spec.Items, err))
	}

	return result
}

// ValidateApplication tests required fields for a Application.
// This should only be called when creating a application (not on update),
// since its name validation is more restrictive than default namespace name validation
func ValidateApplicationProxy(application *applicationapi.Application) field.ErrorList {
	result := validation.ValidateObjectMeta(&application.ObjectMeta, true, oapi.MinimalNameRequirements, field.NewPath("metadata"))
	return result
}

// ValidateApplicationUpdate tests to make sure a application update can be applied.  Modifies newApplication with immutable fields.
func ValidateApplicationUpdate(newApplication *applicationapi.Application, oldApplication *applicationapi.Application) field.ErrorList {
	allErrs := validation.ValidateObjectMetaUpdate(&newApplication.ObjectMeta, &oldApplication.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateApplicationProxy(newApplication)...)

	return allErrs
}
