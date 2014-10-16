package supervisor_test

// import (
// 	hydra_config "github.com/innotech/hydra/config"
// 	. "github.com/innotech/hydra/supervisor"
// 	mock "github.com/innotech/hydra/supervisor/mock"

// 	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	// . "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

// 	// "errors"
// 	// "fmt"
// 	// "net/http"
// 	"time"
// )

// var _ = Describe("EtcdSupervisor", func() {
// 	var (
// 		mockCtrl *gomock.Controller
// 		// mockEtcdClient  *mock.MockEtcdRequester
// 		mockEtcdManager      *mock.MockEtcdController
// 		mockClusterInspector *mock.MockClusterAnalyzer
// 		hydraConfig          *hydra_config.Config
// 		etcdSupervisor       *EtcdSupervisor
// 		stateKeyPath         string
// 	)

// 	const (
// 		// EtcdStoreRootPath string = "/v2/keys/"
// 		peerAddrItself string = "127.0.0.1:7701"
// 	)

// 	// var refreshInterval time.Duration = time.Duration(3000) * time.Millisecond

// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(GinkgoT())
// 		// mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
// 		mockEtcdManager = mock.NewMockEtcdController(mockCtrl)
// 		mockClusterInspector = mock.NewMockClusterAnalyzer(mockCtrl)
// 		hydraConfig = hydra_config.New()
// 		hydraConfig.Peer.Addr = peerAddrItself
// 		etcdSupervisor = NewEtcdSupervisor(hydraConfig)
// 		// etcdSupervisor.EtcdClient = mockEtcdClient
// 		etcdSupervisor.EtcdManager = mockEtcdManager
// 		etcdSupervisor.ClusterInspector = mockClusterInspector

// 		// stateKeyPath = ClusterRootPath + "/" + etcdSupervisor.GetPeerAddr() + "/" + StateKey
// 	})

// 	AfterEach(func() {
// 		// fmt.Println("AfterEach")
// 		// TODO: Supervisor Stop
// 		// time.Sleep(time.Duration(2) * time.Second)
// 		// etcdSupervisor.Stop()
// 		mockCtrl.Finish()
// 	})

// 	Describe("Run", func() {
// 		Context("when a new leader is emited", func() {
// 			It("should try restart etcd as slave", func() {
// 				clusterInspectorChannel := make(chan Peer)
// 				c1 := mockClusterInspector.EXPECT().Run(gomock.Any()).Do(func(ch chan Peer) {
// 					clusterInspectorChannel = ch
// 				}).Times(1)
// 				mockEtcdManager.EXPECT().Restart().Times(1).After(c1)

// 				go func() {
// 					etcdSupervisor.Run()
// 				}()
// 				peer := Peer{
// 					Addr:     "98.245.153.113:4001",
// 					PeerAddr: "98.245.153.113:7001",
// 					State:    PeerStateEnabled,
// 				}
// 				clusterInspectorChannel <- peer
// 				time.Sleep(time.Duration(100) * time.Millisecond)
// 			})
// 		})
// 		Context("when can not set state", func() {
// 			It("should try restart etcd as slave", func() {
// 			})
// 		})
// 	})
// })

// ---------------------------------------------------------------------

// import (
// 	hydra_config "github.com/innotech/hydra/config"
// 	. "github.com/innotech/hydra/supervisor"
// 	mock "github.com/innotech/hydra/supervisor/mock"

// 	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	// . "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// // TODO: change "peer" for "node"
// var _ = Describe("EtcdSupervisor", func() {
// 	var (
// 		mockCtrl        *gomock.Controller
// 		mockEtcdClient  *mock.MockEtcdRequester
// 		mockEtcdManager *mock.MockEtcdController
// 		hydraConfig     *hydra_config.Config
// 		etcdSupervisor  *EtcdSupervisor
// 		stateKeyPath    string
// 	)

// 	const (
// 		// EtcdStoreRootPath string = "/v2/keys/"
// 		peerAddrItself string = "127.0.0.1:7701"
// 	)

// 	// var refreshInterval time.Duration = time.Duration(3000) * time.Millisecond

// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(GinkgoT())
// 		mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
// 		mockEtcdManager = mock.NewMockEtcdController(mockCtrl)
// 		hydraConfig = hydra_config.New()
// 		hydraConfig.Peer.Addr = peerAddrItself
// 		etcdSupervisor = NewEtcdSupervisor(hydraConfig)
// 		etcdSupervisor.EtcdClient = mockEtcdClient
// 		etcdSupervisor.EtcdManager = mockEtcdManager

// 		stateKeyPath = ClusterRootPath + "/" + etcdSupervisor.GetPeerAddr() + "/" + StateKey
// 	})

// 	AfterEach(func() {
// 		fmt.Println("AfterEach")
// 		// TODO: Supervisor Stop
// 		// time.Sleep(time.Duration(2) * time.Second)
// 		// etcdSupervisor.Stop()
// 		mockCtrl.Finish()
// 	})

// 	// It("should be instantiated with default configuration", func() {
// 	// 	// Fail("None test")
// 	// 	// Expect(EtcdSupervisorFactory.GetMaxWriteAttempts()).To(Equal(DefaultWriteAttempts))
// 	// 	// Expect(EtcdSupervisorFactory.GetDurationBetweenAllServersRetry()).To(Equal(DefaultDurationBetweenAllServersRetry))
// 	// 	// Expect(EtcdSupervisorFactory.GetHydraServersCacheDuration()).To(Equal(DefaultHydraServersCacheDuration))
// 	// 	// Expect(EtcdSupervisorFactory.GetMaxNumberOfRetriesPerHydraServer()).To(Equal(DefaultNumberOfRetries))
// 	// })

// 	const modifiedIndex uint64 = 2
// 	successResponse := &Response{
// 		Action: "",
// 		Node: &Node{
// 			Key:           "",
// 			Value:         "",
// 			Dir:           false,
// 			Expiration:    nil,
// 			TTL:           3,
// 			Nodes:         nil,
// 			ModifiedIndex: modifiedIndex,
// 			CreatedIndex:  modifiedIndex - 1,
// 		},
// 		PrevNode:  nil,
// 		EtcdIndex: 0,
// 		RaftIndex: 0,
// 		RaftTerm:  0,
// 	}
// 	successGetLeaderResponse := &RawResponse{
// 		StatusCode: http.StatusOK,
// 		Body:       []byte("http://" + peerAddrItself),
// 		Header:     nil,
// 	}

// 	Context("when the node starts", func() {
// 		It("should be able to post its state", func() {
// 			c1 := mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateKeyPath), gomock.Eq(StateAvailable),
// 				gomock.Eq(etcdSupervisor.GetStateTTL()), gomock.Eq(""), gomock.Eq(uint64(0)), gomock.Eq(False)).
// 				Return(successResponse, nil).Times(1)

// 			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateKeyPath), gomock.Eq(StateAvailable),
// 				gomock.Eq(etcdSupervisor.GetStateTTL()), gomock.Eq(""), gomock.Not(uint64(0)), gomock.Eq(Unknow)).
// 				Return(successResponse, nil).AnyTimes().After(c1)

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 			time.Sleep(time.Duration(1) * time.Second)
// 		})
// 		It("should search for a peer to connect", func() {
// 			c1 := mockEtcdClient.EXPECT().BaseGet(gomock.Not("http://"+peerAddrItself)).Return(successGetLeaderResponse, nil).Times(1)
// 			mockEtcdClient.EXPECT().BaseGet(gomock.Not("http://"+peerAddrItself)).Return(successGetLeaderResponse, nil).AnyTimes().After(c1)

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 			time.Sleep(time.Duration(3) * time.Second)
// 		})
// 	})

// 	Context("when the node can not publish the node state", func() {
// 		It("should retry publishing its state for a finite number of times and then restart the node", func() {
// 			c1 := mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateKeyPath), gomock.Eq(StateAvailable),
// 				gomock.Eq(etcdSupervisor.GetStateTTL()), gomock.Eq(""), gomock.Eq(uint64(0)), gomock.Eq(False)).
// 				Return(nil, errors.New("Write operation impossible")).Times(int(etcdSupervisor.GetNumOfSetStateRetries()))

// 			c2 := mockEtcdManager.EXPECT().Restart().Times(1).After(c1)

// 			c3 := mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateKeyPath), gomock.Eq(StateAvailable),
// 				gomock.Eq(etcdSupervisor.GetStateTTL()), gomock.Eq(""), gomock.Eq(uint64(0)), gomock.Eq(False)).
// 				Return(successResponse, nil).Times(1).After(c2)

// 			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateKeyPath), gomock.Eq(StateAvailable),
// 				gomock.Eq(etcdSupervisor.GetStateTTL()), gomock.Eq(""), gomock.Any(), gomock.Any()).
// 				AnyTimes().After(c3)

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 			time.Sleep(time.Duration(2) * time.Second)
// 		})
// 	})

// 	Context("when node find a foreign peer to connect", func() {
// 		// It("should check if this peer is leader", func() {

// 		// })
// 		Context("when reachable peer is not leader", func() {
// 			It("should continue searching a peer to connect", func() {
// 				foreignPeerAddr := "98.245.153.111:7701"
// 				// res := *RawResponse{
// 				// 	StatusCode: http.StatusOK,
// 				// 	Body:       []byte("http://" + peerAddrItself),
// 				// 	Header:     nil,
// 				// }
// 				c1 := mockEtcdClient.EXPECT().WithMachineAddr("http://" + foreignPeerAddr).Times(1)
// 				c2 := mockEtcdClient.EXPECT().BaseGet(gomock.Eq("leader")).Return(successGetLeaderResponse, nil).Times(1).After(c1)
// 				mockEtcdClient.EXPECT().WithMachineAddr(gomock.Not("http://" + peerAddrItself)).AnyTimes().After(c2)
// 				mockEtcdClient.EXPECT().BaseGet(gomock.Eq("leader")).Return(successGetLeaderResponse, nil).AnyTimes().After(c2)

// 				go func() {
// 					etcdSupervisor.Run()
// 				}()
// 				time.Sleep(time.Duration(2) * time.Second)
// 			})
// 		})
// Context("when reachable peer is leader", func() {
// 	// Next request for machines and compare with local cluster for overwrite local or remote cluster
// 	It("should get registered peers from foreign peer", func() {
// 		f1 := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).
// 			SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 		f2 := mockEtcdClient.EXPECT().Get(gomock.Any(), gomock.Any()).
// 			SetArg("http://98.245.153.111:7001", "/cluster")

// 		go func() {
// 			etcdSupervisor.Run()
// 		}()
// 	})
// 	Context("when the registered peers from foreign peer leader are retrieved", func() {
// 		It("should add unknown peers to foreign leader", func() {

// 		})
// 	})

// 	// TODO: priority = higher cluter size && node priority
// 	Context("when it has lower cluster priority than reachable peer", func() {
// 		It("should try connecting to reachable peer", func() {
// 			successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 			mockEtcdClient.EXPECT().GetLeader(gomock.Any()).AnyTimes().After(successfulCallToSlaveNode)
// 			// TODO: restart with attributes
// 			mockEtcdManager.EXPECT().RestartEtcdService().Times(1).After(successfulCallToSlaveNode)

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 		})
// 	})
// 	Context("when it has higher cluster priority than reachable peer", func() {
// 		It("should remain looking for a foreign peer to connect", func() {
// 			successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 			mockEtcdClient.EXPECT().GetLeader(gomock.Any()).After(successfulCallToSlaveNode)
// 		})
// 	})
// })
// It("should send request to connect", func() {
// 	mockEtcdClient.EXPECT().Set().Times(1)

// 	go func() {
// 		etcdSupervisor.Run()
// 	}()
// })
// Context("when request to connect is accepted", func() {
// 	It("should try to connect", func() {

// 	})
// 	Context("when connect to peer", func() {
// 		It("should change its state to slave", func() {

// 		})
// 	})
// 	Context("when connect to peer is impossible", func() {
// 		It("should remain as master prima", func() {

// 		})
// 	})
// })
// Context("when request to connect is not accepted", func() {
// 	It("should try to connect", func() {

// 	})
// })
// })

// COVER

// 	Context("when the cycle of retries ends unsuccessfully", func() {
// 		It("should restart the node in master prima state", func() {
// 			isMaster := false
// 			key := ClusterPath + etcdSupervisor.GetPeerAddr() + StateKey
// 			// TODO: Config delay or constant
// 			delay := time.Duration("3s")
// 			ttl := etcdSupervisor.GetDurationBetweenPublicationsState() + delay
// 			retriesCall := mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Any(), gomock.Eq(ttl),
// 				gomock.Any(), gomock.Any()).Return(nil, errors.New("Write operation impossible")).Times(DefaultWriteAttempts)
// 			restartCall := mockEtcdManager.EXPECT().restartEtcdService().Times(1).After(retriesCall)
// 			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Eq(StateMasterP), gomock.Eq(ttl),
// 				gomock.Any(), gomock.Any()).Return(&Response{}, nil).After(restartCall)
// 			go func() {
// 				etcdSupervisor.Run()
// 			}()

// 			tolarableDelayDuration := time.Duration("3s")
// 			duration := (etcdSupervisor.GetMaxAttemptsToSetState() * etcdSupervisor.GetDurationBetweenPublicationsState()) +
// 				(etcdSupervisor.GetDurationBetweenPublicationsState() + tolarableDelayDuration)
// 			pollingInterval := etcdSupervisor.GetDurationBetweenPublicationsState()
// 			Consistently(func() uint {
// 				return etcdSupervisor.GetState()
// 			}, duration, pollingInterval).Should(Equal(StateMasterP))
// 		})
// 	})

// Context("when node state is master prima", func() {
// 	It("should have master prima state", func() {
// 		Expect(StateMasterP).To(Equal(etcdSupervisor.GetState()))
// 	})
// 	// It("should publish its state as master prima", func() {
// 	// 	key := ClusterPath + etcdSupervisor.GetPeerAddr() + StateKey
// 	// 	ttl := etcdSupervisor.GetDurationBetweenPublicationsState() + 1
// 	// 	mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Eq(StateMasterP), gomock.Eq(ttl),
// 	// 		gomock.Any(), gomock.Any())
// 	// 	go func() {
// 	// 		etcdSupervisor.Run()
// 	// 	}()
// 	// })
// 	// It("should wait for a change in requestToJoin key", func() {
// 	// 	mockEtcdClient.EXPECT().Watch(gomock.Eq(ClusterRootPath+HydraNodeId+"/"+requestToJoin), gomock.Any(), gomock.Eq(false),
// 	// 		gomock.Any(), gomock.Any()).Times(1)

// 	// 	go func() {
// 	// 		etcdSupervisor.Run()
// 	// 	}()
// 	// })
// 	Context("when none of the foreign peers is accessible", func() {
// 		It("should try to connect with all saved peers that are not part of the cluster", func() {
// 			peers := etcdSupervisor.GetForeignPeers()
// 			var f *gomock.Call
// 			for i, peer := range peers {
// 				if i > 0 {
// 					f = mockEtcdClient.EXPECT().GetLeader(gomock.Eq(peers[i])).After(f).Return(nil, errors.New("Service is not accesible"))
// 				} else {
// 					f = mockEtcdClient.EXPECT().GetLeader(gomock.Eq(peers[i])).Return(nil, errors.New("Service is not accesible"))
// 				}
// 			}

// 			go func() {
// 				etcdSupervisor.Run()
// 			}()
// 		})
// 	})
// 	Context("when node find a foreign peer to connect", func() {
// 		// It("should check if this peer is leader", func() {

// 		// })
// 		Context("when reachable peer is not leader", func() {
// 			It("should continue searching a peer to connect", func() {
// 				successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.112:7001").Return("http://98.245.153.111:7001", nil)
// 				mockEtcdClient.EXPECT().GetLeader(gomock.Any()).After(successfulCallToSlaveNode)

// 				go func() {
// 					etcdSupervisor.Run()
// 				}()
// 			})
// 		})
// 		Context("when reachable peer is leader", func() {
// 			// Next request for machines and compare with local cluster for overwrite local or remote cluster
// 			It("should get registered peers from foreign peer", func() {
// 				f1 := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).
// 					SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 				f2 := mockEtcdClient.EXPECT().Get(gomock.Any(), gomock.Any()).
// 					SetArg("http://98.245.153.111:7001", "/cluster")

// 				go func() {
// 					etcdSupervisor.Run()
// 				}()
// 			})
// 			Context("when the registered peers from foreign peer leader are retrieved", func() {
// 				It("should add unknown peers to foreign leader", func() {

// 				})
// 			})

// 			// TODO: priority = higher cluter size && node priority
// 			Context("when it has lower cluster priority than reachable peer", func() {
// 				It("should try connecting to reachable peer", func() {
// 					successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 					mockEtcdClient.EXPECT().GetLeader(gomock.Any()).AnyTimes().After(successfulCallToSlaveNode)
// 					// TODO: restart with attributes
// 					mockEtcdManager.EXPECT().RestartEtcdService().Times(1).After(successfulCallToSlaveNode)

// 					go func() {
// 						etcdSupervisor.Run()
// 					}()
// 				})
// 			})
// 			Context("when it has higher cluster priority than reachable peer", func() {
// 				It("should remain looking for a foreign peer to connect", func() {
// 					successfulCallToSlaveNode := mockEtcdClient.EXPECT().GetLeader(gomock.Any()).SetArg("http://98.245.153.111:7001").Return("http://98.245.153.111:7001", nil)
// 					mockEtcdClient.EXPECT().GetLeader(gomock.Any()).After(successfulCallToSlaveNode)
// 				})
// 			})
// 		})
// 		// It("should send request to connect", func() {
// 		// 	mockEtcdClient.EXPECT().Set().Times(1)

// 		// 	go func() {
// 		// 		etcdSupervisor.Run()
// 		// 	}()
// 		// })
// 		// Context("when request to connect is accepted", func() {
// 		// 	It("should try to connect", func() {

// 		// 	})
// 		// 	Context("when connect to peer", func() {
// 		// 		It("should change its state to slave", func() {

// 		// 		})
// 		// 	})
// 		// 	Context("when connect to peer is impossible", func() {
// 		// 		It("should remain as master prima", func() {

// 		// 		})
// 		// 	})
// 		// })
// 		// Context("when request to connect is not accepted", func() {
// 		// 	It("should try to connect", func() {

// 		// 	})
// 		// })
// 	})
// 	Context("when node gets a request to connect", func() {
// 		It("should wait for connection", func() {

// 		})
// 		Context("when connection time expires", func() {
// 			It("should finish of waiting for connection", func() {

// 			})
// 		})
// 		Context("when the connection is established", func() {
// 			Context("when cluster is not complete", func() {
// 				It("should remain as master prima", func() {

// 				})
// 			})
// 			Context("when cluster is complete", func() {
// 				It("should remain as master", func() {

// 				})
// 			})
// 		})
// 	})
// })
// })
