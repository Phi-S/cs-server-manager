package main

import (
	"cs-server-manager/event"
	"cs-server-manager/logwrt"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"golang.org/x/net/websocket"
)

type OutgoingWebsocketMessage struct {
	Type    string `json:"type"`
	Message any    `json:"message"`
}

type IncomingWebSocketMessage struct {
	clientConnection *websocket.Conn
	message          string
}

type WebSocketServer struct {
	connectionLock sync.Mutex
	connections    []*websocket.Conn

	OnIncomingMessageEvent    event.InstanceWithData[IncomingWebSocketMessage]
	OnNewClientConnectedEvent event.InstanceWithData[*websocket.Conn]
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		connections: make([]*websocket.Conn, 0),
	}
}

func (s *WebSocketServer) handleWs(con *websocket.Conn) {
	slog.Debug("web socket client connected", "address", con.RemoteAddr())

	s.connectionLock.Lock()
	s.connections = append(s.connections, con)
	s.connectionLock.Unlock()

	s.OnNewClientConnectedEvent.Trigger(con)

	err := s.read(con)
	if err == io.EOF {
		slog.Info("client closed the web socket connection", "address", con.RemoteAddr())
	} else {
		slog.Error("failed to read from client connection", "error", err)
	}
}

func (s *WebSocketServer) read(con *websocket.Conn) error {
	const errorThreshold = 5

	buf := make([]byte, 1024)
	errors := make([]error, 0)
	for {
		n, err := con.Read(buf)
		if err != nil {
			if err == io.EOF {
				return err
			}

			errors = append(errors, err)
			slog.Warn("error reading client message", "error-count", len(errors), "address", con.RemoteAddr(), "error", err)

			if len(errors) >= errorThreshold {
				return fmt.Errorf("%v errors in a row occurred while tying to read from client. Errors %v", len(errors), errors)
			}

			continue
		}
		if len(errors) > 1 {
			errors = make([]error, 0)
		}

		msgBytes := buf[:n]

		msg := IncomingWebSocketMessage{
			clientConnection: con,
			message:          string(msgBytes),
		}
		s.OnIncomingMessageEvent.Trigger(msg)

	}
}

func (s *WebSocketServer) broadcast(msg []byte) error {
	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()

	errorsLock := sync.Mutex{}
	errors := make([]error, 0)
	for _, con := range s.connections {
		go func() {
			if _, err := con.Write(msg); err != nil {
				errorsLock.Lock()
				errors = append(
					errors,
					fmt.Errorf("failed to send message to client. message: %v | address: %v | error %v",
						string(msg), con.RemoteAddr(), err),
				)
				errorsLock.Unlock()
			}
		}()
	}

	if len(errors) > 0 {
		return fmt.Errorf("%v errors occurred. Errors: %v", len(errors), errors)
	}

	return nil
}

func (s *WebSocketServer) Broadcast(messageType string, jsonMessage any) error {
	message := OutgoingWebsocketMessage{
		Type:    messageType,
		Message: jsonMessage,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.broadcast(messageBytes)
}

func (s *WebSocketServer) BroadcastLogMessage(logEntry logwrt.LogEntry) error {
	return s.Broadcast("log", logEntry)
}
