package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"
	mock "github.com/innotech/hydra/supervisor/mock"

	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"errors"
	"time"
)

var _ = Describe("StateManager", func() {
	const (
		peerAddrItself                   string        = "127.0.0.1:7701"
		stateKeyTTL                      uint64        = 5
		durationBetweenPublicationsState time.Duration = time.Duration(500) * time.Millisecond
	)

	var (
		ch             chan StateControllerState
		mockCtrl       *gomock.Controller
		mockEtcdClient *mock.MockEtcdRequester
		stateManager   *StateManager
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEtcdClient = mock.NewMockEtcdRequester(mockCtrl)
		stateManager = NewStateManager(durationBetweenPublicationsState,
			mockEtcdClient, peerAddrItself, stateKeyTTL)
		ch = make(chan StateControllerState)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	const modifiedIndex uint64 = 2
	successResponse := &Response{
		Action: "",
		Node: &Node{
			Key:           "",
			Value:         "",
			Dir:           false,
			Expiration:    nil,
			TTL:           3,
			Nodes:         nil,
			ModifiedIndex: modifiedIndex,
			CreatedIndex:  modifiedIndex - 1,
		},
		PrevNode:  nil,
		EtcdIndex: 0,
		RaftIndex: 0,
		RaftTerm:  0,
	}

	Describe("Run", func() {
		It("should have state Running", func() {
			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateManager.PeerStateKey), gomock.Eq(StateAvailable),
				gomock.Eq(stateManager.PeerStateKeyTTL), gomock.Eq(""), gomock.Any(), gomock.Any()).
				Return(successResponse, nil).AnyTimes()

			go func() {
				stateManager.Run(ch)
			}()
			Eventually(func() StateControllerState {
				return stateManager.GetState()
			}).Should(Equal(Running))
		})

		It("should be able to post the node state", func() {
			c1 := mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateManager.PeerStateKey), gomock.Eq(StateAvailable),
				gomock.Eq(stateManager.PeerStateKeyTTL), gomock.Eq(""), gomock.Eq(uint64(0)), gomock.Eq(False)).
				Return(successResponse, nil).Times(1)

			mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateManager.PeerStateKey), gomock.Eq(StateAvailable),
				gomock.Eq(stateManager.PeerStateKeyTTL), gomock.Eq(""), gomock.Not(uint64(0)), gomock.Eq(Unknow)).
				Return(successResponse, nil).AnyTimes().After(c1)

			go func() {
				stateManager.Run(ch)
			}()
			time.Sleep(time.Duration(1) * time.Second)
		})

		Context("when it can not post the node state", func() {
			It("should report the failure and stop", func() {
				mockEtcdClient.EXPECT().CompareAndSwap(gomock.Eq(stateManager.PeerStateKey), gomock.Eq(StateAvailable),
					gomock.Eq(stateManager.PeerStateKeyTTL), gomock.Eq(""), gomock.Eq(uint64(0)), gomock.Eq(False)).
					Return(nil, errors.New("Write operation impossible")).Times(int(stateManager.NumOfSetStateRetries))

				go func() {
					stateManager.Run(ch)
				}()
				Eventually(ch, 3.0, 0.3).Should(Receive())
				Expect(stateManager.GetState()).To(Equal(Stopped))
			})
		})
	})
})
