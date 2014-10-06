package supervisor_test

// import (
// 	. "github.com/innotech/hydra/supervisor"

// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
// )

// var _ = Describe("EtcdSupervisorFactory", func() {
// 	var (
// 		mockCtrl       *gomock.Controller
// 		mockSupervisor *mock.MockSupervisor
// 	)

// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(GinkgoT())
// 		mockSupervisor = mock.NewMockSupervisor(mockCtrl)
// 	})

// 	AfterEach(func() {
// 		mockCtrl.Finish()
// 	})

// 	It("should be instantiated with default configuration", func() {
// 		Expect(EtcdSupervisorFactory.GetMaxWriteAttempts()).To(Equal(DefaultWriteAttempts))
// 		// Expect(EtcdSupervisorFactory.GetDurationBetweenAllServersRetry()).To(Equal(DefaultDurationBetweenAllServersRetry))
// 		// Expect(EtcdSupervisorFactory.GetHydraServersCacheDuration()).To(Equal(DefaultHydraServersCacheDuration))
// 		// Expect(EtcdSupervisorFactory.GetMaxNumberOfRetriesPerHydraServer()).To(Equal(DefaultNumberOfRetries))
// 	})

// // 	Describe("Build", func() {
// // 		It("should build a EtcdSupervisor", func() {
// // 			var etcdSupervisor *EtcdSupervisor
// // 			Expect(EtcdSupervisorFactory.Build()).To(BeAssignableToTypeOf(etcdSupervisor))
// // 		})
// // 	})

// 	Describe("Config", func() {

// 		// Context("when hydra server list argument is nil", func() {
// 		// 	It("should throw an error", func() {
// 		// 		err := SupervisorFactory.Config(nil)
// 		// 		Expect(err).Should(HaveOccurred())
// 		// 	})
// 		// })
// 		// Context("when hydra server list argument is a empty list", func() {
// 		// 	It("should throw an error", func() {
// 		// 		err := EtcdSupervisorFactory.Config([]string{})
// 		// 		Expect(err).Should(HaveOccurred())
// 		// 	})
// 		// })
// 		// Context("when hydra server list argument is a valid list of servers", func() {
// 		// 	It("should set hydra server list successfully", func() {
// 		// 		err := EtcdSupervisorFactory.Config([]string{"http://localhost:8080"})
// 		// 		Expect(err).ShouldNot(HaveOccurred())
// 		// 	})
// 		// })
// 	})

// // 	Describe("WithWithMaxWriteAttempts", func() {
// // 		Context("when writeAttempts argument is a valid uint number", func() {
// // 			It("should set apps cache duration successfully", func() {
// // 				const attempts uint = 3
// // 				e := EtcdSupervisorFactory.WithWithMaxWriteAttempts(attempts)
// // 				Expect(e).To(Equal(EtcdSupervisorFactory))
// // 				Expect(attempts).To(Equal(EtcdSupervisorFactory.GetMaxWriteAttempts()))
// // 			})
// // 		})
// // 	})

// // 	// Describe("WithAppsCacheDuration", func() {
// // 	// 	Context("when duration argument is a valid uint number", func() {
// // 	// 		It("should set apps cache duration successfully", func() {
// // 	// 			const appsCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
// // 	// 			h := EtcdSupervisorFactory.WithAppsCacheDuration(appsCacheDuration)
// // 	// 			Expect(h).To(Equal(EtcdSupervisorFactory))
// // 	// 			Expect(appsCacheDuration).To(Equal(EtcdSupervisorFactory.GetAppsCacheDuration()))
// // 	// 		})
// // 	// 	})
// // 	// })

// // 	// Describe("WithMaxNumberOfRetriesPerHydraServer", func() {
// // 	// 	Context("when duration argument is a valid uint number", func() {
// // 	// 		It("should set apps cache duration successfully", func() {
// // 	// 			const retries uint = 3
// // 	// 			h := EtcdSupervisorFactory.WithMaxNumberOfRetriesPerHydraServer(retries)
// // 	// 			Expect(h).To(Equal(EtcdSupervisorFactory))
// // 	// 			Expect(retries).To(Equal(EtcdSupervisorFactory.GetMaxNumberOfRetriesPerHydraServer()))
// // 	// 		})
// // 	// 	})
// // 	// })

// // 	// Describe("WaitBetweenAllServersRetry", func() {
// // 	// 	Context("when duration argument is a valid uint number", func() {
// // 	// 		It("should set wait between all servers retry successfully", func() {
// // 	// 			const appsCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
// // 	// 			h := EtcdSupervisorFactory.WaitBetweenAllServersRetry(appsCacheDuration)
// // 	// 			Expect(h).To(Equal(EtcdSupervisorFactory))
// // 	// 			Expect(appsCacheDuration).To(Equal(EtcdSupervisorFactory.GetDurationBetweenAllServersRetry()))
// // 	// 		})
// // 	// 	})
// // 	// })
// // })
