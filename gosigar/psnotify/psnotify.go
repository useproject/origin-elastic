// Copyright (c) 2012 VMware, Inc.

// +build darwin freebsd netbsd openbsd linux

package psnotify

import (
	"errors"
)

type ProcEventFork struct {
	ParentPid int // Pid of the process that called fork()
	ChildPid  int // Child process pid created by fork()
}

type ProcEventExec struct {
	Pid int // Pid of the process that called exec()
}

type ProcEventExit struct {
	Pid int // Pid of the process that called exit()
}

type watch struct {
	flags uint32 // Saved value of Watch() flags param
}

type eventListener interface {
	close() error // Watch.Close() closes the OS specific listener
	getSock() int
}

const (
	PROCEVENT_CHAN_LENGTH     int = 1024
	PROCEVENT_BUF_CHAN_LENGTH int = 1024 * 4
)

type Watcher struct {
	listener eventListener  // OS specifics (kqueue or netlink)
	watches  map[int]*watch // Map of watched process ids
	Error    chan error     // Errors are sent on this channel
	Fork     chan int       // Fork events are sent on this channel
	Exec     chan int       // Exec events are sent on this channel
	Exit     chan int       // Exit events are sent on this channel
	BufList  chan []byte
	done     chan bool // Used to stop the readEvents() goroutine
	isClosed bool      // Set to true when Close() is first called
}

// Initialize event listener and channels
func NewWatcher(rmemMax int) (*Watcher, error) {

	listener, err := createListener(rmemMax)
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		listener: listener,
		watches:  make(map[int]*watch),
		Fork:     make(chan int, PROCEVENT_CHAN_LENGTH),
		Exec:     make(chan int, PROCEVENT_CHAN_LENGTH),
		Exit:     make(chan int, PROCEVENT_CHAN_LENGTH),
		BufList:  make(chan []byte, PROCEVENT_BUF_CHAN_LENGTH),
		Error:    make(chan error),
		done:     make(chan bool, 1),
	}

	go w.ProcessRecvBuf()
	go w.readEvents()
	return w, nil
}

// Close event channels when done message is received
func (w *Watcher) finish() {
	close(w.Fork)
	close(w.Exec)
	close(w.Exit)
	close(w.Error)
	close(w.BufList)
}

// Closes the OS specific event listener,
// removes all watches and closes all event channels.
func (w *Watcher) Close() error {
	if w.isClosed {
		return nil
	}
	w.isClosed = true

	for pid := range w.watches {
		w.RemoveWatch(pid)
	}

	w.done <- true

	w.listener.close()

	return nil
}

// Add pid to the watched process set.
// The flags param is a bitmask of process events to capture,
// must be one or more of: PROC_EVENT_FORK, PROC_EVENT_EXEC, PROC_EVENT_EXIT
func (w *Watcher) Watch(pid int, flags uint32) error {
	//if w.isClosed {
	//	return errors.New("psnotify watcher is closed")
	//}

	//监控所有 PID
	return nil

	watchEntry, found := w.watches[pid]

	if found {
		watchEntry.flags |= flags
	} else {
		if err := w.register(pid, flags); err != nil {
			return err
		}
		w.watches[pid] = &watch{flags: flags}
	}

	return nil
}

// Remove pid from the watched process set.
func (w *Watcher) RemoveWatch(pid int) error {
	//2018-12-04 监听所有进程
	return nil
	_, ok := w.watches[pid]
	if !ok {
		return errors.New("watch for pid does not exist")
	}
	delete(w.watches, pid)
	return w.unregister(pid)
}

// Internal helper to check if there is a message on the "done" channel.
// The "done" message is sent by the Close() method; when received here,
// the Watcher.finish method is called to close all channels and return
// true - in which case the caller should break from the readEvents loop.
func (w *Watcher) isDone() bool {
	var done bool
	select {
	case done = <-w.done:
		w.finish()
	default:
	}
	return done
}

func (w *Watcher) Run() {
	w.Watch(0, PROC_EVENT_ALL)
}
