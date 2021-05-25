package xchain

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/robertkrimen/otto"
	"github.com/xuperchain/xdev/jstest"
	"github.com/xuperchain/xupercore/bcs/contract/evm/abi"
	"github.com/xuperchain/xupercore/kernel/contract/bridge"
)

type contractObject struct {
	env  *environment
	abi  *abi.ABI
	Name string
	Type string
}

// func (c *contractObject) Invoke(method string, args map[string]string, option InvokeOptions) *contract.Response {
func (c *contractObject) Invoke(call otto.FunctionCall) otto.Value {
	var args invokeArgs

	method := call.Argument(0).String()
	args.Method = method

	if !call.Argument(1).IsObject() {
		jstest.Throws("expect method args with object type")
	}
	export, _ := call.Argument(1).Export()
	err := mapstructure.Decode(export, &args.Args)
	if err != nil {
		jstest.Throw(err)
	}
	if c.Type != string(bridge.TypeEvm) {
		args.trueArgs = convertArgs(args.Args)
	} else {
		if method != "" {
			input, err := c.abi.Encode(method, args.Args)
			if err != nil {
				jstest.Throw(fmt.Errorf("abi encode error:%s", err))
			}
			args.trueArgs = map[string][]byte{
				"input": input,
			}
		}
	}

	if call.Argument(2).IsObject() {
		export, _ := call.Argument(2).Export()
		err := mapstructure.Decode(export, &args.Options)
		if err != nil {
			jstest.Throw(err)
		}
	}

	resp, err := c.env.Invoke(c.Name, args)
	if err != nil {
		jstest.Throw(err)
	}
	v, err := call.Otto.ToValue(resp)
	if err != nil {
		jstest.Throw(err)
	}
	return v
}
