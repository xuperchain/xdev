package xchain

import (
	"encoding/json"
	"github.com/xuperchain/xupercore/kernel/contract/sandbox"
	"github.com/xuperchain/xupercore/protos"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/contract/bridge"
	//TODO
	_ "github.com/xuperchain/xupercore/bcs/contract/evm"
	_ "github.com/xuperchain/xupercore/bcs/contract/native"
	_ "github.com/xuperchain/xupercore/bcs/contract/xvm"
)

type environment struct {
	xbridge *bridge.XBridge
	model   contract.StateSandbox
	basedir string
}

func newEnvironment() (*environment, error) {
	basedir, err := ioutil.TempDir("", "xdev-env")
	if err != nil {
		return nil, err
	}
	//store := newMockStore()
	store := sandbox.NewMemXModel()
	if err != nil {
		return nil, err
	}
	vmconfig := contract.DefaultContractConfig()
	wasmConfig := vmconfig.Wasm
	wasmConfig.Driver = "ixvm"

	xbridge, err := bridge.New(&bridge.XBridgeConfig{
		Basedir: basedir,
		VMConfigs: map[bridge.ContractType]bridge.VMConfig{
			bridge.TypeWasm:   &wasmConfig,
			bridge.TypeNative: &vmconfig.Native,
			bridge.TypeEvm:    &vmconfig.EVM,
		},
		Core:      &chainCore{},
		XModel:    store,
		LogWriter: os.Stderr,
	})
	if err != nil {
		os.RemoveAll(basedir)
		return nil, err
	}

	return &environment{
		xbridge: xbridge,
		model:   sandbox.NewXModelCache(store),
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

	//ctx := &bridge.Context{
	//	ContractName:   args.Name,
	//	ResourceLimits: contract.MaxLimits,
	//	Args:           dargs,
	//	CanInitialize:  true,
	//}

	//kctx := &kcontextImpl{
	//	ctx:          ctx,
	//	syscall:      nil,
	//	StateSandbox: e.model,
	//	ChainCore:    new(chainCore),
	//	used:         contract.Limits{0, 0, 0, 0},
	//	limit:        contract.MaxLimits,
	//}
	ctx, err := e.xbridge.NewContext(&contract.ContextConfig{
		//State:                 nil,
		//Initiator:             "",
		//AuthRequire:           nil,
		//Caller:                "",
		Module: "xkernel",
		//ContractName: "deploy",
		//ResourceLimits:        contract.Limits{},
		//CanInitialize:         false,
		//TransferAmount:        "",
		//ContractSet:           nil,
		//ContractCodeFromCache: false,
	})
	resp, err := ctx.Invoke("deploy", map[string][]byte{})
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

func (e *environment) ContractExists(name string) bool {

	ctx, err := e.xbridge.NewContext(&contract.ContextConfig{
		State:          e.model,
		ContractName:   name,
		ResourceLimits: contract.MaxLimits,
	})
	if err != nil {
		//TODO
		return false
	}
	//TODO defer ??
	ctx.Release()
	return true
}

func (e *environment) Invoke(name string, args invokeArgs) (*ContractResponse, error) {
	ctx, err := e.xbridge.NewContext(&contract.ContextConfig{
		State:          e.model,
		Initiator:      args.Options.Account,
		TransferAmount: args.Options.Amount,
		ContractName:   name,

		ResourceLimits: contract.MaxLimits,
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
