package types

import (
	"net/rpc"

	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-netrpc/server"
)

func init() {
	rpc.RegisterName("netrpc/query", Query{})
	rpc.RegisterName("cchat/com.diamondbured.aaa/Server/324982394234", nil)
}

// Query queries for the internal repository of plugins.
type Query struct{}

// Services returns a list of service IDs.
func (q *Query) Services(in struct{}, out *[]cchat.ID) error {
	*out = server.ServiceIDs()
	return nil
}
