// Copyright 2021 The nfgo Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nerrors

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// ErrGroup - A ErrGroup is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero ErrGroup is valid and does not cancel on error.
type ErrGroup struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// NewErrGroup - NewErrGroup returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func NewErrGroup(ctx context.Context) (*ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{
		cancel: cancel,
	}, ctx
}

// Go - Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *ErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go g.doGo(f)
}

func (g *ErrGroup) doGo(f func() error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("nerrors errgroup panic recoverd: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = f()
}

// Wait - Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
