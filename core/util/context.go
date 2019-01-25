package util

import (
	"context"
	"fmt"
	"time"

	"github.com/rjansen/fivecolors/core/validator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultTimeout = time.Minute * 2
)

type Context struct {
	context.Context
	zerolog.Logger
	Cancel    func()
	ID        string
	CreatedAt time.Time
	Timeout   time.Duration
}

func (c Context) String() string {
	return fmt.Sprintf("util.Context:{ ID:%s CreatedAt:%s Timeout:%d }", c.ID, c.CreatedAt.Format(time.RFC3339), c.Timeout)
}

func MakeContext() *Context {
	return CreateContext(context.Background(), defaultTimeout)
}

func NewContext(c context.Context) *Context {
	return CreateContext(c, defaultTimeout)
}

func CreateContext(p context.Context, timeout time.Duration) *Context {
	root, cancel := context.WithTimeout(p, timeout)
	tid, okstr := root.Value("tid").(string)
	if !okstr || validator.IsBlank(tid) != nil {
		tid = NewUUID()
		root = context.WithValue(root, "tid", tid)
	}
	createdAt, oktime := p.Value("createdAt").(time.Time)
	if !oktime || createdAt.IsZero() {
		createdAt = time.Now()
		root = context.WithValue(root, "createdAt", createdAt)
	}
	return &Context{
		ID:        tid,
		CreatedAt: createdAt,
		Timeout:   timeout,
		Logger:    log.With().Str("tid", tid).Logger(),
		Context:   root,
		Cancel:    cancel,
	}
}
