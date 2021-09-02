package pool

import (
	"natpass/code/client/global"
	"natpass/code/network"
	"time"
)

// SendConnectReq send connect request message
func (conn *Conn) SendConnectReq(id string, cfg global.Tunnel) {
	tp := network.ConnectRequest_tcp
	if cfg.Type != "tcp" {
		tp = network.ConnectRequest_udp
	}
	var msg network.Msg
	msg.To = cfg.Target
	msg.XType = network.Msg_connect_req
	msg.Payload = &network.Msg_Creq{
		Creq: &network.ConnectRequest{
			Id:    id,
			Name:  cfg.Name,
			XType: tp,
			Addr:  cfg.RemoteAddr,
			Port:  uint32(cfg.RemotePort),
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendConnectError send connect error response message
func (conn *Conn) SendConnectError(to string, toIdx uint32, id, info string) {
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_connect_rep
	msg.Payload = &network.Msg_Crep{
		Crep: &network.ConnectResponse{
			Id:  id,
			Ok:  false,
			Msg: info,
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendConnectOK send connect success response message
func (conn *Conn) SendConnectOK(to string, toIdx uint32, id string) {
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_connect_rep
	msg.Payload = &network.Msg_Crep{
		Crep: &network.ConnectResponse{
			Id: id,
			Ok: true,
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendDisconnect send disconnect message
func (conn *Conn) SendDisconnect(to string, toIdx uint32, id string) {
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_disconnect
	msg.Payload = &network.Msg_XDisconnect{
		XDisconnect: &network.Disconnect{
			Id: id,
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendData send forward data
func (conn *Conn) SendData(to string, toIdx uint32, id string, data []byte) {
	dup := func(data []byte) []byte {
		ret := make([]byte, len(data))
		copy(ret, data)
		return ret
	}
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_forward
	msg.Payload = &network.Msg_XData{
		XData: &network.Data{
			Lid:  id,
			Data: dup(data),
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendKeepalive send keepalive message
func (conn *Conn) SendKeepalive() {
	var msg network.Msg
	msg.To = "server"
	msg.XType = network.Msg_keepalive
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}