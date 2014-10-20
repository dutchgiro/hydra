package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"time"
)

var _ = Describe("PeersMonitor", func() {
	const (
		peerAddrItself                   string        = "127.0.0.1:7701"
		stateKeyTTL                      uint64        = 5
		durationBetweenPublicationsState time.Duration = time.Duration(500) * time.Millisecond
	)

	var (
		ch             chan []Peer
		mockCtrl       *gomock.Controller
		mockEtcdClient *mock.MockEtcdRequester
		peersMonitor   *PeersMonitor
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
		peersMonitor = NewPeersMonitor(mockEtcdClient)
		ch = make(chan []Peer)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	expectedCluster := []Peer{
		Peer{
			Addr:     "98.245.153.111:4001",
			PeerAddr: "98.245.153.111:7001",
			State:    PeerStateEnabled,
		},
		Peer{
			Addr:     "98.245.153.112:4001",
			PeerAddr: "98.245.153.112:7001",
			State:    PeerStateEnabled,
		},
	}
	var modifiedIndex uint64 = 3
	successResponse := &Response{
		Action: "",
		Node: &Node{
			Key:        ClusterKey,
			Value:      "",
			Dir:        true,
			Expiration: nil,
			TTL:        3,
			Nodes: []*Node{
				&Node{
					Key:        expectedCluster[0].Addr,
					Value:      "",
					Dir:        true,
					Expiration: nil,
					TTL:        3,
					Nodes: []*Node{
						&Node{
							Key:           PeerAddrKey,
							Value:         expectedCluster[0].PeerAddr,
							Dir:           false,
							Expiration:    nil,
							TTL:           3,
							Nodes:         nil,
							ModifiedIndex: modifiedIndex,
							CreatedIndex:  modifiedIndex - 1,
						},
						&Node{
							Key:           StateKey,
							Value:         expectedCluster[0].State,
							Dir:           false,
							Expiration:    nil,
							TTL:           3,
							Nodes:         nil,
							ModifiedIndex: modifiedIndex,
							CreatedIndex:  modifiedIndex - 1,
						},
					},
					ModifiedIndex: modifiedIndex,
					CreatedIndex:  modifiedIndex - 1,
				},
				&Node{
					Key:        expectedCluster[1].Addr,
					Value:      "",
					Dir:        true,
					Expiration: nil,
					TTL:        3,
					Nodes: []*Node{
						&Node{
							Key:           PeerAddrKey,
							Value:         expectedCluster[1].PeerAddr,
							Dir:           false,
							Expiration:    nil,
							TTL:           3,
							Nodes:         nil,
							ModifiedIndex: modifiedIndex,
							CreatedIndex:  modifiedIndex - 1,
						},
						&Node{
							Key:           StateKey,
							Value:         expectedCluster[1].State,
							Dir:           false,
							Expiration:    nil,
							TTL:           3,
							Nodes:         nil,
							ModifiedIndex: modifiedIndex,
							CreatedIndex:  modifiedIndex - 1,
						},
					},
					ModifiedIndex: modifiedIndex,
					CreatedIndex:  modifiedIndex - 1,
				},
			},
			ModifiedIndex: modifiedIndex,
			CreatedIndex:  modifiedIndex - 1,
		},
		PrevNode:  nil,
		EtcdIndex: 0,
		RaftIndex: 0,
		RaftTerm:  0,
	}
	var requestClusterInterval float64 = DefaultRequestClusterInterval.Seconds()

	Describe("Run", func() {
		It("should request the cluster directory periodically", func() {
			c1 := mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
				Return(successResponse, nil).
				Times(1)
			mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
				Return(successResponse, nil).
				AnyTimes().After(c1)

			go func() {
				peersMonitor.Run(ch)
			}()
			time.Sleep(DefaultRequestClusterInterval - time.Duration(1)*time.Second)
		})
		Context("when no change is detected in the cluster", func() {
			It("should not emit any cluster", func() {
				peersMonitor.Peers = expectedCluster
				c1 := mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
					Return(successResponse, nil).
					Times(1)
				mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
					Return(successResponse, nil).
					AnyTimes().After(c1)

				go func() {
					peersMonitor.Run(ch)
				}()
				Consistently(ch, requestClusterInterval*2.2, requestClusterInterval).ShouldNot(Receive())
			})
		})
		Context("when a change is detected in the cluster", func() {
			It("should emit the new cluster", func() {
				peersMonitor.Peers = expectedCluster[:1]
				mockEtcdClient.EXPECT().Get(gomock.Eq("cluster"), gomock.Eq(true), gomock.Eq(true)).
					Return(successResponse, nil).
					Times(1)

				go func() {
					peersMonitor.Run(ch)
				}()
				Eventually(ch, requestClusterInterval-1.0,
					requestClusterInterval/3).Should(Receive())
				Expect(peersMonitor.Peers).To(Equal(expectedCluster))
			})
		})
	})
})
