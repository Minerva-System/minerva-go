package schema

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type BooleanResponse struct {
	Value bool `json:"value"`
}

func (BooleanResponse) FromMessage(b *wrapperspb.BoolValue) BooleanResponse {
	return BooleanResponse{ Value: b.Value }
}
