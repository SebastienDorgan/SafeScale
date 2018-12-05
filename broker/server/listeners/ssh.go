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

	"github.com/CS-SI/SafeScale/system"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/CS-SI/SafeScale/broker"
	"github.com/CS-SI/SafeScale/broker/server/handlers"
)

// SSHHandler exists to ease integration tests
var SSHHandler = handlers.NewSSHHandler

// broker ssh connect host2
// broker ssh run host2 -c "uname -a"
// broker ssh copy /file/test.txt host1://tmp
// broker ssh copy host1:/file/test.txt /tmp

// SSHListener SSH service server grpc
type SSHListener struct{}

// Run executes an ssh command an an host
func (s *SSHListener) Run(ctx context.Context, in *pb.SshCommand) (*pb.SshResponse, error) {
	log.Infof("Listeners: ssh run '%s' -c '%s'", in.Host, in.Command)

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't execute ssh command: no tenant set")
	}

	handler := SSHHandler(tenant.Service)
	retcode, stdout, stderr, err := handler.Run(in.GetHost().GetName(), in.GetCommand())
	if err != nil {
		err = grpc.Errorf(codes.Internal, err.Error())
	}
	return &pb.SshResponse{
		Status:    int32(retcode),
		OutputStd: stdout,
		OutputErr: stderr,
	}, err
}

// Copy copy file from/to an host
func (s *SSHListener) Copy(ctx context.Context, in *pb.SshCopyCommand) (*pb.SshResponse, error) {
	log.Printf("Ssh copy called '%s', '%s'", in.Source, in.Destination)

	tenant := GetCurrentTenant()
	if tenant == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "can't copy by ssh: no tenant set")
	}

	handler := SSHHandler(tenant.Service)
	retcode, stdout, stderr, err := handler.Copy(in.GetSource(), in.GetDestination())
	if err != nil {
		return nil, err
	}
	if retcode != 0 {
		return nil, fmt.Errorf("Can't copy by ssh: copy failed: retcode=%d (=%s): %s", retcode, system.SCPErrorString(retcode), stderr)
	}
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	return &pb.SshResponse{
		Status:    int32(retcode),
		OutputStd: stdout,
		OutputErr: stderr,
	}, nil
}