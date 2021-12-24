package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kazimanzurrashid/consents-api-go/models"
)

var _ = Describe("Event", func() {
	Describe("Create", func() {
		var (
			userID string
			req    *models.EventCreateRequest

			db    *sql.DB
			mock  sqlmock.Sqlmock
			event Event
		)

		BeforeEach(func() {
			userID = generateID()
			req = &models.EventCreateRequest{
				User: models.EventCreateUser{
					ID: userID,
				},
				Consents: []models.Consent{
					{
						ID:      models.ConsentEmail,
						Enabled: true,
					},
					{
						ID:      models.ConsentSMS,
						Enabled: false,
					},
				},
			}

			db, mock = NewSQLMock()
			event = NewEvent(db)
		})

		Context("success", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"public\".\"events\"").
					WithArgs(
						sqlmock.AnyArg(),
						userID,
						models.ConsentEmail,
						sqlmock.AnyArg(),
						true).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("INSERT INTO \"public\".\"events\"").
					WithArgs(
						sqlmock.AnyArg(),
						userID,
						models.ConsentSMS,
						sqlmock.AnyArg(),
						false).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				e = event.Create(context.TODO(), req)
			})

			It("does not return any error", func() {
				Expect(e).To(BeNil())
			})
		})

		Context("error in transaction begin", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("begin error"))

				e = event.Create(context.TODO(), req)
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})

		Context("error in record insert", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"public\".\"events\"").
					WithArgs(
						sqlmock.AnyArg(),
						userID,
						models.ConsentEmail,
						sqlmock.AnyArg(),
						true).
					WillReturnError(fmt.Errorf("insert error"))
				mock.ExpectRollback()

				e = event.Create(context.TODO(), req)
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})
	})
})
