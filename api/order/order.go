package order

import (
	"io/ioutil"

	"github.com/morzhanov/async-api/api/payment"

	"github.com/swaggest/go-asyncapi/reflector/asyncapi-2.0.0"
	"github.com/swaggest/go-asyncapi/spec-2.0.0"
)

type CreateOrderMessage struct {
	Name   string `json:"name" path:"name" description:"Order Name"`
	Amount int    `json:"amount" path:"amount" description:"Order Amount"`
}

type ProcessOrderMessage struct {
	ID string `json:"id" path:"id" description:"Order ID"`
}

func Build(url string, protocol string, protocolVersion string) error {
	reflector := asyncapi.Reflector{
		Schema: &spec.AsyncAPI{
			Servers: map[string]spec.Server{
				"live": {
					URL:             url,
					ProtocolVersion: protocolVersion,
					Protocol:        protocol,
				},
			},
		},
	}

	err := reflector.AddChannel(asyncapi.ChannelInfo{
		Name: "order.create",
		Subscribe: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Create Order consumer",
				Summary:     "Create Order consumer",
			},
			MessageSample: new(CreateOrderMessage),
		},
	})
	if err != nil {
		return err
	}
	err = reflector.AddChannel(asyncapi.ChannelInfo{
		Name: "order.process",
		Subscribe: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Process Order consumer",
				Summary:     "Process Order consumer",
			},
			MessageSample: new(ProcessOrderMessage),
		},
	})
	if err != nil {
		return err
	}
	err = reflector.AddChannel(asyncapi.ChannelInfo{
		Name: "payment.process",
		Subscribe: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Process Payment producer",
				Summary:     "Process Payment producer",
			},
			MessageSample: new(payment.ProcessPaymentMessage),
		},
	})
	if err != nil {
		return err
	}

	yaml, err := reflector.Schema.MarshalYAML()
	if err != nil {
		return err
	}
	return ioutil.WriteFile("order_asyncapi_spec.yaml", yaml, 0644)
}
