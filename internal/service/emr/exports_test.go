// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package emr

// Exports for use in tests only.
var (
	ResourceBlockPublicAccessConfiguration = resourceBlockPublicAccessConfiguration
	ResourceCluster                        = resourceCluster
	ResourceInstanceFleet                  = resourceInstanceFleet
	ResourceInstanceGroup                  = resourceInstanceGroup
	ResourceManagedScalingPolicy           = resourceManagedScalingPolicy
	ResourceSecurityConfiguration          = resourceSecurityConfiguration
	ResourceStudio                         = resourceStudio
	ResourceStudioSessionMapping           = resourceStudioSessionMapping

	FindBlockPublicAccessConfiguration = findBlockPublicAccessConfiguration
	FindClusterByID                    = findClusterByID
	FindInstanceFleetByTwoPartKey      = findInstanceFleetByTwoPartKey
	FindInstanceGroupByTwoPartKey      = findInstanceGroupByTwoPartKey
	FindManagedScalingPolicyByID       = findManagedScalingPolicyByID
	FindSecurityConfigurationByName    = findSecurityConfigurationByName
	FindStudioByID                     = findStudioByID
	FindStudioSessionMappingByIDOrName = findStudioSessionMappingByIDOrName
)
