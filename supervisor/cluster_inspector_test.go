package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"errors"
	"net/http"
	"time"
)

var _ = Describe("ClusterInspector", func() {
	const (
		etcdAddrItself string = "127.0.0.1:7401"
		peerAddrItself string = "127.0.0.1:7701"
	)

	var (
		ch               chan string
		mockCtrl         *gomock.Controller
		mockEtcdClient   *mock.MockEtcdRequester
		mockPeersMonitor *mock.MockFolderMonitor
		clusterInspector *ClusterInspector
	)

	var configPeers []string = []string{}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
		mockPeersMonitor = mock.NewMockFolderMonitor(mockCtrl)
		clusterInspector = NewClusterInspector(etcdAddrItself, configPeers)
		clusterInspector.EtcdClient = mockEtcdClient
		clusterInspector.PeersMonitor = mockPeersMonitor
		ch = make(chan string)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var modifiedIndex uint64 = 3
	makeGetClusterResponse := func(nodes []*Node) *Response {
		return &Response{
			Action: "",
			Node: &Node{
				Key:           ClusterKey,
				Value:         "",
				Dir:           true,
				Expiration:    nil,
				TTL:           3,
				Nodes:         nodes,
				ModifiedIndex: modifiedIndex,
				CreatedIndex:  modifiedIndex - 1,
			},
			PrevNode:  nil,
			EtcdIndex: 0,
			RaftIndex: 0,
			RaftTerm:  0,
		}
	}

	peer0 := Peer{
		Addr:     "98.245.153.111:4001",
		PeerAddr: "98.245.153.111:7001",
		State:    PeerStateEnabled,
	}
	peer1 := Peer{
		Addr:     "98.245.153.112:4001",
		PeerAddr: "98.245.153.112:7001",
		State:    PeerStateEnabled,
	}
	peer2 := Peer{
		Addr:     "98.245.153.113:4001",
		PeerAddr: "98.245.153.113:7001",
		State:    PeerStateEnabled,
	}
	expectedPeers := []Peer{peer0, peer1}

	node0 := &Node{
		Key:        peer0.Addr,
		Value:      "",
		Dir:        true,
		Expiration: nil,
		TTL:        3,
		Nodes: []*Node{
			&Node{
				Key:           PeerAddrKey,
				Value:         peer0.PeerAddr,
				Dir:           false,
				Expiration:    nil,
				TTL:           3,
				Nodes:         nil,
				ModifiedIndex: modifiedIndex,
				CreatedIndex:  modifiedIndex - 1,
			},
			&Node{
				Key:           StateKey,
				Value:         peer0.State,
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
	}
	node1 := &Node{
		Key:        peer1.Addr,
		Value:      "",
		Dir:        true,
		Expiration: nil,
		TTL:        3,
		Nodes: []*Node{
			&Node{
				Key:           PeerAddrKey,
				Value:         peer1.PeerAddr,
				Dir:           false,
				Expiration:    nil,
				TTL:           3,
				Nodes:         nil,
				ModifiedIndex: modifiedIndex,
				CreatedIndex:  modifiedIndex - 1,
			},
			&Node{
				Key:           StateKey,
				Value:         peer1.State,
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
	}
	node2 := &Node{
		Key:        peer2.Addr,
		Value:      "",
		Dir:        true,
		Expiration: nil,
		TTL:        3,
		Nodes: []*Node{
			&Node{
				Key:           PeerAddrKey,
				Value:         peer2.PeerAddr,
				Dir:           false,
				Expiration:    nil,
				TTL:           3,
				Nodes:         nil,
				ModifiedIndex: modifiedIndex,
				CreatedIndex:  modifiedIndex - 1,
			},
			&Node{
				Key:           StateKey,
				Value:         peer2.State,
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
	}

	cErr := errors.New("Unreachable peer")

	Describe("Run", func() {
		It("should run a Peers Monitor", func() {
			mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)

			go func() {
				clusterInspector.Run(ch)
			}()
			time.Sleep(time.Duration(1) * time.Second)
		})
		Context("when peers monitor emit new cluster", func() {
			It("should update the peer cluster", func() {
				var peersMonitorChannel chan []Peer
				mockPeersMonitor.EXPECT().Run(gomock.Any()).Do(func(ch chan []Peer) {
					peersMonitorChannel = ch
				}).Times(1)
				mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).
					AnyTimes()
				mockEtcdClient.EXPECT().BaseGet(gomock.Any()).
					Return(nil, cErr).AnyTimes()
				go func() {
					clusterInspector.Run(ch)
				}()
				Expect(clusterInspector.PeerCluster.Peers).To(Equal([]Peer{}))
				time.Sleep(time.Duration(200) * time.Millisecond)
				peersMonitorChannel <- expectedPeers
				Expect(clusterInspector.PeerCluster.Peers).To(Equal(expectedPeers))
			})
		})
		It("should search for a peer to connect", func() {
			clusterInspector.PeerCluster.Peers = expectedPeers
			mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)
			mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).AnyTimes()
			c1 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
				Return(nil, cErr).Times(1)
			mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
				Return(nil, cErr).AnyTimes().After(c1)

			go func() {
				clusterInspector.Run(ch)
			}()
			time.Sleep(time.Duration(1) * time.Second)
		})
		Context("when none peer is reachable", func() {
			It("should consult for all peers and repeat it", func() {
				clusterInspector.PeerCluster.Peers = expectedPeers
				mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)

				addr0 := AddHttpProtocol(expectedPeers[0].Addr)
				addr1 := AddHttpProtocol(expectedPeers[1].Addr)
				c1 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addr0)).
					Times(1)
				c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
					Return(nil, cErr).Times(1).After(c1)
				c3 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addr1)).
					Times(1).After(c2)
				c4 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
					Return(nil, cErr).Times(1).After(c3)
				c5 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addr0)).
					Times(1).After(c4)
				c6 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
					Return(nil, cErr).Times(1).After(c5)
				mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).
					AnyTimes().After(c6)
				mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
					Return(nil, cErr).AnyTimes().After(c6)

				go func() {
					clusterInspector.Run(ch)
				}()
				time.Sleep(time.Duration(1) * time.Second)
			})
		})
		Context("when find a foreign peer to connect", func() {
			Context("when reachable peer is not leader", func() {
				It("should continue searching a peer to connect", func() {
					clusterInspector.PeerCluster.Peers = expectedPeers
					mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)

					addr0 := AddHttpProtocol(expectedPeers[0].Addr)
					c1 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addr0)).
						Times(1)
					res := &RawResponse{
						StatusCode: http.StatusOK,
						Body:       []byte("http://" + peerAddrItself),
						Header:     nil,
					}
					c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
						Return(res, nil).Times(1).After(c1)
					mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).
						AnyTimes().After(c2)
					mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
						Return(nil, cErr).AnyTimes().After(c2)

					go func() {
						clusterInspector.Run(ch)
					}()
					time.Sleep(time.Duration(1) * time.Second)
				})
			})
			Context("when reachable peer is leader", func() {
				Context("when size of foreign cluster is less than own", func() {
					It("should continue searching for leader", func() {
						clusterInspector.PeerCluster.Peers = expectedPeers
						addrURL := AddHttpProtocol(expectedPeers[0].Addr)
						peerURL := AddHttpProtocol(expectedPeers[0].PeerAddr)
						getLeaderResponse := &RawResponse{
							StatusCode: http.StatusOK,
							Body:       []byte(peerURL),
							Header:     nil,
						}
						nodes := []*Node{node0}
						getClusterResponse := makeGetClusterResponse(nodes)
						mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)
						c1 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
							Times(1)
						c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
							Return(getLeaderResponse, nil).Times(1).After(c1)
						c3 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
							Times(1).After(c2)
						endCh := make(chan bool)
						c4 := mockEtcdClient.EXPECT().Get(gomock.Eq(ClusterKey), gomock.Eq(true), gomock.Eq(true)).
							Do(func(_ interface{}, _ interface{}, _ interface{}) {
							endCh <- true
						}).
							Return(getClusterResponse, nil).Times(1).After(c3)
						mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).
							AnyTimes().After(c4)
						mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
							Return(nil, cErr).AnyTimes().After(c4)

						go func() {
							clusterInspector.Run(ch)
						}()
						<-endCh
					})
				})
				Context("when size of foreign cluster is greater than own", func() {
					It("should emit new leader found", func() {
						clusterInspector.PeerCluster.Peers = expectedPeers
						addrURL := AddHttpProtocol(expectedPeers[0].Addr)
						peerURL := AddHttpProtocol(expectedPeers[0].PeerAddr)
						getLeaderResponse := &RawResponse{
							StatusCode: http.StatusOK,
							Body:       []byte(peerURL),
							Header:     nil,
						}
						nodes := []*Node{node0, node1, node2}
						getClusterResponse := makeGetClusterResponse(nodes)
						mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)
						c1 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
							Times(1)
						c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
							Return(getLeaderResponse, nil).Times(1).After(c1)
						c3 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
							Times(1).After(c2)
						mockEtcdClient.EXPECT().Get(gomock.Eq(ClusterKey), gomock.Eq(true), gomock.Eq(true)).
							Return(getClusterResponse, nil).Times(1).After(c3)

						go func() {
							clusterInspector.Run(ch)
						}()
						Expect(<-ch).To(Equal(peerURL))
					})
				})
				Context("when size of foreign cluster is equal than own", func() {
					Context("when foreign leader has less priority than it", func() {
						It("should continue searching for leader", func() {
							clusterInspector.PeerCluster.Peers = expectedPeers
							addrURL := AddHttpProtocol(expectedPeers[1].Addr)
							peerURL := AddHttpProtocol(expectedPeers[1].PeerAddr)
							getLeaderResponse := &RawResponse{
								StatusCode: http.StatusOK,
								Body:       []byte(peerURL),
								Header:     nil,
							}
							nodes := []*Node{node0, node1}
							getClusterResponse := makeGetClusterResponse(nodes)
							mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)
							c1 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(AddHttpProtocol(expectedPeers[0].Addr))).
								Times(1)
							c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
								Return(nil, cErr).Times(1).After(c1)
							c3 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
								Times(1).After(c2)
							c4 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
								Return(getLeaderResponse, nil).Times(1).After(c3)
							c5 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
								Times(1).After(c4)
							c6 := mockEtcdClient.EXPECT().Get(gomock.Eq(ClusterKey), gomock.Eq(true), gomock.Eq(true)).
								Return(getClusterResponse, nil).Times(1).After(c5)
							mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).
								AnyTimes().After(c6)
							mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
								Return(nil, cErr).AnyTimes().After(c6)

							go func() {
								clusterInspector.Run(ch)
							}()
							time.Sleep(time.Duration(1) * time.Second)
						})
					})
					Context("when foreign leader has greater priority than it", func() {
						It("should emit new leader found", func() {
							clusterInspector.PeerCluster.Peers = []Peer{expectedPeers[1:2][0], expectedPeers[0:1][0]}
							addrURL := AddHttpProtocol(expectedPeers[1].Addr)
							peerURL := AddHttpProtocol(expectedPeers[1].PeerAddr)
							getLeaderResponse := &RawResponse{
								StatusCode: http.StatusOK,
								Body:       []byte(peerURL),
								Header:     nil,
							}
							nodes := []*Node{node1, node0}
							getClusterResponse := makeGetClusterResponse(nodes)
							mockPeersMonitor.EXPECT().Run(gomock.Any()).Times(1)
							c1 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
								Times(1)
							c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
								Return(getLeaderResponse, nil).Times(1).After(c1)
							c3 := mockEtcdClient.EXPECT().WithMachineAddr(gomock.Eq(addrURL)).
								Times(1).After(c2)
							c4 := mockEtcdClient.EXPECT().Get(gomock.Eq(ClusterKey), gomock.Eq(true), gomock.Eq(true)).
								Return(getClusterResponse, nil).Times(1).After(c3)
							mockEtcdClient.EXPECT().WithMachineAddr(gomock.Any()).
								AnyTimes().After(c4)
							mockEtcdClient.EXPECT().BaseGet(gomock.Eq(LeaderKey)).
								Return(nil, cErr).AnyTimes().After(c4)

							go func() {
								clusterInspector.Run(ch)
							}()
							Expect(<-ch).To(Equal(peerURL))
						})
					})
				})
			})
		})
	})
})
