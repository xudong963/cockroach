// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package log

import "github.com/cockroachdb/cockroach/pkg/util/syncutil"

type loggerRegistry struct {
	mu struct {
		syncutil.Mutex
		loggers []*loggerT
	}
}

var registry = loggerRegistry{}

// debugLog is the logger instance for “general” logging messages,
// that is, those not going to a specialized secondary logger.
var debugLog loggerT

func init() {
	// Make debugLog known to the registry.
	// this ensures that all iterations also traverse
	// debugLog.
	registry.put(&debugLog)
}

// stderrLog is the logger where writes performed directly
// to the stderr file descriptor (such as that performed
// by the go runtime) *may* be redirected.
// NB: whether they are actually redirected is determined
// by stderrLog.redirectInternalStderrWrites().
var stderrLog = &debugLog

// len returns the number of known loggers.
func (r *loggerRegistry) len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.mu.loggers)
}

// iterate iterates over all the loggers and stops at the first error
// encountered.
func (r *loggerRegistry) iter(fn func(l *loggerT) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, l := range r.mu.loggers {
		if err := fn(l); err != nil {
			return err
		}
	}
	return nil
}

// iterLocked is like iter but it also locks each logger visited.
func (r *loggerRegistry) iterLocked(fn func(l *loggerT) error) error {
	return r.iter(func(l *loggerT) error {
		l.mu.Lock()
		defer l.mu.Unlock()
		return fn(l)
	})
}

// put adds a logger into the registry.
func (r *loggerRegistry) put(l *loggerT) {
	r.mu.Lock()
	r.mu.loggers = append(r.mu.loggers, l)
	r.mu.Unlock()
}

// del removes one logger from the registry.
func (r *loggerRegistry) del(l *loggerT) {
	// Make the registry forget about this logger. This avoids
	// stacking many secondary loggers together when there are
	// subsequent tests starting servers in the same package.
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, thatLogger := range r.mu.loggers {
		if thatLogger != l {
			continue
		}
		r.mu.loggers = append(r.mu.loggers[:i], r.mu.loggers[i+1:]...)
		return
	}
}
