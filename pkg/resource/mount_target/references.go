// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package mount_target

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ec2apitypes "github.com/aws-controllers-k8s/ec2-controller/apis/v1alpha1"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"

	svcapitypes "github.com/aws-controllers-k8s/efs-controller/apis/v1alpha1"
)

// +kubebuilder:rbac:groups=ec2.services.k8s.aws,resources=securitygroups,verbs=get;list
// +kubebuilder:rbac:groups=ec2.services.k8s.aws,resources=securitygroups/status,verbs=get;list

// +kubebuilder:rbac:groups=ec2.services.k8s.aws,resources=subnets,verbs=get;list
// +kubebuilder:rbac:groups=ec2.services.k8s.aws,resources=subnets/status,verbs=get;list

// ClearResolvedReferences removes any reference values that were made
// concrete in the spec. It returns a copy of the input AWSResource which
// contains the original *Ref values, but none of their respective concrete
// values.
func (rm *resourceManager) ClearResolvedReferences(res acktypes.AWSResource) acktypes.AWSResource {
	ko := rm.concreteResource(res).ko.DeepCopy()

	if ko.Spec.FileSystemRef != nil {
		ko.Spec.FileSystemID = nil
	}

	if len(ko.Spec.SecurityGroupRefs) > 0 {
		ko.Spec.SecurityGroups = nil
	}

	if ko.Spec.SubnetRef != nil {
		ko.Spec.SubnetID = nil
	}

	return &resource{ko}
}

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve those
// reference field(s) into their respective target field(s). It returns a
// copy of the input AWSResource with resolved reference(s), a boolean which
// is set to true if the resource contains any references (regardless of if
// they are resolved successfully) and an error if the passed AWSResource's
// reference field(s) could not be resolved.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (acktypes.AWSResource, bool, error) {
	ko := rm.concreteResource(res).ko

	resourceHasReferences := false
	err := validateReferenceFields(ko)
	if fieldHasReferences, err := rm.resolveReferenceForFileSystemID(ctx, apiReader, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	if fieldHasReferences, err := rm.resolveReferenceForSecurityGroups(ctx, apiReader, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	if fieldHasReferences, err := rm.resolveReferenceForSubnetID(ctx, apiReader, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	return &resource{ko}, resourceHasReferences, err
}

// validateReferenceFields validates the reference field and corresponding
// identifier field.
func validateReferenceFields(ko *svcapitypes.MountTarget) error {

	if ko.Spec.FileSystemRef != nil && ko.Spec.FileSystemID != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("FileSystemID", "FileSystemRef")
	}
	if ko.Spec.FileSystemRef == nil && ko.Spec.FileSystemID == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("FileSystemID", "FileSystemRef")
	}

	if len(ko.Spec.SecurityGroupRefs) > 0 && len(ko.Spec.SecurityGroups) > 0 {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SecurityGroups", "SecurityGroupRefs")
	}

	if ko.Spec.SubnetRef != nil && ko.Spec.SubnetID != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SubnetID", "SubnetRef")
	}
	if ko.Spec.SubnetRef == nil && ko.Spec.SubnetID == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("SubnetID", "SubnetRef")
	}
	return nil
}

// resolveReferenceForFileSystemID reads the resource referenced
// from FileSystemRef field and sets the FileSystemID
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForFileSystemID(
	ctx context.Context,
	apiReader client.Reader,
	ko *svcapitypes.MountTarget,
) (hasReferences bool, err error) {
	if ko.Spec.FileSystemRef != nil && ko.Spec.FileSystemRef.From != nil {
		hasReferences = true
		arr := ko.Spec.FileSystemRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: FileSystemRef")
		}
		namespace := ko.ObjectMeta.GetNamespace()
		if arr.Namespace != nil && *arr.Namespace != "" {
			namespace = *arr.Namespace
		}
		obj := &svcapitypes.FileSystem{}
		if err := getReferencedResourceState_FileSystem(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.FileSystemID = (*string)(obj.Status.FileSystemID)
	}

	return hasReferences, nil
}

// getReferencedResourceState_FileSystem looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_FileSystem(
	ctx context.Context,
	apiReader client.Reader,
	obj *svcapitypes.FileSystem,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"FileSystem",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"FileSystem",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"FileSystem",
			namespace, name)
	}
	if obj.Status.FileSystemID == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"FileSystem",
			namespace, name,
			"Status.FileSystemID")
	}
	return nil
}

// resolveReferenceForSecurityGroups reads the resource referenced
// from SecurityGroupRefs field and sets the SecurityGroups
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForSecurityGroups(
	ctx context.Context,
	apiReader client.Reader,
	ko *svcapitypes.MountTarget,
) (hasReferences bool, err error) {
	for _, f0iter := range ko.Spec.SecurityGroupRefs {
		if f0iter != nil && f0iter.From != nil {
			hasReferences = true
			arr := f0iter.From
			if arr.Name == nil || *arr.Name == "" {
				return hasReferences, fmt.Errorf("provided resource reference is nil or empty: SecurityGroupRefs")
			}
			namespace := ko.ObjectMeta.GetNamespace()
			if arr.Namespace != nil && *arr.Namespace != "" {
				namespace = *arr.Namespace
			}
			obj := &ec2apitypes.SecurityGroup{}
			if err := getReferencedResourceState_SecurityGroup(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
				return hasReferences, err
			}
			if ko.Spec.SecurityGroups == nil {
				ko.Spec.SecurityGroups = make([]*string, 0, 1)
			}
			ko.Spec.SecurityGroups = append(ko.Spec.SecurityGroups, (*string)(obj.Status.ID))
		}
	}

	return hasReferences, nil
}

// getReferencedResourceState_SecurityGroup looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_SecurityGroup(
	ctx context.Context,
	apiReader client.Reader,
	obj *ec2apitypes.SecurityGroup,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"SecurityGroup",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"SecurityGroup",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"SecurityGroup",
			namespace, name)
	}
	if obj.Status.ID == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"SecurityGroup",
			namespace, name,
			"Status.ID")
	}
	return nil
}

// resolveReferenceForSubnetID reads the resource referenced
// from SubnetRef field and sets the SubnetID
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForSubnetID(
	ctx context.Context,
	apiReader client.Reader,
	ko *svcapitypes.MountTarget,
) (hasReferences bool, err error) {
	if ko.Spec.SubnetRef != nil && ko.Spec.SubnetRef.From != nil {
		hasReferences = true
		arr := ko.Spec.SubnetRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: SubnetRef")
		}
		namespace := ko.ObjectMeta.GetNamespace()
		if arr.Namespace != nil && *arr.Namespace != "" {
			namespace = *arr.Namespace
		}
		obj := &ec2apitypes.Subnet{}
		if err := getReferencedResourceState_Subnet(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.SubnetID = (*string)(obj.Status.SubnetID)
	}

	return hasReferences, nil
}

// getReferencedResourceState_Subnet looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_Subnet(
	ctx context.Context,
	apiReader client.Reader,
	obj *ec2apitypes.Subnet,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"Subnet",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"Subnet",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"Subnet",
			namespace, name)
	}
	if obj.Status.SubnetID == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"Subnet",
			namespace, name,
			"Status.SubnetID")
	}
	return nil
}
