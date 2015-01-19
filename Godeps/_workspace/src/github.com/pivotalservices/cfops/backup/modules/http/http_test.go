package http_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/backup/modules/http"
)

var (
	roundTripSuccess bool
	mockResponse     *http.Response
	requestCatcher   *http.Request
	handlerSuccess   bool
)

type MockRoundTripper struct {
}

func (roundTripper *MockRoundTripper) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	if !roundTripSuccess {
		err = errors.New("Mock error")
		return
	}
	*requestCatcher = *request
	resp = mockResponse
	return resp, err
}

type MockHandler struct {
}

func (handler *MockHandler) Handle(resp *http.Response) (val interface{},
	err error) {
	if !handlerSuccess {
		return nil, errors.New("Mock error")
	}
	return "Success", nil
}

var _ = Describe("Http", func() {
	var (
		handler  *MockHandler
		executor *HttpGateway
	)
	BeforeEach(func() {
		requestCatcher = &http.Request{}
		handler = &MockHandler{}
		executor = NewHttpGateway("http://endpoint/test", "username", "password", "contentType", handler)
		NewRoundTripper = func() http.RoundTripper {
			return &MockRoundTripper{}
		}
	})

	Context("The http is request and handled successfully", func() {
		BeforeEach(func() {
			roundTripSuccess = true
			handlerSuccess = true
		})
		It("Should return nil error on success", func() {
			_, err := executor.Execute("Get")
			Ω(err).Should(BeNil())
		})
		It("Should execute correct request", func() {
			val, _ := executor.Execute("Get")
			Ω(requestCatcher.URL.Host).Should(Equal("endpoint"))
			Ω(requestCatcher.Method).Should(Equal("Get"))
			Ω(requestCatcher.Header["Content-Type"][0]).Should(Equal("contentType"))
			Ω(requestCatcher.Header["Authorization"][0]).Should(Equal("Basic dXNlcm5hbWU6cGFzc3dvcmQ="))
			Ω(val).Should(Equal("Success"))
		})
	})

	Context("The round trip request failed", func() {
		BeforeEach(func() {
			roundTripSuccess = false
			handlerSuccess = true
		})
		It("Should return error", func() {
			_, err := executor.Execute("Get")
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("The handler failed", func() {
		BeforeEach(func() {
			roundTripSuccess = true
			handlerSuccess = false
		})
		It("Should return error", func() {
			_, err := executor.Execute("Get")
			Ω(err).ShouldNot(BeNil())
		})
	})

})
