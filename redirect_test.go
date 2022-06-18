package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dberror "github.com/syndtr/goleveldb/leveldb/errors"

	mocks "github.com/upsideon/bajo/mocks"
)

var _ = Describe("URL redirects", func() {
	var router *gin.Engine
	var mockURLDatabase *mocks.MockURLDatabase
	var urlKey string
	var writer *httptest.ResponseRecorder

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		mockURLDatabase = mocks.NewMockURLDatabase(ctrl)
		router = initializeRouter(mockURLDatabase)
		writer = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", urlKey), nil)
		router.ServeHTTP(writer, request)
	})

	When("no URL key is provided", func() {
		It("returns a 404", func() {
			Expect(writer.Code).To(Equal(http.StatusNotFound))
		})

		It("return error message", func() {
			Expect(writer.Body.String()).To(Equal("404 page not found"))
		})
	})

	When("a URL key is provided", func() {
		BeforeEach(func() {
			urlKey = "Z-f1Fbe7"
		})

		Context("and there is an error retrieving it from database", func() {
			BeforeEach(func() {
				mockURLDatabase.EXPECT().Get(
					[]byte(urlKey), nil,
				).Return(nil, errors.New("failed to query database"))
			})

			It("returns a 500", func() {
				Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			})

			It("return error message", func() {
				Expect(writer.Body.String()).To(Equal("Internal Server Error"))
			})
		})

		Context("and the URL key is not present in database", func() {
			BeforeEach(func() {
				mockURLDatabase.EXPECT().Get(
					[]byte(urlKey), nil,
				).Return(nil, dberror.ErrNotFound)
			})

			It("returns a 404", func() {
				Expect(writer.Code).To(Equal(http.StatusNotFound))
			})

			It("return error message", func() {
				Expect(writer.Body.String()).To(Equal("Not Found"))
			})
		})

		Context("and the URL key is present in database", func() {
			BeforeEach(func() {
				mockURLDatabase.EXPECT().Get(
					[]byte(urlKey), nil,
				).Return([]byte("https://duckduckgo.com/"), nil)
			})

			It("returns a 302", func() {
				Expect(writer.Code).To(Equal(http.StatusFound))
			})
		})
	})
})
