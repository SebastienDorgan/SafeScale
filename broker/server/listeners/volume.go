/*
 * Copyright 2018, CS Systemes d'Information, http://www.c-s.fr
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

package listeners

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"

	pb "github.com/CS-SI/SafeScale/broker"
	"github.com/CS-SI/SafeScale/broker/server/handlers"
	"github.com/CS-SI/SafeScale/broker/utils"
	conv "github.com/CS-SI/SafeScale/broker/utils"
	"github.com/CS-SI/SafeScale/providers/model/enums/VolumeSpeed"
)

// broker volume create v1 --speed="SSD" --size=2000 (par default HDD, possible SSD, HDD, COLD)
// broker volume attach v1 host1 --path="/shared/data" --format="xfs" (par default /shared/v1 et ext4)
// broker volume detach v1
// broker volume delete v1
// broker volume inspect v1
// broker volume update v1 --speed="HDD" --size=1000

// VolumeHandler ...
var VolumeHandler = handlers.NewVolumeHandler

// VolumeListener is the volume service grps server
type VolumeListener struct{}

// List the available volumes
func (s *VolumeListener) List(ctx context.Context, in *pb.VolumeListRequest) (*pb.VolumeList, error) {
	log.Printf("Volume List called")

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't list volumes: no tenant set")
	}

	handler := VolumeHandler(tenant.Service)
	volumes, err := handler.List(in.GetAll())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Map model.Volume to pb.Volume
	var pbvolumes []*pb.Volume
	for _, volume := range volumes {
		pbvolumes = append(pbvolumes, conv.ToPBVolume(&volume))
	}
	rv := &pb.VolumeList{Volumes: pbvolumes}
	return rv, nil
}

// Create a new volume
func (s *VolumeListener) Create(ctx context.Context, in *pb.VolumeDefinition) (*pb.Volume, error) {
	log.Infof("Listeners: volume create '%v' called", in)
	defer log.Debugf("Listeners: volume create '%v' done", in)

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't create volume: no tenant set")
	}

	handler := VolumeHandler(tenant.Service)
	volume, err := handler.Create(in.GetName(), int(in.GetSize()), VolumeSpeed.Enum(in.GetSpeed()))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	log.Printf("Volume '%s' created: %v", in.GetName(), volume.Name)
	return conv.ToPBVolume(volume), nil
}

// Attach a volume to an host and create a mount point
func (s *VolumeListener) Attach(ctx context.Context, in *pb.VolumeAttachment) (*google_protobuf.Empty, error) {
	log.Printf("Attach volume called '%s', '%s'", in.Host.Name, in.MountPath)

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't attach volume: no tenant set")
	}

	handler := VolumeHandler(tenant.Service)
	err := handler.Attach(in.GetVolume().GetName(), in.GetHost().GetName(), in.GetMountPath(), in.GetFormat())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	return &google_protobuf.Empty{}, nil
}

// Detach a volume from an host. It umount associated mountpoint
func (s *VolumeListener) Detach(ctx context.Context, in *pb.VolumeDetachment) (*google_protobuf.Empty, error) {
	log.Debugf("broker.server.listeners.VolumeListener.Detach(%v) called", in)
	defer log.Debugf("broker.server.listeners.VolumeListener.Detach(%v) done", in)

	volumeName := in.GetVolume().GetName()
	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't detach volume: no tenant set")
	}

	hostName := in.GetHost().GetName()
	handler := VolumeHandler(tenant.Service)
	err := handler.Detach(volumeName, hostName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	log.Println(fmt.Sprintf("Volume '%s' detached from '%s'", volumeName, hostName))
	return &google_protobuf.Empty{}, nil
}

// Delete a volume
func (s *VolumeListener) Delete(ctx context.Context, in *pb.Reference) (*google_protobuf.Empty, error) {
	log.Printf("Volume delete called '%s'", in.Name)

	ref := utils.GetReference(in)
	if ref == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "can't inspect volume: neither name nor id given as reference")
	}

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't delete volume: no tenant set")
	}

	handler := VolumeHandler(tenant.Service)
	err := handler.Delete(ref)
	if err != nil {
		return &google_protobuf.Empty{}, grpc.Errorf(codes.Internal, fmt.Sprintf("can't delete volume '%s': %s", ref, err.Error()))
	}
	log.Printf("Volume '%s' successfully deleted.", ref)
	return &google_protobuf.Empty{}, nil
}

// Inspect a volume
func (s *VolumeListener) Inspect(ctx context.Context, in *pb.Reference) (*pb.VolumeInfo, error) {
	log.Printf("Inspect Volume called '%s'", in.Name)

	ref := utils.GetReference(in)
	if ref == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "can't inspect volume: neither name nor id given as reference")
	}

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't inspect volume: no tenant set")
	}

	handler := VolumeHandler(tenant.Service)
	volume, mounts, err := handler.Inspect(ref)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	if volume == nil {
		return nil, grpc.Errorf(codes.NotFound, fmt.Sprintf("can't inspect volume: no volume '%s' found", ref))
	}

	return conv.ToPBVolumeInfo(volume, mounts), nil
}