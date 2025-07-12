package net

import (
	"io"
	"net"
	"time"

	"github.com/Dreamacro/clash/log"
)

// Relay copies between left and right bidirectionally.
func Relay(leftConn, rightConn net.Conn) {
	ch := make(chan error)

	go func() {
		// Wrapping to avoid using *net.TCPConn.(ReadFrom)
		// See also https://github.com/Dreamacro/clash/pull/1209
		n, err := io.Copy(WriteOnlyWriter{Writer: leftConn}, ReadOnlyReader{Reader: rightConn})
		log.Infoln("copy rightConn[%v<-%v] to leftConn[%v->%v] %d B",
			rightConn.LocalAddr(), rightConn.RemoteAddr(), leftConn.LocalAddr(), leftConn.RemoteAddr(), n)
		leftConn.SetReadDeadline(time.Now())
		ch <- err
	}()

	n, _ := io.Copy(WriteOnlyWriter{Writer: rightConn}, ReadOnlyReader{Reader: leftConn})
	log.Infoln("copy leftConn[%v<-%v] to rightConn[%v->%v] %d B",
		leftConn.LocalAddr(), leftConn.RemoteAddr(), rightConn.LocalAddr(), rightConn.RemoteAddr(), n)
	rightConn.SetReadDeadline(time.Now())
	<-ch
}
