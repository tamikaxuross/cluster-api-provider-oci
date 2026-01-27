/*
Copyright (c) 2021, 2022 Oracle and/or its affiliates.

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

package util

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1beta2 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// ConvertClusterV1Beta2ToV1Beta1 converts a v1beta2 Cluster to v1beta1
func ConvertClusterV1Beta2ToV1Beta1(cluster *clusterv1beta2.Cluster) (*clusterv1.Cluster, error) {
	if cluster == nil {
		return nil, nil
	}

	fmt.Printf("Converting cluster: %s/%s\n", cluster.Namespace, cluster.Name)

	// Save and nil out failureDomains to avoid unmarshal error
	originalFailureDomains := cluster.Status.FailureDomains
	cluster.Status.FailureDomains = nil

	clusterV1Beta1 := &clusterv1.Cluster{}
	data, err := json.Marshal(cluster)
	if err != nil {
		cluster.Status.FailureDomains = originalFailureDomains
		return nil, errors.Wrap(err, "failed to marshal v1beta2 cluster")
	}

	if err := json.Unmarshal(data, clusterV1Beta1); err != nil {
		cluster.Status.FailureDomains = originalFailureDomains
		return nil, errors.Wrap(err, "failed to unmarshal to v1beta1 cluster")
	}

	// Restore original
	cluster.Status.FailureDomains = originalFailureDomains

	// Convert failureDomains from array to map
	if len(originalFailureDomains) > 0 {
		clusterV1Beta1.Status.FailureDomains = make(clusterv1.FailureDomains)
		for _, fd := range originalFailureDomains {
			// Handle pointer to bool conversion
			controlPlane := false
			if fd.ControlPlane != nil {
				controlPlane = *fd.ControlPlane
			}
			clusterV1Beta1.Status.FailureDomains[fd.Name] = clusterv1.FailureDomainSpec{
				ControlPlane: controlPlane,
				Attributes:   fd.Attributes,
			}
		}
		fmt.Printf("Converted %d failure domains to map\n", len(originalFailureDomains))
	}

	return clusterV1Beta1, nil
}

// ConvertMachineV1Beta2ToV1Beta1 converts a v1beta2 Machine to v1beta1
func ConvertMachineV1Beta2ToV1Beta1(machine *clusterv1beta2.Machine) (*clusterv1.Machine, error) {
	if machine == nil {
		return nil, nil
	}

	machineV1Beta1 := &clusterv1.Machine{}
	data, err := json.Marshal(machine)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal v1beta2 machine")
	}

	if err := json.Unmarshal(data, machineV1Beta1); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal to v1beta1 machine")
	}

	return machineV1Beta1, nil
}

// ConvertMachinePoolV1Beta2ToV1Beta1 converts a v1beta2 MachinePool to v1beta1
func ConvertMachinePoolV1Beta2ToV1Beta1(machinePool *clusterv1beta2.MachinePool) (*clusterv1.MachinePool, error) {
	if machinePool == nil {
		return nil, nil
	}

	machinePoolV1Beta1 := &clusterv1.MachinePool{}
	data, err := json.Marshal(machinePool)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal v1beta2 machinepool")
	}

	if err := json.Unmarshal(data, machinePoolV1Beta1); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal to v1beta1 machinepool")
	}

	return machinePoolV1Beta1, nil
}
