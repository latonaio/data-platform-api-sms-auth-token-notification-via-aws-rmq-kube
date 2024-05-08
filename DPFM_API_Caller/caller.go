package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-sms-auth-token-notification-via-aws-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-sms-auth-token-notification-via-aws-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-sms-auth-token-notification-via-aws-rmq-kube/config"
	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

type DPFMAPICaller struct {
	ctx  context.Context
	conf *config.Conf
	rmq  *rabbitmq.RabbitmqClient
}

func NewDPFMAPICaller(
	conf *config.Conf,
	rmq *rabbitmq.RabbitmqClient,
) *DPFMAPICaller {
	return &DPFMAPICaller{
		ctx:  context.Background(),
		conf: conf,
		rmq:  rmq,
	}
}

func (c *DPFMAPICaller) AsyncCreates(
	accepter []string,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) (interface{}, []error) {
	var smsAuthToken *[]dpfm_api_output_formatter.SMSAuthToken
	var errs []error

	for _, fn := range accepter {
		switch fn {
		case "SMSAuthToken":
			smsAuthToken = c.SMSAuthToken(input, &errs, log, c.conf)
		default:
		}
	}

	data := &dpfm_api_output_formatter.Message{
		SMSAuthToken: smsAuthToken,
	}

	return data, nil
}
