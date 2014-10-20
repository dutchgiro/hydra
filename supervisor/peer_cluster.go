package supervisor

import (
	"errors"
)

type Peer struct {
	Addr     string
	PeerAddr string
	State    string
}

type IterableCluster interface {
	HasNext() bool
	Next() (Peer, error)
	SetPeers(peers []Peer)
}

var peerNotFoundError error = errors.New("Peer not found")

type PeerCluster struct {
	current int
	Peers   []Peer
}

func NewPeerCluster(peers []Peer) *PeerCluster {
	return &PeerCluster{
		current: -1,
		Peers:   peers,
	}
}

func NewPeerClusterFromNodes(nodes Nodes) *PeerCluster {
	return &PeerCluster{
		current: -1,
		Peers:   parseRawPeerCluster(nodes),
	}
}

func (p *PeerCluster) GetEnabledPeers() []*Peer {
	enabledPeers := []*Peer{}
	for _, peer := range p.Peers {
		if peer.State == PeerStateEnabled {
			enabledPeers = append(enabledPeers, &peer)
		}
	}
	return enabledPeers
}

func (p *PeerCluster) GetPeerByAddr(addr string) (*Peer, error) {
	for _, peer := range p.Peers {
		if peer.Addr == addr {
			return &peer, nil
		}
	}
	return nil, peerNotFoundError
}

func (p *PeerCluster) GetPeerPosition(addr string) (int, error) {
	for i := 0; i < len(p.Peers); i++ {
		if p.Peers[i].Addr == addr {
			return i, nil
		}
	}
	return -1, peerNotFoundError
}

func (p *PeerCluster) HasNext() bool {
	if p.current >= len(p.Peers)-1 {
		return false
	}
	return true
}

func (p *PeerCluster) Next() (Peer, error) {
	p.current++
	if p.current >= len(p.Peers) {
		return Peer{}, errors.New("No such element")
	}
	return p.Peers[p.current], nil
}

func (p *PeerCluster) Reset() {
	p.current = -1
}

func parseRawPeerCluster(nodes Nodes) []Peer {
	peerCluster := []Peer{}
	for _, rawPeer := range nodes {
		peerCluster = append(peerCluster, parseRawPeer(rawPeer))
	}
	return peerCluster
}

func parseRawPeer(rawPeer *Node) Peer {
	peer := Peer{Addr: rawPeer.Key}
	for _, attr := range rawPeer.Nodes {
		switch attr.Key {
		case PeerAddrKey:
			peer.PeerAddr = attr.Value
		case StateKey:
			peer.State = attr.Value
		}
	}
	return peer
}
