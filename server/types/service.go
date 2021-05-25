// Package types contains RPC server wrapper types.
package types

import "github.com/diamondburned/cchat"

type Service struct {
	V cchat.Service
}

func (s Service) ID(_ struct{}, out *cchat.ID) error {
	*out = s.V.ID()
	return nil
}
