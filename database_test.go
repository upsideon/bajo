package main

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/syndtr/goleveldb/leveldb"

	mocks "github.com/upsideon/bajo/mocks"
)

var _ = Describe("Database", func() {
	var ctrl *gomock.Controller
	var mockDatabaseManager *mocks.MockDatabaseManager

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	Describe("GetURLDatabase", func() {
		BeforeEach(func() {
			mockDatabaseManager = mocks.NewMockDatabaseManager(ctrl)
		})

		Context("and an error occurs retrieving URL database", func() {
			BeforeEach(func() {
				mockDatabaseManager.EXPECT().OpenFile(
					defaultDatabasePath, nil,
				).Return(nil, errors.New("failed to retrieve database"))
			})

			It("should panic with an error message", func() {
				defer func() {
					if r := recover(); r == nil {
						Fail("did not panic when unable to access database")
					} else {
						Expect(r).To(
							Equal("Error: Unable to access URL database: failed to retrieve database"),
						)
					}
				}()

				GetURLDatabase(mockDatabaseManager)
			})
		})

		Context("and the URL database is retrieved successfully", func() {
			var urlDatabase *leveldb.DB

			BeforeEach(func() {
				urlDatabase = &leveldb.DB{}

				mockDatabaseManager.EXPECT().OpenFile(
					defaultDatabasePath, nil,
				).Return(urlDatabase, nil)
			})

			It("should return the URL database", func() {
				returnedDatabase := GetURLDatabase(mockDatabaseManager)
				Expect(returnedDatabase).To(Equal(urlDatabase))
			})
		})
	})
})
