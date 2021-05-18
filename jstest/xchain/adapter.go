package xchain

import (
	"testing"

	"github.com/xuperchain/xdev/jstest"
)

type xchainAdapter struct {
}

// NewAdapter is the xchain adapter
func NewAdapter() jstest.Adapter {
	return new(xchainAdapter)
}

func (x *xchainAdapter) OnSetup(r *jstest.Runner) {
	r.GlobalObject().Set("Xchain", func() *xchainObject {
		x, err := newXchainObject()
		if err != nil {
			jstest.Throw(err)
		}
		return x
	})
}

func (x *xchainAdapter) OnTeardown(r *jstest.Runner) {
}

func (x *xchainAdapter) OnTestCase(r *jstest.Runner, test jstest.TestCase) jstest.TestCase {
	body := func(t *testing.T) {
		xctx, err := newXchainObject()
		if err != nil {
			t.Fatal(err)
		}
		defer xctx.env.Close()

		if !r.Option.Quiet {
			// TODO: add log output
		}
		// reset xchain environment
		r.GlobalObject().Set("xchain", xctx)

		test.F(t)
	}
	return jstest.TestCase{
		Name: test.Name,
		F:    body,
	}
}
