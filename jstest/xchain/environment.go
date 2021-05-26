package xchain

import (
	"encoding/json"
	"github.com/xuperchain/xupercore/kernel/common/xcontext"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/permission/acl"
	actx "github.com/xuperchain/xupercore/kernel/permission/acl/context"
	"github.com/xuperchain/xupercore/lib/logs"
	"github.com/xuperchain/xupercore/protos"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"

	_ "github.com/xuperchain/xupercore/bcs/consensus/pow"
	_ "github.com/xuperchain/xupercore/bcs/consensus/single"
	_ "github.com/xuperchain/xupercore/bcs/consensus/tdpos"
	_ "github.com/xuperchain/xupercore/bcs/consensus/xpoa"
	_ "github.com/xuperchain/xupercore/bcs/contract/evm"
	_ "github.com/xuperchain/xupercore/bcs/contract/native"
	_ "github.com/xuperchain/xupercore/bcs/contract/xvm"
	_ "github.com/xuperchain/xupercore/bcs/network/p2pv1"
	_ "github.com/xuperchain/xupercore/bcs/network/p2pv2"
	_ "github.com/xuperchain/xupercore/kernel/contract/kernel"
	_ "github.com/xuperchain/xupercore/kernel/contract/manager"
	_ "github.com/xuperchain/xupercore/lib/crypto/client"
	_ "github.com/xuperchain/xupercore/lib/storage/kvdb/leveldb"
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
	_, _ = acl.NewACLManager(&actx.AclCtx{
		BaseCtx:  xcontext.BaseCtx{},
		BcName:   "xuper",
		Ledger:   &MockLedgerRely{&XMSnapshotReader{}},
		Contract: m,
	})

	e := &environment{
		manager: m,
		store:   store,
		basedir: basedir,
	}
	if err := e.InitAccount(); err != nil {
		return nil, err
	}
	if err := e.InitLog(); err != nil {
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

//InitLog init xupercore logger config to ignore non-crit logs and disable console output
func (e *environment) InitLog() error {
	logDir, err := ioutil.TempDir("", "xdev-log")
	if err != nil {
		return err
	}
	confDir, err := ioutil.TempDir("", "xdev-conf")
	if err != nil {
		return err
	}

	confPath := confDir + "/logs.yaml"
	xdevlog := `
level: crit
console: false
`
	if err := ioutil.WriteFile(confPath, []byte(xdevlog), 0755); err != nil {
		return err
	}

	logs.InitLog(confPath, logDir)
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
	resp, err := ctx.Invoke(args.Method, args.trueArgs)
	if err != nil {
		return nil, err
	}
	defer ctx.Release()

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
