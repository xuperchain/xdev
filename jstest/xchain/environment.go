package xchain

import (
	"encoding/json"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/protos"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	_ "github.com/xuperchain/xupercore/bcs/contract/evm"
	_ "github.com/xuperchain/xupercore/bcs/contract/native"
	_ "github.com/xuperchain/xupercore/bcs/contract/xvm"
	_ "github.com/xuperchain/xupercore/kernel/contract/kernel"
	_ "github.com/xuperchain/xupercore/kernel/contract/manager"
	_ "github.com/xuperchain/xupercore/lib/storage/kvdb/leveldb"
)

type environment struct {
	manager contract.Manager
	model   contract.StateSandbox
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

	state, err := m.NewStateSandbox(&contract.SandboxConfig{
		XMReader: store.State(),
	})
	if err != nil {
		return nil, err
	}
	return &environment{
		manager: m,
		model:   state,
		basedir: basedir,
	}, nil
}

type deployArgs struct {
	Name     string                 `json:"name"`
	Code     string                 `json:"code"`
	Lang     string                 `json:"lang"`
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
	descpb.Runtime = args.Lang
	descpb.ContractType = args.Type
	desc, err := proto.Marshal(descpb)
	if err != nil {
		return nil, err
	}
	dargs["contract_desc"] = desc

	ctx, err := e.manager.NewContext(&contract.ContextConfig{
		State:                 e.model,
		Initiator:             args.Options.Account,
		AuthRequire:           nil,
		Caller:                "",
		Module:                "xkernel",
		ContractName:          "$contract",
		ResourceLimits:        contract.MaxLimits,
		CanInitialize:         false,
		TransferAmount:        args.Options.Amount,
		ContractSet:           nil,
		ContractCodeFromCache: false,
	})
	resp, err := ctx.Invoke("deployContract", dargs)
	if err != nil {
		return nil, err
	}

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

//func (e *environment) ContractExists(name string) bool {
//
//	//ctx, err := e.xbridge.NewContext(&contract.ContextConfig{
//	//	State:          e.model,
//	//	ContractName:   name,
//	//	ResourceLimits: contract.MaxLimits,
//	//})
//	//if err != nil {
//	//	//TODO
//	//	return false
//	//}
//	////TODO defer ??
//	//ctx.Release()
//	//return true
//	return false
//}

func (e *environment) Invoke(name string, args invokeArgs) (*ContractResponse, error) {
	ctx, err := e.manager.NewContext(&contract.ContextConfig{
		State:                 e.model,
		Initiator:             "",
		AuthRequire:           nil,
		Caller:                "",
		Module:                "",
		ContractName:          name,
		ResourceLimits:        contract.MaxLimits,
		CanInitialize:         false,
		TransferAmount:        "",
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
	//TODO
	defer ctx.Release()

	if resp.Status >= contract.StatusErrorThreshold {
		return newContractResponse(resp), nil
	}

	if err := e.model.Flush(); err != nil {
		return nil, err
	}

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
