package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	// . "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	// "errors"
	"net/http"
	"time"
)

var _ = FDescribe("ClusterInspector", func() {
	const (
		etcdAddrItself string = "127.0.0.1:7401"
		peerAddrItself string = "127.0.0.1:7701"
	)

	var (
		// ch               chan PeerCluster
		mockCtrl         *gomock.Controller
		mockEtcdClient   *mock.MockEtcdRequester
		mockPeersMonitor *mock.MockFolderMonitor
		clusterInspector *ClusterInspector
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
		mockPeersMonitor = mock.NewMockFolderMonitor(mockCtrl)
		clusterInspector = NewClusterInspector(etcdAddrItself, peerAddrItself, []string{})
		clusterInspector.EtcdClient = mockEtcdClient
		clusterInspector.PeersMonitor = mockPeersMonitor
		// ch = make(chan StateControllerState)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	successGetLeaderResponse := &RawResponse{
		StatusCode: http.StatusOK,
		Body:       []byte("http://" + peerAddrItself),
		Header:     nil,
	}

	Describe("Run", func() {
		It("should run a Peers Monitor", func() {
			mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)

			go func() {
				clusterInspector.Run()
			}()
			time.Sleep(time.Duration(1) * time.Second)
		})
		It("should search for a peer to connect", func() {
			mockPeersMonitor.EXPECT().Run(gomock.Any()).AnyTimes()
			c1 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).Return(successGetLeaderResponse, nil).Times(1)
			mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).Return(successGetLeaderResponse, nil).AnyTimes().After(c1)

			go func() {
				clusterInspector.Run()
			}()
			time.Sleep(time.Duration(3) * time.Second)
		})
	})
})

// **************************************************************************************************************

// Context("when node find a foreign peer to connect", func() {
// 	// It("should check if this peer is leader", func() {

// 	// })
// 	Context("when reachable peer is not leader", func() {
// 		It("should continue searching a peer to connect", func() {
// 			foreignPeerAddr := "98.245.153.111:7701"
// 			// res := *RawResponse{
// 			// 	StatusCode: http.StatusOK,
// 			// 	Body:       []byte("http://" + peerAddrItself),
// 			// 	Header:     nil,
// 			// }
// 			c1 := mockEtcdClient.EXPECT().WithMachineAddr("http://" + foreignPeerAddr).Times(1)
// 			c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq("leader")).Return(successGetLeaderResponse, nil).Times(1).After(c1)
// 			mockEtcdClient.EXPECT().WithMachineAddr(gomock.Not("http://" + peerAddrItself)).AnyTimes().After(c2)
// 			mockEtcdClient.EXPECT().BaseGet(gomock.Eq("leader")).Return(successGetLeaderResponse, nil).AnyTimes().After(c2)

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 			time.Sleep(time.Duration(2) * time.Second)
// 		})
// 	})
// 	Context("when reachable peer is leader", func() {
// 		// Next request for machines and compare with local cluster for overwrite local or remote cluster
// 		It("should get registered peers from foreign peer", func() {
// 			f1 := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).
// 				SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 			f2 := mockEtcdClient.EXPECT().Get(gomock.Any(), gomock.Any()).
// 				SetArg("http://98.245.153.111:7001", "/cluster")

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 		})
// 		Context("when the registered peers from foreign peer leader are retrieved", func() {
// 			It("should add unknown peers to foreign leader", func() {

// 			})
// 		})

// 		// TODO: priority = higher cluter size && node priority
// 		Context("when it has lower cluster priority than reachable peer", func() {
// 			It("should try connecting to reachable peer", func() {
// 				successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 				mockEtcdClient.EXPECT().GetLeader(gomock.Any()).AnyTimes().After(successfulCallToSlaveNode)
// 				// TODO: restart with attributes
// 				mockEtcdManager.EXPECT().RestartEtcdService().Times(1).After(successfulCallToSlaveNode)

// 				go func() {
// 					etcdSupervisor.Run()
// 				}()
// 			})
// 		})
// 		Context("when it has higher cluster priority than reachable peer", func() {
// 			It("should remain looking for a foreign peer to connect", func() {
// 				successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 				mockEtcdClient.EXPECT().GetLeader(gomock.Any()).After(successfulCallToSlaveNode)
// 			})
// 		})
// 	})
// 	It("should send request to connect", func() {
// 		mockEtcdClient.EXPECT().Set().Times(1)

// 		go func() {
// 			etcdSupervisor.Run()
// 		}()
// 	})
// 	Context("when request to connect is accepted", func() {
// 		It("should try to connect", func() {

// 		})
// 		Context("when connect to peer", func() {
// 			It("should change its state to slave", func() {

// 			})
// 		})
// 		Context("when connect to peer is impossible", func() {
// 			It("should remain as master prima", func() {

// 			})
// 		})
// 	})
// 	Context("when request to connect is not accepted", func() {
// 		It("should try to connect", func() {

// 		})
// 	})
// })
// 	})
// })
