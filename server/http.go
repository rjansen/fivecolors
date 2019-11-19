package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rjansen/l"
)

func Start(ctx context.Context, httpServer *http.Server) {
	l.Info(ctx, "server.starting")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		l.Info(ctx, "listener.starting", l.NewValue("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			l.Error(ctx, "listener.error", l.NewValue("error", err), l.NewValue("address", httpServer.Addr))
			panic(err)
		}
	}()

	<-stop

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l.Info(ctx, "server.shutdown")
	httpServer.Shutdown(shutDownCtx)
	l.Info(ctx, "shutdown.gracefully")
}
