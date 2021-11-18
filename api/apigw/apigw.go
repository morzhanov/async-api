package apigw

import (
	"io/ioutil"

	"github.com/morzhanov/async-api/api/order"

	"github.com/swaggest/go-asyncapi/reflector/asyncapi-2.0.0"
	"github.com/swaggest/go-asyncapi/spec-2.0.0"
)

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
				Description: "Create Order producer",
				Summary:     "Create Order producer",
			},
			MessageSample: new(order.CreateOrderMessage),
		},
	})
	if err != nil {
		return err
	}
	err = reflector.AddChannel(asyncapi.ChannelInfo{
		Name: "order.process",
		Subscribe: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Process Order producer",
				Summary:     "Process Order producer",
			},
			MessageSample: new(order.ProcessOrderMessage),
		},
	})
	if err != nil {
		return err
	}

	yaml, err := reflector.Schema.MarshalYAML()
	if err != nil {
		return err
	}
	return ioutil.WriteFile("apigw_asyncapi_spec.yaml", yaml, 0644)
}
