// Copyright 2013-2017 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	. "github.com/aerospike/aerospike-client-go/types"
)

type operateCommand struct {
	readCommand

	policy     *WritePolicy
	operations []*Operation

	hasWrite bool
}

func newOperateCommand(cluster *Cluster, policy *WritePolicy, key *Key, operations []*Operation) operateCommand {
	return operateCommand{
		readCommand: newReadCommand(cluster, &policy.BasePolicy, key, nil),
		policy:      policy,
		operations:  operations,
	}
}

func (cmd *operateCommand) writeBuffer(ifc command) (err error) {
	cmd.hasWrite, err = cmd.setOperate(cmd.policy, cmd.key, cmd.operations)
	return err
}

func (cmd *operateCommand) getNode(ifc command) (*Node, error) {
	return cmd.cluster.getMasterNode(&cmd.partition)
}

func (cmd *operateCommand) Execute() error {
	return cmd.execute(cmd)
}

func (cmd *operateCommand) handleWriteKeyNotFoundError(resultCode ResultCode) error {
	if cmd.hasWrite {
		return NewAerospikeError(resultCode)
	}
	return nil
}
