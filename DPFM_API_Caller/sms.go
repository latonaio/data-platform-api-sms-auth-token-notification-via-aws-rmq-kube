package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-sms-auth-token-notification-via-aws-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-sms-auth-token-notification-via-aws-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-sms-auth-token-notification-via-aws-rmq-kube/config"
	"fmt"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	pinpoint "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	"golang.org/x/xerrors"
)

func (c *DPFMAPICaller) SMSAuthToken(
	input *dpfm_api_input_reader.SDC,
	errs *[]error,
	log *logger.Logger,
	conf *config.Conf,
) *[]dpfm_api_output_formatter.SMSAuthToken {
	inputSmsAuthToken := input.SMSAuthToken

	if inputSmsAuthToken == nil && len(*inputSmsAuthToken) == 0 {
		err := xerrors.Errorf("SMSAuthToken is empty")
		*errs = append(*errs, err)
		return nil
	}

	mobilePhoneNumber := (*inputSmsAuthToken)[0].MobilePhoneNumber
	authenticationCode := (*inputSmsAuthToken)[0].AuthenticationCode

	err := postSmsAws(mobilePhoneNumber,
		fmt.Sprintf("あなたの認証コードは: %d です。", authenticationCode),
		conf,
	)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	var sMSAuthTokens []dpfm_api_output_formatter.SMSAuthToken

	sMSAuthTokens = append(sMSAuthTokens, dpfm_api_output_formatter.SMSAuthToken{
		MobilePhoneNumber:  mobilePhoneNumber,
		AuthenticationCode: authenticationCode,
	})

	return &sMSAuthTokens
}

func postSmsAws(
	recipient string,
	message string,
	conf *config.Conf,
) error {
	cfg, err := aws_config.LoadDefaultConfig(context.TODO())
	if err != nil {
		err = xerrors.Errorf("Error loading AWS configuration: %v", err)
		return err
	}

	client := pinpoint.NewFromConfig(cfg)

	senderId := conf.AWS.AWSPinpointSenderID

	input := pinpoint.SendTextMessageInput{
		DestinationPhoneNumber: &recipient,
		MessageBody:            &message,
		//		MessageType:            pinpoint.MessageTypeTransactional,
		OriginationIdentity: &senderId,
		DryRun:              false,
	}

	_, err = client.SendTextMessage(context.TODO(), &input)

	if err != nil {
		err = xerrors.Errorf("Error SendTextMessage: %v", err)
		return err
	}

	return nil
}
