package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/finove/goutils/logger"
	"github.com/spf13/cobra"
)

var (
	ncUseUDP     bool
	ncListenPort uint16
	ncRemote     string
)

var ncCmd = &cobra.Command{
	Use:     "nc",
	Short:   "same as nc",
	Version: "0.0.1",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			ncRemote = fmt.Sprintf("%s:%s", args[0], args[1])
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var conn net.Conn
		if ncListenPort > 0 {
			if ncUseUDP == true {
				var uaddr *net.UDPAddr
				if uaddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", ncListenPort)); err != nil {
					logger.Error("parse udp address fail:%v", err)
					return
				}
				if conn, err = net.ListenUDP("udp", uaddr); err != nil {
					logger.Error("listen udp address %s fail:%v", uaddr.String(), err)
					return
				}
			} else {
				var laddr *net.TCPAddr
				var ln *net.TCPListener
				if laddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", ncListenPort)); err != nil {
					logger.Error("parse tcp address fail:%v", err)
					return
				}
				if ln, err = net.ListenTCP("tcp", laddr); err != nil {
					logger.Error("listen tcp address %s fail:%v", laddr.String(), err)
					return
				}
				conn, err = ln.Accept()
				if err != nil {
					logger.Error("accept tcp connection fail:%v", err)
					return
				}
			}
		} else if ncRemote != "" {
			if ncUseUDP == true {
				conn, err = net.DialTimeout("udp", ncRemote, 10*time.Second)
			} else {
				conn, err = net.DialTimeout("tcp", ncRemote, 10*time.Second)
			}
			if err != nil {
				logger.Error("connect to %s as udp %v fail:%v", ncRemote, ncUseUDP, err)
				return
			}
		} else {
			logger.Error("missing paramater")
			return
		}
		go processInput(conn)
		processRemote(nil, conn)
	},
}

func processRemote(ctx context.Context, conn net.Conn) {
	var err error
	var n int
	var buff = make([]byte, 2048)
	for {
		n, err = conn.Read(buff)
		if err != nil {
			break
		} else if n <= 0 {
			continue
		}
		fmt.Printf("%s", string(buff[:n]))
	}
	return
}

func processInput(conn net.Conn) {
	var err error
	var n int
	var buff = make([]byte, 2048)
	for {
		n, err = os.Stdin.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				conn.Close()
				break
			}
			continue
		} else if n <= 0 {
			continue
		}
		n, err = conn.Write(buff[:n])
		if err != nil {
			logger.Error("send to remote fail:%v", err)
			break
		}
	}
	return
}

func transUDPAddr(addr string) (uaddr *net.UDPAddr) {
	uaddr, _ = net.ResolveUDPAddr("udp", addr)
	return
}

func transTCPAddr(addr string) (uaddr *net.TCPAddr) {
	uaddr, _ = net.ResolveTCPAddr("tcp", addr)
	return
}

func init() {
	rootCmd.AddCommand(ncCmd)
	ncCmd.Flags().BoolVarP(&ncUseUDP, "udp", "u", false, "use udp")
	ncCmd.Flags().Uint16VarP(&ncListenPort, "listen", "l", 0, "listen port")
}
