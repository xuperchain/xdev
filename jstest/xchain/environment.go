package xchain

import (
	"encoding/json"
	log15 "github.com/xuperchain/log15"
	"github.com/xuperchain/xupercore/kernel/common/xcontext"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/permission/acl"
	actx "github.com/xuperchain/xupercore/kernel/permission/acl/context"
	"github.com/xuperchain/xupercore/protos"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	_ "github.com/xuperchain/xupercore/bcs/contract/evm"
	_ "github.com/xuperchain/xupercore/bcs/contract/native"
	_ "github.com/xuperchain/xupercore/bcs/contract/xvm"
	_ "github.com/xuperchain/xupercore/kernel/contract/kernel"
	_ "github.com/xuperchain/xupercore/kernel/contract/manager"
)

type environment struct {
	manager contract.Manager
	store   *mockStore
	basedir string
}

func newEnvironment() (*environment, error) {
	basedir, err := ioutil.TempDir("", "xdev-env")
	if err != nil {
		return nil, err
	}
	store := NewmockStore()

	config := contract.DefaultContractConfig()
	config.Wasm.Driver = "ixvm"

	logger := log15.New()
	logger.SetHandler(log15.StreamHandler(os.Stderr, log15.LogfmtFormat()))

	config.LogDriver = logger

	m, err := contract.CreateManager("default", &contract.ManagerConfig{
		Basedir:  basedir,
		BCName:   "xuper",
		EnvConf:  nil,
		Core:     &chainCore{},
		XMReader: store.State(),
		Config:   config,
	})
	if err != nil {
		return nil, err
	}
	// To Register kernel contract $acl
	_, err = acl.NewACLManager(&actx.AclCtx{
		BaseCtx:  xcontext.BaseCtx{},
		BcName:   "xuper",
		Ledger:   &MockLedgerRely{&XMSnapshotReader{}},
		Contract: m,
	})

	if err != nil {
		return nil, err
	}
	e := &environment{
		manager: m,
		store:   store,
		basedir: basedir,
	}
	if err := e.InitAccount(); err != nil {
		return nil, err
	}
	return e, nil
}

type deployArgs struct {
	Name string `json:"name"`
	Code string `json:"code"`
	//Deprecated: using Runtime instead
	Lang string `json:"lang"`
	// Runtime specify runtime, has priority than lang
	Runtime  string                 `json:"runtime"`
	InitArgs map[string]interface{} `json:"init_args"`
	Type     string                 `json:"type"`
	ABIFile  string                 `json:"abi"`
	Options  invokeOptions          `json:"options"`

	trueArgs map[string][]byte
	codeBuf  []byte
}

func convertArgs(ori map[string]interface{}) map[string][]byte {
	ret := make(map[string][]byte)
	for k, v := range ori {
		ret[k] = []byte(v.(string))
	}
	return ret
}
func (e *environment) InitAccount() error {
	state, err := e.manager.NewStateSandbox(&contract.SandboxConfig{
		XMReader: e.store.State(),
	})
	ctx, err := e.manager.NewContext(&contract.ContextConfig{
		State:                 state,
		Initiator:             "",
		AuthRequire:           nil,
		Caller:                "",
		Module:                "xkernel",
		ContractName:          "$acl",
		ResourceLimits:        contract.MaxLimits,
		CanInitialize:         true,
		TransferAmount:        "0",
		ContractSet:           nil,
		ContractCodeFromCache: false,
	})
	if err != nil {
		return err
	}
	_, err = ctx.Invoke("NewAccount", map[string][]byte{
		"acl":          []byte(defaultAccountACL),
		"account_name": []byte(defaultTestingAccount),
	})
	if err != nil {
		return err
	}
	e.store.Commit(state)
	return nil
}

func (e *environment) Deploy(args deployArgs) (*ContractResponse, error) {
	dargs := make(map[string][]byte)
	dargs["contract_name"] = []byte(args.Name)
	dargs["contract_code"] = args.codeBuf
	dargs["account_name"] = []byte(args.Options.Account)

	initArgs, err := json.Marshal(args.trueArgs)
	if err != nil {
		return nil, err
	}
	dargs["init_args"] = initArgs

	descpb := new(protos.WasmCodeDesc)
	descpb.Runtime = args.Runtime
	if descpb.Runtime == "" {
		descpb.Runtime = args.Lang
	}
	descpb.ContractType = args.Type
	desc, err := proto.Marshal(descpb)
	if err != nil {
		return nil, err
	}
	dargs["contract_desc"] = desc

	state, err := e.manager.NewStateSandbox(&contract.SandboxConfig{
		XMReader: e.store.State(),
	})
	ctx, err := e.manager.NewContext(&contract.ContextConfig{
		State:                 state,
		Initiator:             args.Options.Account,
		AuthRequire:           nil,
		Caller:                "",
		Module:                "xkernel",
		ContractName:          "$contract",
		ResourceLimits:        contract.MaxLimits,
		CanInitialize:         true,
		TransferAmount:        args.Options.Amount,
		ContractSet:           nil,
		ContractCodeFromCache: false,
	})
	defer ctx.Release()
	resp, err := ctx.Invoke("deployContract", dargs)
	if err != nil {
		return nil, err
	}
	e.store.Commit(state)
	return newContractResponse(resp), nil
}

type invokeOptions struct {
	Account string `json:"account"`
	Amount  string `json:"amount"`
}

type invokeArgs struct {
	Method   string                 `json:"method"`
	Args     map[string]interface{} `json:"args"`
	trueArgs map[string][]byte
	Options  invokeOptions
}

func (e *environment) ContractExists(name string) bool {
	ctx, err := e.manager.NewContext(&contract.ContextConfig{
		State:                 nil,
		Initiator:             "",
		AuthRequire:           nil,
		Caller:                "",
		Module:                "",
		ContractName:          name,
		ResourceLimits:        contract.Limits{},
		CanInitialize:         false,
		TransferAmount:        "",
		ContractSet:           nil,
		ContractCodeFromCache: false,
	})
	defer ctx.Release()
	if err != nil {
		return false
	}
	return true
}

func (e *environment) Invoke(name string, args invokeArgs) (*ContractResponse, error) {
	state, err := e.manager.NewStateSandbox(&contract.SandboxConfig{
		XMReader: e.store.State(),
	})
	ctx, err := e.manager.NewContext(&contract.ContextConfig{
		State:                 state,
		Initiator:             args.Options.Account,
		AuthRequire:           nil,
		Caller:                "",
		Module:                "",
		ContractName:          name,
		ResourceLimits:        contract.MaxLimits,
		CanInitialize:         false,
		TransferAmount:        args.Options.Amount,
		ContractSet:           nil,
		ContractCodeFromCache: false,
	})
	if err != nil {
		return nil, err
	}
	defer ctx.Release()

	resp, err := ctx.Invoke(args.Method, args.trueArgs)
	if err != nil {
		return nil, err
	}

	if resp.Status >= contract.StatusErrorThreshold {
		return newContractResponse(resp), nil
	}
	e.store.Commit(state)
	return newContractResponse(resp), nil
}

func (e *environment) Close() {
	os.RemoveAll(e.basedir)
}

type ContractResponse struct {
	Status  int
	Message string
	Body    string
}

func newContractResponse(resp *contract.Response) *ContractResponse {
	return &ContractResponse{
		Status:  resp.Status,
		Message: resp.Message,
		Body:    string(resp.Body),
	}
}
