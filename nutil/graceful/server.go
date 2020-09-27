package graceful

import "context"

// ShutdownServer -
type ShutdownServer interface {
	Serve() error
	MustServe()
	Shutdown(ctx context.Context) error
}
