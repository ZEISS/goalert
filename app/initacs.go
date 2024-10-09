package app

import (
	"context"

	"github.com/target/goalert/notification/acs"

	"github.com/pkg/errors"
)

func (app *App) initAcs(ctx context.Context) error {
	app.acsConfig = &acs.Config{}

	var err error
	app.acsSMS, err = acs.NewSMS(ctx, acs.DefaultClient, app.db, app.acsConfig)
	if err != nil {
		return errors.Wrap(err, "init AcsSMS")
	}

	// app.twilioVoice, err = twilio.NewVoice(ctx, app.db, app.twilioConfig)
	// if err != nil {
	// 	return errors.Wrap(err, "init TwilioVoice")
	// }

	return nil
}
