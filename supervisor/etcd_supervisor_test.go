package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"errors"
	"time"
)

var _ = Describe("EtcdSupervisor", func() {
	var (
		mockCtrl        *gomock.Controller
		mockEtcdClient  *mock.MockEtcdClient
		mockEtcdManager *mock.MockEtcdManager
		etcdSupervisor  *EtcdSupervisor
	)

	const (
		EtcdStoreRootPath string = "/v2/keys/"
	)

	// var refreshInterval time.Duration = time.Duration(3000) * time.Millisecond

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdClient = mock.NewMockEtcdClient(mockCtrl)
		mockEtcdManager = mock.NewMockEtcdManager(mockCtrl)
		etcdSupervisor = NewEtcdSupervisor(mockEtcdClient)
	})

	AfterEach(func() {
		// TODO: Supervisor Stop
		mockCtrl.Finish()
	})

	Context("when node can not publish the node state", func() {
		// TODO: check if request remains waiting instead of throw timeout
		It("should retry publishing", func() {
			key := ClusterPath + etcdSupervisor.GetPeerAddr() + StateKey
			ttl := etcdSupervisor.GetDurationBetweenPublicationsState() + 1
			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Any(), gomock.Eq(ttl),
				gomock.Any(), gomock.Any()).Times(DefaultWriteAttempts).Return(nil, errors.New("Write operation impossible"))

			go func() {
				etcdSupervisor.Run()
			}()
		})
		Context("when the cycle of retries ends unsuccessfully", func() {
			It("should restart the node in master prima state", func() {
				isMaster := false
				key := ClusterPath + etcdSupervisor.GetPeerAddr() + StateKey
				ttl := etcdSupervisor.GetDurationBetweenPublicationsState() + 1
				retriesCall := mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Any(), gomock.Eq(ttl),
					gomock.Any(), gomock.Any()).Times(DefaultWriteAttempts).Return(nil, errors.New("Write operation impossible"))
				restartCall := mockEtcdManager.EXPECT().restartEtcdService().Times(1).After(retriesCall)
				mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Eq(StateMasterP), gomock.Eq(ttl),
					gomock.Any(), gomock.Any()).After(restartCall)
				go func() {
					etcdSupervisor.Run()
				}()

				tolarableGapDuration := time.Duration("3s")
				duration := (etcdSupervisor.GetMaxAttemptsToSetState() * etcdSupervisor.GetDurationBetweenPublicationsState()) +
					(etcdSupervisor.GetDurationBetweenPublicationsState() + tolarableGapDuration)
				pollingInterval := etcdSupervisor.GetDurationBetweenPublicationsState()
				Consistently(func() uint {
					return etcdSupervisor.GetState()
				}, duration, pollingInterval).Should(Equal(StateMasterP))
			})
		})
	})
	Context("when node state is master prima", func() {
		It("should publish its state as master prima", func() {
			key := ClusterPath + etcdSupervisor.GetPeerAddr() + StateKey
			ttl := etcdSupervisor.GetDurationBetweenPublicationsState() + 1
			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(key), gomock.Eq(StateMasterP), gomock.Eq(ttl),
				gomock.Any(), gomock.Any())
			go func() {
				etcdSupervisor.Run()
			}()
		})
		// It("should wait for a change in requestToJoin key", func() {
		// 	mockEtcdClient.EXPECT().Watch(gomock.Eq(ClusterRootPath+HydraNodeId+"/"+requestToJoin), gomock.Any(), gomock.Eq(false),
		// 		gomock.Any(), gomock.Any()).Times(1)

		// 	go func() {
		// 		etcdSupervisor.Run()
		// 	}()
		// })
		Context("when none of the foreign peers is accessible", func() {
			It("should try to connect with all saved peers that are not part of the cluster", func() {
				peers := etcdSupervisor.GetForeignPeers()
				var f *gomock.Call
				for i, peer := range peers {
					if i > 0 {
						f = mockEtcdClient.EXPECT().GetLeader(gomock.Eq(peers[i])).After(f).Return(nil, errors.New("Service is not accesible"))
					} else {
						f = mockEtcdClient.EXPECT().GetLeader(gomock.Eq(peers[i])).Return(nil, errors.New("Service is not accesible"))
					}
				}

				go func() {
					etcdSupervisor.Run()
				}()
			})
		})
		Context("when node find a peer to connect", func() {
			// It("should check if this peer is leader", func() {

			// })
			Context("when accesible peer is not leader", func() {

			})
			Context("when accesible peer is leader", func() {

			})
			It("should send request to connect", func() {
				mockEtcdClient.EXPECT().Set().Times(1)

				go func() {
					etcdSupervisor.Run()
				}()
			})
			Context("when request to connect is accepted", func() {
				It("should try to connect", func() {

				})
				Context("when connect to peer", func() {
					It("should change its state to slave", func() {

					})
				})
				Context("when connect to peer is impossible", func() {
					It("should remain as master prima", func() {

					})
				})
			})
			Context("when request to connect is not accepted", func() {
				It("should try to connect", func() {

				})
			})
		})
		Context("when node gets a request to connect", func() {
			It("should wait for connection", func() {

			})
			Context("when connection time expires", func() {
				It("should finish of waiting for connection", func() {

				})
			})
			Context("when the connection is established", func() {
				Context("when cluster is not complete", func() {
					It("should remain as master prima", func() {

					})
				})
				Context("when cluster is complete", func() {
					It("should remain as master", func() {

					})
				})
			})
		})
	})
})
