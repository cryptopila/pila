package coin

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// handshakeMessage is exchanged by peers when establishing a connection.
type handshakeMessage struct {
	Protocol uint32 `json:"protocol"`
	NodeID   string `json:"node"`
}

// PerformHandshake initiates a handshake with the remote peer. It sends a
// handshakeMessage containing the local node ID and waits for the remote
// response. A protocol mismatch results in an error.
// PerformHandshake initiates a handshake with the remote peer and returns the
// remote node ID if successful.
func PerformHandshake(conn net.Conn, nodeID string) (string, error) {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	hello := handshakeMessage{Protocol: P2PProtocolVersion, NodeID: nodeID}
	if err := enc.Encode(&hello); err != nil {
		return "", err
	}

	// Set a small deadline so tests won't hang forever.
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var resp handshakeMessage
	if err := dec.Decode(&resp); err != nil {
		return "", err
	}
	// Reset the deadline so future reads are not affected.
	_ = conn.SetReadDeadline(time.Time{})

	if resp.Protocol != hello.Protocol {
		return "", fmt.Errorf("unexpected protocol %d", resp.Protocol)
	}
	return resp.NodeID, nil
}

// HandleHandshake responds to a handshake request from a peer and returns the
// remote node ID.
func HandleHandshake(conn net.Conn, nodeID string) (string, error) {
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	var req handshakeMessage
	if err := dec.Decode(&req); err != nil {
		return "", err
	}

	resp := handshakeMessage{Protocol: P2PProtocolVersion, NodeID: nodeID}
	if err := enc.Encode(&resp); err != nil {
		return "", err
	}
	return req.NodeID, nil
}
