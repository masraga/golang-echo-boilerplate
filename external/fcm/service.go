package fcm

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/masraga/kerp-api/internal/service/notification"
	"github.com/rs/zerolog"
)

type FcmService struct {
	Ctx              context.Context
	ServiceAccountId ConfigServiceAccountId
	Logger           zerolog.Logger
	client           *messaging.Client
}

type FcmServiceStructOpts struct {
	Ctx              context.Context
	ServiceAccountId ConfigServiceAccountId
	Logger           zerolog.Logger
}

func NewFcmService(opts FcmServiceStructOpts) *FcmService {
	config := &firebase.Config{ProjectID: string(opts.ServiceAccountId)}
	app, err := firebase.NewApp(opts.Ctx, config)
	if err != nil {
		err = errors.Join(err, notification.ErrInitiateNotifProvider)
		fmt.Println(err)
		opts.Logger.Err(err).Msg(err.Error())
		return nil
	}
	client, err := app.Messaging(opts.Ctx)
	if err != nil {
		err = errors.Join(err, notification.ErrInitiateNotifProvider)
		fmt.Println(err)
		opts.Logger.Err(err).Msg(err.Error())
		return nil
	}
	return &FcmService{
		Ctx:              opts.Ctx,
		ServiceAccountId: opts.ServiceAccountId,
		client:           client,
	}
}
