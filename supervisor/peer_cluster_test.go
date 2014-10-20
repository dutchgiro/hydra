package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = FDescribe("PeerCluster", func() {
	var (
		initPeers   []Peer
		peerCluster *PeerCluster
	)

	BeforeEach(func() {
		initPeers = []Peer{
			Peer{
				Addr:     "98.245.153.111:4001",
				PeerAddr: "98.245.153.111:7001",
				State:    "",
			},
			Peer{
				Addr:     "98.245.153.112:4001",
				PeerAddr: "98.245.153.112:7001",
				State:    "",
			},
		}
		peerCluster = NewPeerCluster(initPeers)
	})

	Describe("GetEnabledPeers", func() {
		It("should return the enabled peers only", func() {
			peerCluster.Peers[0].State = ""
			peerCluster.Peers[1].State = PeerStateEnabled
			Expect(peerCluster.Peers).To(HaveLen(2))
			enabledPeers := peerCluster.GetEnabledPeers()
			Expect(enabledPeers).To(HaveLen(1))
			Expect(enabledPeers[0].State).To(Equal(PeerStateEnabled))
		})
	})

	Describe("GetPeerByAddr", func() {
		Context("when peer doesn't exist", func() {
			It("should return an error", func() {
				const nonexistentPeerAddr string = "98.245.153.113:4001"
				foundPeer, err := peerCluster.GetPeerByAddr(nonexistentPeerAddr)
				Expect(err).To(HaveOccurred())
				Expect(foundPeer).To(Equal(Peer{}))
			})
		})
		Context("when peer exists", func() {
			It("should return the associated peer", func() {
				foundPeer, err := peerCluster.GetPeerByAddr(peerCluster.Peers[1].Addr)
				Expect(err).ToNot(HaveOccurred())
				Expect(foundPeer).To(Equal(peerCluster.Peers[1]))
			})
		})
	})

	Describe("GetPeerPosition", func() {
		Context("when peer exists", func() {
			It("should return the enabled peers only", func() {
				const expectedPosition int = 1
				pos, err := peerCluster.GetPeerPosition(peerCluster.Peers[expectedPosition].Addr)
				Expect(err).ToNot(HaveOccurred())
				Expect(pos).To(Equal(expectedPosition))
			})
		})
		Context("when peer exists", func() {
			It("should return the enabled peers only", func() {
				_, err := peerCluster.GetPeerPosition("127.0.0.1:4001")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Next", func() {
		It("should return the next value until end the collection", func() {
			for i := 0; i < len(initPeers); i++ {
				peer, err := peerCluster.Next()
				Expect(peer).To(Equal(initPeers[i]))
				Expect(err).NotTo(HaveOccurred())
			}
			peer, err := peerCluster.Next()
			Expect(peer).To(Equal(Peer{}))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("HasNext", func() {
		It("should return true until all peers have been requested", func() {
			for i := 0; i < len(initPeers); i++ {
				Expect(peerCluster.HasNext()).To(BeTrue())
				_, err := peerCluster.Next()
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(peerCluster.HasNext()).To(BeFalse())
		})
	})

	Describe("Reset", func() {
		It("should return the iterator to the beginning of the collection", func() {
			_, _ = peerCluster.Next()
			_, _ = peerCluster.Next()
			peerCluster.Reset()
			for i := 0; i < len(initPeers); i++ {
				Expect(peerCluster.HasNext()).To(BeTrue())
				_, err := peerCluster.Next()
				Expect(err).ToNot(HaveOccurred())
			}
			Expect(peerCluster.HasNext()).To(BeFalse())
		})
	})
})
