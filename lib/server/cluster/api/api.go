/*
 * Copyright 2018-2019, CS Systemes d'Information, http://www.c-s.fr
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	pb "github.com/CS-SI/SafeScale/lib"
	clusterpropsv1 "github.com/CS-SI/SafeScale/lib/server/cluster/control/properties/v1"
	"github.com/CS-SI/SafeScale/lib/server/cluster/enums/ClusterState"
	"github.com/CS-SI/SafeScale/lib/server/cluster/identity"
	"github.com/CS-SI/SafeScale/lib/server/iaas"
	"github.com/CS-SI/SafeScale/lib/utils/concurrency"
	"github.com/CS-SI/SafeScale/lib/utils/serialize"
)

//go:generate mockgen -destination=../mocks/mock_cluster.go -package=mocks github.com/CS-SI/SafeScale/lib/server/cluster/api Cluster

// Cluster is an interface of methods associated to Cluster-like structs
type Cluster interface {
	// GetService ...
	GetService(task concurrency.Task) iaas.Service
	// GetIdentity returns the identity of the cluster (name, flavor, complexity)
	GetIdentity(task concurrency.Task) identity.Identity
	// GetNetworkConfig returns network configuration of the cluster
	GetNetworkConfig(concurrency.Task) clusterpropsv1.Network
	// GetProperties returns the extension of the cluster
	GetProperties(concurrency.Task) *serialize.JSONProperties

	// Start starts the cluster
	Start(concurrency.Task) error
	// Stop stops the cluster
	Stop(concurrency.Task) error
	// GetState returns the current state of the cluster
	GetState(concurrency.Task) (ClusterState.Enum, error)
	// AddNode adds a node
	AddNode(concurrency.Task, *pb.HostDefinition) (string, error)
	// AddNodes adds several nodes
	AddNodes(concurrency.Task, int, *pb.HostDefinition) ([]string, error)
	// DeleteLastNode deletes a node
	DeleteLastNode(concurrency.Task, string) error
	// DeleteSpecificNode deletes a node identified by its ID
	DeleteSpecificNode(concurrency.Task, string, string) error
	// ListMasterIDs lists the IDs of masters (if there is such masters in the flavor...)
	ListMasterIDs(concurrency.Task) []string
	// ListMasterIPs lists the IPs of masters (if there is such masters in the flavor...)
	ListMasterIPs(concurrency.Task) []string
	// FindAvailableMaster returns ID of the first master available to execute order
	FindAvailableMaster(concurrency.Task) (string, error)
	// ListNodeIDs lists IDs of the nodes in the cluster
	ListNodeIDs(concurrency.Task) []string
	// ListNodeIPs lists the IPs of the nodes in the cluster
	ListNodeIPs(concurrency.Task) []string
	// FindAvailableNode returns ID of the first node available to execute order
	FindAvailableNode(concurrency.Task) (string, error)
	// SearchNode tells if the ID of the host passed as parameter is a node
	SearchNode(concurrency.Task, string) bool
	// GetNode returns a node based on its ID
	GetNode(concurrency.Task, string) (*pb.Host, error)
	// CountNodes counts the nodes of the cluster
	CountNodes(concurrency.Task) uint

	// Delete allows to destroy infrastructure of cluster
	Delete(concurrency.Task) error
}
