package context

import (
	"errors"
	"farm.e-pedion.com/repo/logger"
)

const (
	ContentTypeHeader   = "Content-Type"
	ContentLengthHeader = "Content-Length"
	AcceptHeader        = "Accept"
)

var (
	ErrInvalidContentType = errors.New("Invalid ContentType. Only: aplication/json, application/octet-stream are valid")
	ErrInvalidAccept      = errors.New("Invalid Accept. Only: aplication/json, application/octet-stream are valid")
)

//SetupAll calls all provided setup functions and return all raised errors
func SetupAll(setupFuncs ...SetupFunc) []error {
	var errs []error
	for i, v := range setupFuncs {
		if err := v(); err != nil {
			logger.Warn("contex.SetupSilent",
				logger.Int("index", i),
				logger.Struct("func", v),
				logger.Err(err),
			)
			errs = append(errs, err)
		}
	}
	return errs
}

//Setup calls the provided setup functions and return at the first raised error
func Setup(setupFuncs ...SetupFunc) error {
	for i, v := range setupFuncs {
		if err := v(); err != nil {
			logger.Error("contex.Setup",
				logger.Int("index", i),
				logger.Struct("func", v),
				logger.Err(err),
			)
			return err
		}
	}
	return nil
}
