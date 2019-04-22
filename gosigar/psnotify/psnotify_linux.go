// Copyright (c) 2012 VMware, Inc.

// Go interface to the Linux netlink process connector.
// See Documentation/connector/connector.txt in the linux kernel source tree.
package psnotify

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ErwinYy/origin-elastic/gosigar/sys"
	"os"
	"syscall"
)

const (
	// internal flags (from <linux/connector.h>)
	_CN_IDX_PROC = 0x1
	_CN_VAL_PROC = 0x1

	// internal flags (from <linux/cn_proc.h>)
	_PROC_CN_MCAST_LISTEN = 1
	_PROC_CN_MCAST_IGNORE = 2

	// Flags (from <linux/cn_proc.h>)
	PROC_EVENT_FORK = 0x00000001 // fork() events
	PROC_EVENT_EXEC = 0x00000002 // exec() events
	PROC_EVENT_EXIT = 0x80000000 // exit() events

	// Watch for all process events
	PROC_EVENT_ALL = PROC_EVENT_FORK | PROC_EVENT_EXEC | PROC_EVENT_EXIT
)

var (
	byteOrder = sys.GetEndian()
)

// linux/connector.h: struct cb_id
type cbId struct {
	Idx uint32
	Val uint32
}

// linux/connector.h: struct cb_msg
type cnMsg struct {
	Id    cbId
	Seq   uint32
	Ack   uint32
	Len   uint16
	Flags uint16
}

// linux/cn_proc.h: struct proc_event.{what,cpu,timestamp_ns}
type procEventHeader struct {
	What      uint32
	Cpu       uint32
	Timestamp uint64
}

// linux/cn_proc.h: struct proc_event.fork
type forkProcEvent struct {
	ParentPid  uint32
	ParentTgid uint32
	ChildPid   uint32
	ChildTgid  uint32
}

// linux/cn_proc.h: struct proc_event.exec
type execProcEvent struct {
	ProcessPid  uint32
	ProcessTgid uint32
}

// linux/cn_proc.h: struct proc_event.exit
type exitProcEvent struct {
	ProcessPid  uint32
	ProcessTgid uint32
	ExitCode    uint32
	ExitSignal  uint32
}

// standard netlink header + connector header
type netlinkProcMessage struct {
	Header syscall.NlMsghdr
	Data   cnMsg
}

type netlinkListener struct {
	addr *syscall.SockaddrNetlink // Netlink socket address
	sock int                      // The syscall.Socket() file descriptor
	seq  uint32                   // struct cn_msg.seq
}

func (this *netlinkListener) getSock() int {
	return this.sock
}

// Initialize linux implementation of the eventListener interface
func createListener(rmemMax int) (eventListener, error) {
	listener := &netlinkListener{}
	err := listener.bind(rmemMax)
	return listener, err
}

// noop on linux
func (w *Watcher) unregister(pid int) error {
	return nil
}

// noop on linux
func (w *Watcher) register(pid int, flags uint32) error {
	return nil
}

func (w *Watcher) ProcessRecvBuf() {
	var msgs []syscall.NetlinkMessage
	var err error
	var bufInfo = make([]byte, syscall.Getpagesize())
	var chanOK bool
	
	for {
		select {
		case bufInfo, chanOK = <-w.BufList:
			if !chanOK || w.isDone() {
				return
			}

			//fmt.Printf("******* ProcessRecvBuf len(bufInfo)): %v \n", len(string(bufInfo)))

			msgs, err = syscall.ParseNetlinkMessage(bufInfo)
			if err != nil {
				w.Error <- fmt.Errorf("error processRecvBuf got %v ", err)
				bufInfo = bufInfo[:0]
				continue
			}
			bufInfo = bufInfo[:0]

			//fmt.Printf("******* ProcessRecvBuf got msgs: %+v \n", msgs)

			for _, m := range msgs {
				if m.Header.Type == syscall.NLMSG_DONE {
					w.handleEvent(m.Data)
				}
			}
		}
	}
}

// Read events from the netlink socket
func (w *Watcher) readEvents() {
	var bufRecv []byte
	bufRecv = make([]byte, syscall.Getpagesize())

	var listener, _ = w.listener.(*netlinkListener)

	for {
		if w.isDone() {
			return
		}

		nr, _, err := syscall.Recvfrom(listener.sock, bufRecv, 0)

		//fmt.Printf("******* readEvents got len: %v \n", nr)

		if err != nil {
			w.Error <- err
			continue
		}

		if nr < syscall.NLMSG_HDRLEN {
			w.Error <- syscall.EINVAL
			continue
		}

		//fmt.Printf("******* readEvents Recvfrom got nr = %d \n",  nr)

		//var b = make([]byte, nr)
		//copy(b, bufRecv[:nr])
		select {
		case w.BufList <- bufRecv[:nr]:
		default:
		}
	}
}

// Internal helper to check if pid && event is being watched
func (w *Watcher) isWatching(pid int, event uint32) bool {
	//2018-12-04 返回全部 pid 事件
	return true
	if watch, ok := w.watches[pid]; ok {
		return (watch.flags & event) == event
	}
	return false
}

// Dispatch events from the netlink socket to the Event channels.
// Unlike bsd kqueue, netlink receives events for all pids,
// so we apply filtering based on the watch table via isWatching()
func (w *Watcher) handleEvent(data []byte) {
	var err error
	buf := bytes.NewBuffer(data)
	msg := AcquireCnMsgPool()
	hdr := AcquireProcEvent()

	err = binary.Read(buf, byteOrder, msg)
	if err != nil {
		return
	}
	err = binary.Read(buf, byteOrder, hdr)

	if err != nil {
		return
	}
	switch hdr.What {
	case PROC_EVENT_FORK:
		event := AcquireForkProcEvent()
		err = binary.Read(buf, byteOrder, event)
		if err != nil {
			return
		}
		ppid := int(event.ParentTgid)
		//pid := int(event.ChildTgid)
		pid := int(event.ChildPid)

		//fmt.Printf("******* PROC_EVENT_FORK got event: %+v \n", event)

		//if w.isWatching(ppid, PROC_EVENT_EXEC) {
			// follow forks
			//watch, _ := w.watches[ppid]
			//w.Watch(pid, watch.flags)
		//}

		if w.isWatching(ppid, PROC_EVENT_FORK) {
			//var proceventFork = AcquireProcEventFork()
			//proceventFork.ParentPid = ppid
			//proceventFork.ChildPid = pid

			//fmt.Printf("******* PROC_EVENT_FORK got proceventFork: %+v \n", proceventFork)

			select {
			case w.Fork <- pid: //proceventFork.ChildPid:
			default:
			}
			//ReleaseProcEventFork(proceventFork)
		}
		ReleaseForkProcEvent(event)
	case PROC_EVENT_EXEC:
		event := AcquireExecProcEvent()
		err = binary.Read(buf, byteOrder, event)
		if err != nil {
			return
		}
		pid := int(event.ProcessTgid)

		//fmt.Printf("******* PROC_EVENT_EXEC got event: %+v \n", event)

		if w.isWatching(pid, PROC_EVENT_EXEC) {
			//var proceventExec = AcquireProcEventExec()
			//proceventExec.Pid = pid

			//fmt.Printf("******* PROC_EVENT_EXEC got proceventExec: %+v \n", proceventExec)

			select {
			case w.Exec <- pid: //proceventExec.Pid:
			default:
			}
			//ReleaseProcEventExec(proceventExec)
		}
		ReleaseExecProcEvent(event)
	case PROC_EVENT_EXIT:
		event := AcquireExitProcEvent()
		err = binary.Read(buf, byteOrder, event)
		if err != nil {
			return
		}
		pid := int(event.ProcessTgid)

		//fmt.Printf("******* PROC_EVENT_EXIT got event: %+v \n", event)

		if w.isWatching(pid, PROC_EVENT_EXIT) {
			w.RemoveWatch(pid)
			//var proceventExit = AcquireProcEventExit()
			//proceventExit.Pid = pid

			//fmt.Printf("******* PROC_EVENT_EXIT got proceventExit: %+v \n", proceventExit)

			select {
			case w.Exit <- pid://proceventExit.Pid:
			default:
			}
			//ReleaseProcEventExit(proceventExit)
		}
		ReleaseExitProcEvent(event)
	}

	//释放 pool
	ReleaseCnMsgPool(msg)
	ReleaseProcEvent(hdr)
}

// Bind our netlink socket and
// send a listen control message to the connector driver.
func (listener *netlinkListener) bind(rmemMax int) error {
	sock, err := syscall.Socket(
		syscall.AF_NETLINK,
		syscall.SOCK_DGRAM,
		syscall.NETLINK_CONNECTOR)

	if err != nil {
		return err
	}

	listener.sock = sock
	listener.addr = &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Groups: _CN_IDX_PROC,
	}

	//if rmemMax != 0 {
	//	err = syscall.SetsockoptInt(listener.sock, syscall.SOL_SOCKET, syscall.SO_RCVBUF, rmemMax)
	//	if err != nil {
	//		return err
	//	}
	//}

	err = syscall.Bind(listener.sock, listener.addr)

	if err != nil {
		return err
	}

	return listener.send(_PROC_CN_MCAST_LISTEN)
}

// Send an ignore control message to the connector driver
// and close our netlink socket.
func (listener *netlinkListener) close() error {
	err := listener.send(_PROC_CN_MCAST_IGNORE)
	syscall.Close(listener.sock)
	return err
}

// Generic method for sending control messages to the connector
// driver; where op is one of PROC_CN_MCAST_{LISTEN,IGNORE}
func (listener *netlinkListener) send(op uint32) error {
	listener.seq++
	pr := &netlinkProcMessage{}
	plen := binary.Size(pr.Data) + binary.Size(op)
	pr.Header.Len = syscall.NLMSG_HDRLEN + uint32(plen)
	pr.Header.Type = uint16(syscall.NLMSG_DONE)
	pr.Header.Flags = 0
	pr.Header.Seq = listener.seq
	pr.Header.Pid = uint32(os.Getpid())

	pr.Data.Id.Idx = _CN_IDX_PROC
	pr.Data.Id.Val = _CN_VAL_PROC

	pr.Data.Len = uint16(binary.Size(op))

	buf := bytes.NewBuffer(make([]byte, 0, pr.Header.Len))
	binary.Write(buf, byteOrder, pr)
	binary.Write(buf, byteOrder, op)

	return syscall.Sendto(listener.sock, buf.Bytes(), 0, listener.addr)
}
