package xchain

import (
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/contract/bridge"
)

type kcontextImpl struct {
	ctx     *bridge.Context
	syscall *bridge.SyscallService
	contract.StateSandbox
	contract.ChainCore
	used, limit contract.Limits
}

func newKContext(ctx *bridge.Context, syscall *bridge.SyscallService) *kcontextImpl {
	return &kcontextImpl{
		ctx:          ctx,
		syscall:      syscall,
		limit:        ctx.ResourceLimits,
		StateSandbox: ctx.State,
		ChainCore:    ctx.Core,
	}
}

// 交易相关数据
func (k *kcontextImpl) Args() map[string][]byte {
	return k.ctx.Args
}

func (k *kcontextImpl) Initiator() string {
	return k.ctx.Initiator
}

func (k *kcontextImpl) Caller() string {
	return k.ctx.Caller
}

func (k *kcontextImpl) AuthRequire() []string {
	return k.ctx.AuthRequire
}

func (k *kcontextImpl) AddResourceUsed(delta contract.Limits) {
	k.used.Add(delta)
}

func (k *kcontextImpl) ResourceLimit() contract.Limits {
	return k.limit
}

func (k *kcontextImpl) Call(module, contractName, method string, args map[string][]byte) (*contract.Response, error) {
	return &contract.Response{
		Status:  200,
		Message: "ok",
		Body:    nil,
	}, nil
}