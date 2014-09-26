package supervisor_test

import (
	. "github.com/innotech/hydra/supervisor"

	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("Supervisor", func() {
	Context("when node can not write", func() {
		It("should retry writing", func() {

		})
		Context("when the cycle of retries ends unsuccessfully", func() {
			It("should restart the node in master prima state", func() {

			})
		})
	})
	Context("when node state is master prima", func() {
		It("should try to connect to saved peers", func() {

		})
		Context("when node find a peer to connect", func() {
			It("should send request to connect", func() {

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
