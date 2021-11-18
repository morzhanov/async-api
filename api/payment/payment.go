package payment

import (
	"io/ioutil"

	"github.com/swaggest/go-asyncapi/reflector/asyncapi-2.0.0"
	"github.com/swaggest/go-asyncapi/spec-2.0.0"
)

type ProcessPaymentMessage struct {
	OrderID string `json:"orderId" path:"orderId" description:"Order ID"`
	Name    string `json:"name" path:"name" description:"Order Name"`
	Amount  int    `json:"amount" path:"amount" description:"Order Amount"`
	Status  string `json:"status" path:"status" description:"Order Status"`
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
		Name: "payment.process",
		Subscribe: &asyncapi.MessageSample{
			MessageEntity: spec.MessageEntity{
				Description: "Process Payment consumer",
				Summary:     "Process Payment consumer",
			},
			MessageSample: new(ProcessPaymentMessage),
		},
	})
	if err != nil {
		return err
	}

	yaml, err := reflector.Schema.MarshalYAML()
	if err != nil {
		return err
	}
	return ioutil.WriteFile("payment_asyncapi_spec.yaml", yaml, 0644)
}
