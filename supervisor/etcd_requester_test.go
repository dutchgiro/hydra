package supervisor_test

// import (
// 	. "github.com/innotech/hydra/supervisor"

// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

// 	"net/http"
// 	"time"
// )

// var _ = Describe("EtcdRequester", func() {
// 	var (
// 		etcdRequester *EtcdRequester
// 		hydraNodeId   string
// 	)

// 	BeforeEach(func() {
// 		etcdRequester = NewEtcdRequester()
// 		hydraNodeId = "hydra_0"
// 	})

// 	Describe("CheckWritingAbility", func() {
// 		Context("when etcd is not accesible", func() {
// 			It("should throw an error", func() {
// 				response, err := etcdRequester.CheckWritingAbility("htttp://localhost:3537"+ClusterRootPath, hydraNodeId)

// 				Expect(response).To(BeNil())
// 				Expect(err).To(HaveOccurred())
// 			})
// 		})
// 		Context("when etcd server is accessible", func() {
// 			Context("when response status code is different from 200 OK", func() {
// 				It("should throw an error", func() {
// 					routes := []Route{
// 						Route{
// 							Pattern: ClusterRootPath + hydraNodeId,
// 							Handler: func(w http.ResponseWriter, r *http.Request) {
// 								w.WriteHeader(http.StatusNotFound)
// 							},
// 						},
// 					}
// 					ts := RunHydraServerMock(routes)
// 					defer ts.Close()
// 					time.Sleep(time.Duration(200) * time.Millisecond)

// 					response, err := etcdRequester.CheckWritingAbility(ts.URL+ClusterRootPath, hydraNodeId)

// 					Expect(response).To(BeNil())
// 					Expect(err).To(HaveOccurred())
// 				})
// 			})
// 			Context("when response status code is 200 OK", func() {
// 				It("should return a server list", func() {
// 					jsonOutput = []byte(`{
// 					    "action": "compareAndSwap",
// 					    "node": {
// 					        "createdIndex": 8,
// 					        "key": "/foo",
// 					        "modifiedIndex": 9,
// 					        "value": "two"
// 					    },
// 					    "prevNode": {
// 					        "createdIndex": 8,
// 					        "key": "/foo",
// 					        "modifiedIndex": 8,
// 					        "value": "one"
// 					    }
// 					}`)
// 					routes := []Route{
// 						Route{
// 							Pattern: ClusterRootPath + hydraNodeId,
// 							Handler: func(w http.ResponseWriter, r *http.Request) {
// 								// jsonOutput, _ := json.Marshal(outputServers)
// 								w.WriteHeader(http.StatusOK)
// 								w.Header().Set("Content-Type", "application/json")
// 								w.Header().Set("X-Etcd-Index", "4")
// 								w.Header().Set("X-Raft-Index", "106")
// 								w.Header().Set("X-Raft-Term", "0")
// 								w.Write(jsonOutput)
// 							},
// 						},
// 					}
// 					ts := RunHydraServerMock(routes)
// 					defer ts.Close()
// 					time.Sleep(time.Duration(200) * time.Millisecond)

// 					response, err := etcdRequester.CheckWritingAbility(ts.URL+ClusterRootPath, hydraNodeId)

// 					Expect(err).ToNot(HaveOccurred())
// 					var expected Response
// 					Expect(response).To(BeAssignableToTypeOf(expected))
// 					// TODO: check struct attributes
// 				})
// 			})
// 		})
// 	})
// })
