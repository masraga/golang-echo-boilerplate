package ctxerr

import (
	"github.com/rs/zerolog"
)

type CtxErr struct {
	Logger zerolog.Logger
}

type CtxErrOpts struct {
	Logger zerolog.Logger
}

func NewCtxErr(opt CtxErrOpts) *CtxErr {
	return &CtxErr{
		Logger: opt.Logger,
	}
}

func (c *CtxErr) Wrap(err error) error {
	if err == nil {
		return nil
	}
	c.Logger.Error().Msg(err.Error())
	return err
}
