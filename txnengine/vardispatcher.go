package txnengine

import (
	"fmt"
	"goshawkdb.io/common"
	cmsgs "goshawkdb.io/common/capnp"
	"goshawkdb.io/server"
	msgs "goshawkdb.io/server/capnp"
	"goshawkdb.io/server/configuration"
	"goshawkdb.io/server/db"
	"goshawkdb.io/server/dispatcher"
)

type TopologyPublisher interface {
	AddTopologySubscriber(obs TopologySubscriber) *configuration.Topology
	RemoveTopologySubscriberAsync(obs TopologySubscriber)
}

type TopologySubscriber interface {
	TopologyChanged(TopologyChange)
}

type TopologyChangeSubscriberType uint8

const (
	VarSubscriber                     TopologyChangeSubscriberType = iota
	ProposerSubscriber                TopologyChangeSubscriberType = iota
	AcceptorSubscriber                TopologyChangeSubscriberType = iota
	ConnectionManagerSubscriber       TopologyChangeSubscriberType = iota
	TopologyChangeSubscriberTypeLimit int                          = iota
)

type TopologyChange interface {
	Topology() *configuration.Topology
	Done(TopologyChangeSubscriberType)
	AddOne(TopologyChangeSubscriberType)
	HasCallbackFor(TopologyChangeSubscriberType) bool
}

type VarDispatcher struct {
	dispatcher.Dispatcher
	varmanagers []*VarManager
}

func NewVarDispatcher(count uint8, rmId common.RMId, cm TopologyPublisher, db *db.Databases, lc LocalConnection) *VarDispatcher {
	vd := &VarDispatcher{
		varmanagers: make([]*VarManager, count),
	}
	vd.Dispatcher.Init(count)
	for idx, exe := range vd.Executors {
		vd.varmanagers[idx] = NewVarManager(exe, rmId, cm, db, lc)
	}
	return vd
}

func (vd *VarDispatcher) ApplyToVar(fun func(*Var, error), createIfMissing bool, vUUId *common.VarUUId) {
	vd.withVarManager(vUUId, func(vm *VarManager) { vm.ApplyToVar(fun, createIfMissing, vUUId) })
}

func (vd *VarDispatcher) Status(sc *server.StatusConsumer) {
	sc.Emit("Vars")
	for idx, executor := range vd.Executors {
		s := sc.Fork()
		s.Emit(fmt.Sprintf("Var Manager %v", idx))
		manager := vd.varmanagers[idx]
		executor.Enqueue(func() { manager.Status(s) })
	}
	sc.Join()
}

func (vd *VarDispatcher) withVarManager(vUUId *common.VarUUId, fun func(*VarManager)) bool {
	idx := uint8(vUUId[server.MostRandomByteIndex]) % vd.ExecutorCount
	executor := vd.Executors[idx]
	manager := vd.varmanagers[idx]
	return executor.Enqueue(func() { fun(manager) })
}

type LocalConnection interface {
	RunClientTransaction(txn *cmsgs.ClientTxn, varPosMap map[common.VarUUId]*common.Positions, assignTxnId bool) (*msgs.Outcome, error)
	Status(*server.StatusConsumer)
}
