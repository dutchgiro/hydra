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
	Merge(newPeers []Peer)
	Next() bool
	Value() Peer
}

type PeerCluster struct {
	configPeers []string
	current     int
	peers       []Peer
}

func NewPeerCluster(peers []Peer) *PeerCluster {
	return &PeerCluster{peers: peers, current: -1}
}

func (p *PeerCluster) Merge(newPeers []Peer) {
	// return p.peers[p.current]
}

func (p *PeerCluster) HasNext() bool {
	if p.current >= len(p.peers) {
		return false
	}
	return true
}

func (p *PeerCluster) Next() (Peer, error) {
	p.current++
	if p.current >= len(p.peers) {
		return Peer{}, errors.New("No such element")
	}
	return p.peers[p.current], nil
}

// func (p *PeerCluster) Value() Peer {
// 	return p.peers[p.current]
// }
