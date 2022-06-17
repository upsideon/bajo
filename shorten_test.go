package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mocks "github.com/upsideon/bajo/mocks"
)

const shortenURL = "/shorten"

var _ = Describe(shortenURL, func() {
	var router *gin.Engine
	var writer *httptest.ResponseRecorder

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		urlDatabase := mocks.NewMockURLDatabase(ctrl)
		router = initializeRouter(urlDatabase)
		writer = httptest.NewRecorder()
	})

	When("a URL is not provided", func() {
		JustBeforeEach(func() {
			request, _ := http.NewRequest("POST", shortenURL, nil)
			router.ServeHTTP(writer, request)
		})

		It("returns a 400", func() {
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
		})

		It("return error message", func() {
			Expect(writer.Body.String()).To(Equal("Bad Request"))
		})
	})

	When("a URL is provided", func() {
		var url string

		JustBeforeEach(func() {
			request_content := map[string]string{
				"url": url,
			}
			request_body, _ := json.Marshal(request_content)

			request, _ := http.NewRequest("POST", shortenURL, bytes.NewReader(request_body))
			router.ServeHTTP(writer, request)
		})

		Context("and the URL is valid", func() {
			BeforeEach(func() {
				url = "https://en.wikipedia.org/wiki/URL_shortening"
			})

			It("returns a 200", func() {
				Expect(writer.Code).To(Equal(http.StatusOK))
			})

			It("returns a shortened URL", func() {
				expected_response_content := map[string]string{
					"shortened_url": "https://bajo/oROh-p8o",
				}
				expected_json, _ := json.Marshal(expected_response_content)

				Expect(writer.Body.String()).To(Equal(string(expected_json)))
			})
		})
	})
})
