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

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = &aws.JSONValue{}
	_ = ackv1alpha1.AWSAccountID("")
)

// Provides a description of an EFS file system access point.
type AccessPointDescription struct {
	AccessPointARN *string `json:"accessPointARN,omitempty"`
	AccessPointID  *string `json:"accessPointID,omitempty"`
	FileSystemID   *string `json:"fileSystemID,omitempty"`
	LifeCycleState *string `json:"lifeCycleState,omitempty"`
	Name           *string `json:"name,omitempty"`
	OwnerID        *string `json:"ownerID,omitempty"`
	// The full POSIX identity, including the user ID, group ID, and any secondary
	// group IDs, on the access point that is used for all file system operations
	// performed by NFS clients using the access point.
	PosixUser *PosixUser `json:"posixUser,omitempty"`
	// Specifies the directory on the Amazon EFS file system that the access point
	// provides access to. The access point exposes the specified file system path
	// as the root directory of your file system to applications using the access
	// point. NFS clients using the access point can only access data in the access
	// point's RootDirectory and it's subdirectories.
	RootDirectory *RootDirectory `json:"rootDirectory,omitempty"`
	Tags          []*Tag         `json:"tags,omitempty"`
}

// The backup policy for the file system used to create automatic daily backups.
// If status has a value of ENABLED, the file system is being automatically
// backed up. For more information, see Automatic backups (https://docs.aws.amazon.com/efs/latest/ug/awsbackup.html#automatic-backups).
type BackupPolicy struct {
	Status *string `json:"status,omitempty"`
}

// Required if the RootDirectory > Path specified does not exist. Specifies
// the POSIX IDs and permissions to apply to the access point's RootDirectory
// > Path. If the access point root directory does not exist, EFS creates it
// with these settings when a client connects to the access point. When specifying
// CreationInfo, you must include values for all properties.
//
// Amazon EFS creates a root directory only if you have provided the CreationInfo:
// OwnUid, OwnGID, and permissions for the directory. If you do not provide
// this information, Amazon EFS does not create the root directory. If the root
// directory does not exist, attempts to mount using the access point will fail.
//
// If you do not provide CreationInfo and the specified RootDirectory does not
// exist, attempts to mount the file system using the access point will fail.
type CreationInfo struct {
	OwnerGID    *int64  `json:"ownerGID,omitempty"`
	OwnerUID    *int64  `json:"ownerUID,omitempty"`
	Permissions *string `json:"permissions,omitempty"`
}

// Describes the destination file system in the replication configuration.
type Destination struct {
	FileSystemID            *string      `json:"fileSystemID,omitempty"`
	LastReplicatedTimestamp *metav1.Time `json:"lastReplicatedTimestamp,omitempty"`
}

// Describes the new or existing destination file system for the replication
// configuration.
type DestinationToCreate struct {
	AvailabilityZoneName *string `json:"availabilityZoneName,omitempty"`
	FileSystemID         *string `json:"fileSystemID,omitempty"`
	KMSKeyID             *string `json:"kmsKeyID,omitempty"`
}

// A description of the file system.
type FileSystemDescription struct {
	AvailabilityZoneID   *string      `json:"availabilityZoneID,omitempty"`
	AvailabilityZoneName *string      `json:"availabilityZoneName,omitempty"`
	CreationTime         *metav1.Time `json:"creationTime,omitempty"`
	Encrypted            *bool        `json:"encrypted,omitempty"`
	FileSystemARN        *string      `json:"fileSystemARN,omitempty"`
	FileSystemID         *string      `json:"fileSystemID,omitempty"`
	// Describes the protection on a file system.
	FileSystemProtection         *FileSystemProtectionDescription `json:"fileSystemProtection,omitempty"`
	KMSKeyID                     *string                          `json:"kmsKeyID,omitempty"`
	LifeCycleState               *string                          `json:"lifeCycleState,omitempty"`
	Name                         *string                          `json:"name,omitempty"`
	NumberOfMountTargets         *int64                           `json:"numberOfMountTargets,omitempty"`
	OwnerID                      *string                          `json:"ownerID,omitempty"`
	PerformanceMode              *string                          `json:"performanceMode,omitempty"`
	ProvisionedThroughputInMiBps *float64                         `json:"provisionedThroughputInMiBps,omitempty"`
	// The latest known metered size (in bytes) of data stored in the file system,
	// in its Value field, and the time at which that size was determined in its
	// Timestamp field. The value doesn't represent the size of a consistent snapshot
	// of the file system, but it is eventually consistent when there are no writes
	// to the file system. That is, the value represents the actual size only if
	// the file system is not modified for a period longer than a couple of hours.
	// Otherwise, the value is not necessarily the exact size the file system was
	// at any instant in time.
	SizeInBytes    *FileSystemSize `json:"sizeInBytes,omitempty"`
	Tags           []*Tag          `json:"tags,omitempty"`
	ThroughputMode *string         `json:"throughputMode,omitempty"`
}

// Describes the protection on a file system.
type FileSystemProtectionDescription struct {
	ReplicationOverwriteProtection *string `json:"replicationOverwriteProtection,omitempty"`
}

// The latest known metered size (in bytes) of data stored in the file system,
// in its Value field, and the time at which that size was determined in its
// Timestamp field. The value doesn't represent the size of a consistent snapshot
// of the file system, but it is eventually consistent when there are no writes
// to the file system. That is, the value represents the actual size only if
// the file system is not modified for a period longer than a couple of hours.
// Otherwise, the value is not necessarily the exact size the file system was
// at any instant in time.
type FileSystemSize struct {
	Timestamp       *metav1.Time `json:"timestamp,omitempty"`
	Value           *int64       `json:"value,omitempty"`
	ValueInArchive  *int64       `json:"valueInArchive,omitempty"`
	ValueInIA       *int64       `json:"valueInIA,omitempty"`
	ValueInStandard *int64       `json:"valueInStandard,omitempty"`
}

// Describes a policy used by Lifecycle management that specifies when to transition
// files into and out of storage classes. For more information, see Managing
// file system storage (https://docs.aws.amazon.com/efs/latest/ug/lifecycle-management-efs.html).
//
// When using the put-lifecycle-configuration CLI command or the PutLifecycleConfiguration
// API action, Amazon EFS requires that each LifecyclePolicy object have only
// a single transition. This means that in a request body, LifecyclePolicies
// must be structured as an array of LifecyclePolicy objects, one object for
// each transition. For more information, see the request examples in PutLifecycleConfiguration.
type LifecyclePolicy struct {
	TransitionToArchive             *string `json:"transitionToArchive,omitempty"`
	TransitionToIA                  *string `json:"transitionToIA,omitempty"`
	TransitionToPrimaryStorageClass *string `json:"transitionToPrimaryStorageClass,omitempty"`
}

// Provides a description of a mount target.
type MountTargetDescription struct {
	AvailabilityZoneID   *string `json:"availabilityZoneID,omitempty"`
	AvailabilityZoneName *string `json:"availabilityZoneName,omitempty"`
	FileSystemID         *string `json:"fileSystemID,omitempty"`
	IPAddress            *string `json:"ipAddress,omitempty"`
	LifeCycleState       *string `json:"lifeCycleState,omitempty"`
	MountTargetID        *string `json:"mountTargetID,omitempty"`
	NetworkInterfaceID   *string `json:"networkInterfaceID,omitempty"`
	OwnerID              *string `json:"ownerID,omitempty"`
	SubnetID             *string `json:"subnetID,omitempty"`
	VPCID                *string `json:"vpcID,omitempty"`
}

// The full POSIX identity, including the user ID, group ID, and any secondary
// group IDs, on the access point that is used for all file system operations
// performed by NFS clients using the access point.
type PosixUser struct {
	GID           *int64   `json:"gid,omitempty"`
	SecondaryGIDs []*int64 `json:"secondaryGIDs,omitempty"`
	UID           *int64   `json:"uid,omitempty"`
}

// Describes the replication configuration for a specific file system.
type ReplicationConfigurationDescription struct {
	CreationTime                *metav1.Time `json:"creationTime,omitempty"`
	OriginalSourceFileSystemARN *string      `json:"originalSourceFileSystemARN,omitempty"`
	SourceFileSystemARN         *string      `json:"sourceFileSystemARN,omitempty"`
	SourceFileSystemID          *string      `json:"sourceFileSystemID,omitempty"`
}

// Specifies the directory on the Amazon EFS file system that the access point
// provides access to. The access point exposes the specified file system path
// as the root directory of your file system to applications using the access
// point. NFS clients using the access point can only access data in the access
// point's RootDirectory and it's subdirectories.
type RootDirectory struct {
	// Required if the RootDirectory > Path specified does not exist. Specifies
	// the POSIX IDs and permissions to apply to the access point's RootDirectory
	// > Path. If the access point root directory does not exist, EFS creates it
	// with these settings when a client connects to the access point. When specifying
	// CreationInfo, you must include values for all properties.
	//
	// Amazon EFS creates a root directory only if you have provided the CreationInfo:
	// OwnUid, OwnGID, and permissions for the directory. If you do not provide
	// this information, Amazon EFS does not create the root directory. If the root
	// directory does not exist, attempts to mount using the access point will fail.
	//
	// If you do not provide CreationInfo and the specified RootDirectory does not
	// exist, attempts to mount the file system using the access point will fail.
	CreationInfo *CreationInfo `json:"creationInfo,omitempty"`
	Path         *string       `json:"path,omitempty"`
}

// A tag is a key-value pair. Allowed characters are letters, white space, and
// numbers that can be represented in UTF-8, and the following characters:+
// - = . _ : /.
type Tag struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type UpdateFileSystemProtectionInput struct {
	ReplicationOverwriteProtection *string `json:"replicationOverwriteProtection,omitempty"`
}
