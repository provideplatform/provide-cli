/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shell

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/provideplatform/provide-go/common"
)

const replTickerInterval = 500 * time.Millisecond
const replSleepInterval = 250 * time.Millisecond
const replCmdSleepInterval = 100 * time.Millisecond

type REPL struct {
	buf         *bytes.Buffer
	cancelF     context.CancelFunc
	closing     uint32
	cmd         *exec.Cmd
	executing   bool
	finished    bool
	fn          func(*sync.WaitGroup) error
	io          io.ReadWriteCloser // FIXME
	mutex       *sync.Mutex
	shutdownCtx context.Context
	wg          *sync.WaitGroup
}

func NewREPL(fn func(*sync.WaitGroup) error) (*REPL, error) {
	repl := &REPL{
		fn:    fn,
		mutex: &sync.Mutex{},
		wg:    &sync.WaitGroup{},
	}

	return repl, nil
}

func NewREPLWithCmd(cmd exec.Cmd, buf *bytes.Buffer) (*REPL, error) {
	repl := &REPL{
		buf:   buf,
		cmd:   &cmd,
		mutex: &sync.Mutex{},
		wg:    &sync.WaitGroup{},
	}

	return repl, nil
}

func (r *REPL) installSignalHandlers() chan os.Signal {
	common.Log.Tracef("installing subshell signal handlers")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return sigs
}

func (r *REPL) run() {
	sigs := r.installSignalHandlers()
	r.shutdownCtx, r.cancelF = context.WithCancel(context.Background())

	timer := time.NewTicker(replTickerInterval)
	defer timer.Stop()

	for !r.shuttingDown() {
		select {
		case <-timer.C:
			if r.cmd != nil && !r.executing && !r.finished {
				err := r.exec()
				if err != nil {
					common.Log.Warningf("runloop exec() returned err; %s", err.Error())
				}
			} else if r.fn != nil {
				err := r.fn(r.wg)
				if err != nil {
					common.Log.Warningf("runloop fn() returned err; %s", err.Error())
				}
			}
		case sig := <-sigs:
			common.Log.Infof("received signal: %s", sig)
			r.shutdown()
		case <-r.shutdownCtx.Done():
			close(sigs)
		default:
			time.Sleep(replSleepInterval)
		}
	}

	r.cancelF()
}

func (r *REPL) exec() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.executing || r.finished {
		return nil
	}

	r.cmd.Stdin = os.Stdin
	r.cmd.Stderr = nil
	r.cmd.Stdout = r.buf

	go func() {
		for r.executing {
			if r.buf.Len() > 0 {
				eraseCurrentLine()
				writeRaw(r.buf.Bytes(), true)
				r.buf.Reset()
			}

			time.Sleep(replCmdSleepInterval)
		}

		r.shutdown()
	}()

	r.executing = true
	err := r.cmd.Run()
	r.executing = false
	r.finished = true

	if err != nil {
		return err
	}

	return nil
}

func (r *REPL) shutdown() {
	if atomic.AddUint32(&r.closing, 1) == 1 {
		common.Log.Tracef("shutting down")
		r.cancelF()
	}
}

func (r *REPL) shuttingDown() bool {
	return (atomic.LoadUint32(&r.closing) > 0)
}
