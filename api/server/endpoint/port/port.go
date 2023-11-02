package port

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/xiayesuifeng/gopanel/api/server/router"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Port struct {
}

func (p *Port) Name() string {
	return "port"
}

func (p *Port) Run(r router.Router) {
	r.GET("/:port/available", p.Available)
}

func (p *Port) Available(ctx *router.Context) error {
	port, err := strconv.Atoi(ctx.Param("port"))
	if err != nil {
		return ctx.Error(http.StatusBadRequest, err.Error())
	}

	tcpAvailable := false
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(port)), time.Second*3)
	if err != nil {
		tcpAvailable = true
	} else {
		conn.Close()
	}

	udpAvailable := true
	conn, err = net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		udpAvailable = false
	} else {
		conn.Close()
	}

	network := "unknown"
	if tcpAvailable && !udpAvailable {
		network = "tcp"
	} else if !tcpAvailable && udpAvailable {
		network = "udp"
	} else if tcpAvailable && udpAvailable {
		network = "both"
	}

	return ctx.JSON(gin.H{
		"available": tcpAvailable || udpAvailable,
		"network":   network,
	})
}
