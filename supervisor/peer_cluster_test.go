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
})
