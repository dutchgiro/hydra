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

type PeerCluster struct {
	current int
	peers   []Peer
}

func NewPeerCluster(peers []Peer) *PeerCluster {
	return &PeerCluster{
		current: -1,
		peers:   peers,
	}
}

func NewPeerClusterFromNodes(nodes Nodes) {
	return &PeerCluster{
		current: -1,
		peers:   parseRawPeerCluster(nodes),
	}

}

func (p *PeerCluster) HasNext() bool {
	if p.current >= len(p.peers)-1 {
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

func (p *PeerCluster) Reset() {
	p.current = -1
}

func (p *PeerCluster) SetPeers(peers []Peer) {
	p.peers = peers
}

func (p *PeerCluster) parseRawPeerCluster(nodes Nodes) []Peer {
	peerCluster := []Peer{}
	for _, rawPeer := range nodes {
		peerCluster = append(peerCluster, p.parseRawPeer(rawPeer))
	}
	return peerCluster
}

func (p *PeerCluster) parseRawPeer(rawPeer *Node) Peer {
	peer := Peer{PeerAddr: rawPeer.Key}
	for _, attr := range rawPeer.Nodes {
		switch attr.Key {
		case AddrKey:
			peer.Addr = attr.Value
		case StateKey:
			peer.State = attr.Value
		}
	}
	return peer
}
