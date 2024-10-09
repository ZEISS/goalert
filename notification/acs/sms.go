package acs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/target/goalert/config"
	"github.com/target/goalert/gadb"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/nfydest"
	"github.com/target/goalert/notification/nfymsg"
	"github.com/target/goalert/util/log"
	acs "github.com/zeiss/go-acs"
	"github.com/zeiss/go-acs/sms"
)

var (
	ErrProviderDisabled          = errors.New("acsc: provider is disabled")
	ErrProviderRefusedFromNumber = errors.New("acsc: refusing to send message to FromNumber")
)

func NewSMSDest(number string) gadb.DestV1 {
	return gadb.NewDestV1(DestTypeAcsSMS, FieldPhoneNumber, number)
}

// SMS implements a notification.Sender for Azure Communication Services SMS.
type SMS struct {
	cfg *Config
	c   *acs.Client
	db  *dbSMS

	r notification.Receiver
}

var (
	_ notification.ReceiverSetter = &SMS{}
	_ nfydest.MessageSender       = &SMS{}
	_ nfydest.MessageStatuser     = &SMS{}
)

// NewSMS performs operations like validating essential parameters,
// registering the Azur Communication Services SMS providerm
// and adding routes for successful and unsuccessful message delivery to Azure Communication Services.
func NewSMS(ctx context.Context, client *http.Client, db *sql.DB, c *Config) (*SMS, error) {
	s := new(SMS)

	cfg := config.FromContext(ctx)

	b, err := newDB(ctx, db)
	if err != nil {
		return nil, err
	}

	s.db = b
	s.cfg = c

	s.c = acs.New(cfg.AzureCommunicationServices.Endpoint, cfg.AzureCommunicationServices.Key, client)

	return s, nil
}

// SetReceiver sets the notification.Receiver for incoming messages and status updates.
func (s *SMS) SetReceiver(r notification.Receiver) { s.r = r }

// Status provides the current status of a message.
func (s *SMS) MessageStatus(ctx context.Context, externalID string) (*notification.Status, error) {
	return nil, nil
}

// Send implements the notification.Sender interface.
func (s *SMS) SendMessage(ctx context.Context, msg notification.Message) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)
	if !cfg.AzureCommunicationServices.Enable {
		return nil, ErrProviderDisabled
	}

	if msg.DestType() != DestTypeAcsSMS {
		return nil, fmt.Errorf("unsupported destination type %s; expected SMS", msg.DestType())
	}

	destNumber := msg.DestArg(FieldPhoneNumber)
	if destNumber == cfg.AzureCommunicationServices.FromNumber {
		return nil, ErrProviderRefusedFromNumber
	}

	ctx = log.WithFields(ctx, log.Fields{
		"Phone": destNumber,
		"Type":  "AcsSMS",
	})

	var message string
	var err error
	switch t := msg.(type) {
	case notification.Verification:
		message = fmt.Sprintf("%s: Verification code: %s", cfg.ApplicationName(), t.Code)
	default:
		return nil, fmt.Errorf("unhandled message type %T", t)
	}

	req := &sms.Request{
		From: cfg.AzureCommunicationServices.FromNumber,
		SMSRecipients: []sms.SMSRecipients{
			{
				To: destNumber,
			},
		},
		Message: message,
		SMSSendOptions: sms.SMSSendOptions{
			EnableDeliveryReport: true,
		},
	}

	resp, err := s.c.SMS.SendSMS(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("send message: %v", err)
	}

	if len(resp.Value) == 0 {
		return nil, errors.New("could not sent sms")
	}

	v := resp.Value[0]
	if !v.Successful {
		return nil, fmt.Errorf("send message failed: %s", v.ErrorMessage)
	}

	return &notification.SentMessage{
		ExternalID:   v.MessageID,
		State:        nfymsg.StateSent,
		StateDetails: v.ErrorMessage,
		SrcValue:     cfg.AzureCommunicationServices.FromNumber,
	}, nil
}

func (s *SMS) ServeStatusCallback(w http.ResponseWriter, req *http.Request) {
}

func (s *SMS) ServeMessage(w http.ResponseWriter, req *http.Request) {
}
