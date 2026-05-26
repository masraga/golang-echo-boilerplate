package ctxerr

import (
	"github.com/rs/zerolog"
)

type CtxErr struct {
	ConfigShowErrMode ShowErrMode
	Logger            zerolog.Logger
}

type CtxErrOpts struct {
	ConfigShowErrMode ShowErrMode
	Logger            zerolog.Logger
}

func NewCtxErr(opt CtxErrOpts) *CtxErr {
	return &CtxErr{
		ConfigShowErrMode: opt.ConfigShowErrMode,
		Logger:            opt.Logger,
	}
}

func (c *CtxErr) Wrap(err error) error {
	if err == nil {
		return nil
	}
	c.Logger.Error().CallerSkipFrame(1).Caller().Msg(err.Error())
	return err
}
