package main

import (
	"bytes"
	"encoding/json"
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

const shortenURL = "/shorten"

var _ = Describe(shortenURL, func() {
	var router *gin.Engine
	var mockURLDatabase *mocks.MockURLDatabase
	var writer *httptest.ResponseRecorder

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		mockURLDatabase = mocks.NewMockURLDatabase(ctrl)
		router = initializeRouter(mockURLDatabase)
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
		const (
			computedUrlKey = "oROh-p8o"
			customUrlKey   = "custom"
			exampleUrl     = "https://en.wikipedia.org/wiki/URL_shortening"
			invalidUrlKey  = "thiskeyistoolongfortherouteisitnotmyfriend?"
		)

		var requestContent map[string]string

		BeforeEach(func() {
			requestContent = map[string]string{}
		})

		JustBeforeEach(func() {
			requestBody, _ := json.Marshal(requestContent)
			request, _ := http.NewRequest("POST", shortenURL, bytes.NewReader(requestBody))
			router.ServeHTTP(writer, request)
		})

		Context("and the URL is valid", func() {
			BeforeEach(func() {
				requestContent["url"] = exampleUrl
			})

			Context("and a custom URL key is specified", func() {
				Context("and the custom URL key is invalid", func() {
					BeforeEach(func() {
						requestContent["key"] = invalidUrlKey
					})

					It("returns a 400", func() {
						Expect(writer.Code).To(Equal(http.StatusBadRequest))
					})

					It("return error message", func() {
						Expect(writer.Body.String()).To(Equal("Bad Request"))
					})
				})

				Context("and the custom URL key is valid", func() {
					BeforeEach(func() {
						requestContent["key"] = customUrlKey
					})

					Context("and there is an error checking for an existing URL key", func() {
						BeforeEach(func() {
							mockURLDatabase.EXPECT().Get(
								[]byte(customUrlKey), nil,
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
								[]byte(customUrlKey), nil,
							).Return(nil, dberror.ErrNotFound)
						})

						Context("and inserting the URL key fails", func() {
							BeforeEach(func() {
								mockURLDatabase.EXPECT().Put(
									[]byte(customUrlKey), []byte(exampleUrl), nil,
								).Return(errors.New("failed to insert URL key"))
							})

							It("returns a 500", func() {
								Expect(writer.Code).To(Equal(http.StatusInternalServerError))
							})

							It("return error message", func() {
								Expect(writer.Body.String()).To(Equal("Internal Server Error"))
							})
						})

						Context("and inserting the URL key succeeds", func() {
							BeforeEach(func() {
								mockURLDatabase.EXPECT().Put(
									[]byte(customUrlKey), []byte(exampleUrl), nil,
								).Return(nil)
							})

							It("returns a 200", func() {
								Expect(writer.Code).To(Equal(http.StatusOK))
							})

							It("returns a shortened URL", func() {
								expectedResponseContent := map[string]string{
									"shortened_url": fmt.Sprintf("https://bajo/%s", customUrlKey),
								}
								expectedJson, _ := json.Marshal(expectedResponseContent)

								Expect(writer.Body.String()).To(Equal(string(expectedJson)))
							})
						})
					})

					Context("and the URL key is present in database", func() {
						BeforeEach(func() {
							mockURLDatabase.EXPECT().Get(
								[]byte(customUrlKey), nil,
							).Return([]byte(exampleUrl), nil)
						})

						It("returns a 200", func() {
							Expect(writer.Code).To(Equal(http.StatusOK))
						})

						It("returns a shortened URL", func() {
							expectedResponseContent := map[string]string{
								"shortened_url": fmt.Sprintf("https://bajo/%s", customUrlKey),
							}
							expectedJson, _ := json.Marshal(expectedResponseContent)

							Expect(writer.Body.String()).To(Equal(string(expectedJson)))
						})
					})
				})
			})

			Context("and a custom URL key is not specified", func() {
				Context("and there is an error checking for an existing URL key", func() {
					BeforeEach(func() {
						mockURLDatabase.EXPECT().Get(
							[]byte(computedUrlKey), nil,
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
							[]byte(computedUrlKey), nil,
						).Return(nil, dberror.ErrNotFound)
					})

					Context("and inserting the URL key fails", func() {
						BeforeEach(func() {
							mockURLDatabase.EXPECT().Put(
								[]byte(computedUrlKey), []byte(exampleUrl), nil,
							).Return(errors.New("failed to insert URL key"))
						})

						It("returns a 500", func() {
							Expect(writer.Code).To(Equal(http.StatusInternalServerError))
						})

						It("return error message", func() {
							Expect(writer.Body.String()).To(Equal("Internal Server Error"))
						})
					})

					Context("and inserting the URL key succeeds", func() {
						BeforeEach(func() {
							mockURLDatabase.EXPECT().Put(
								[]byte(computedUrlKey), []byte(exampleUrl), nil,
							).Return(nil)
						})

						It("returns a 200", func() {
							Expect(writer.Code).To(Equal(http.StatusOK))
						})

						It("returns a shortened URL", func() {
							expectedResponseContent := map[string]string{
								"shortened_url": fmt.Sprintf("https://bajo/%s", computedUrlKey),
							}
							expectedJson, _ := json.Marshal(expectedResponseContent)

							Expect(writer.Body.String()).To(Equal(string(expectedJson)))
						})
					})
				})

				Context("and the URL key is present in database", func() {
					BeforeEach(func() {
						mockURLDatabase.EXPECT().Get(
							[]byte(computedUrlKey), nil,
						).Return([]byte(exampleUrl), nil)
					})

					It("returns a 200", func() {
						Expect(writer.Code).To(Equal(http.StatusOK))
					})

					It("returns a shortened URL", func() {
						expectedResponseContent := map[string]string{
							"shortened_url": fmt.Sprintf("https://bajo/%s", computedUrlKey),
						}
						expectedJson, _ := json.Marshal(expectedResponseContent)

						Expect(writer.Body.String()).To(Equal(string(expectedJson)))
					})
				})
			})
		})
	})
})
