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

var _ = Describe("User", func() {
	var (
		id    string
		email string

		db   *sql.DB
		mock sqlmock.Sqlmock
		user User
	)

	BeforeEach(func() {
		id = generateID()
		email = "user@example.com"

		db, mock = NewSQLMock()
		user = NewUser(db)
	})

	Describe("Create", func() {
		Context("success", func() {
			var res *models.User

			BeforeEach(func() {
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(sqlmock.AnyArg(), email).
					WillReturnResult(sqlmock.NewResult(0, 1))

				res, _ = user.Create(
					context.TODO(),
					&models.UserCreateRequest{Email: email})
			})

			It("returns newly created user", func() {
				Expect(res).NotTo(BeNil())
				Expect(res.Email).To(Equal(email))
				Expect(res.Consents).To(BeEmpty())
			})
		})

		Context("error inserting", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(sqlmock.AnyArg(), email).
					WillReturnError(fmt.Errorf("insert error"))

				_, e = user.Create(
					context.TODO(),
					&models.UserCreateRequest{Email: email})
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})
	})

	Describe("Delete", func() {
		Context("success", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"events\"").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("DELETE FROM \"users\"").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				e = user.Delete(context.TODO(), id)
			})

			It("does not return any error", func() {
				Expect(e).To(BeNil())
			})
		})

		Context("error in transaction begin", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("begin error"))

				e = user.Delete(context.TODO(), id)
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})

		Context("error in events record delete", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"events\"").
					WithArgs(id).
					WillReturnError(fmt.Errorf("delete error"))
				mock.ExpectRollback()

				e = user.Delete(context.TODO(), id)
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})

		Context("error in users record delete", func() {
			var e error

			BeforeEach(func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"events\"").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectExec("DELETE FROM \"users\"").
					WithArgs(id).
					WillReturnError(fmt.Errorf("delete error"))
				mock.ExpectRollback()

				e = user.Delete(context.TODO(), id)
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})
	})

	Describe("Detail", func() {
		Context("existent", func() {
			var res *models.User

			BeforeEach(func() {
				userRow := mock.NewRows([]string{"id", "email"}).
					AddRow(id, email)

				mock.ExpectQuery("FROM \"users\"").
					WithArgs(id).
					WillReturnRows(userRow)

				eventRows := mock.NewRows([]string{"consent_id", "enabled"}).
					AddRow(models.ConsentEmail, true).
					AddRow(models.ConsentSMS, false)

				mock.ExpectQuery("FROM \"events\"").
					WithArgs(id, models.ConsentEmail, models.ConsentSMS).
					WillReturnRows(eventRows).
					RowsWillBeClosed()

				res, _ = user.Detail(context.TODO(), id)
			})

			It("returns matching user", func() {
				Expect(res).NotTo(BeNil())
				Expect(res.ID).To(Equal(id))
				Expect(res.Email).To(Equal(email))
				Expect(res.Consents).NotTo(BeEmpty())
			})
		})

		Context("non-existent", func() {
			var res *models.User

			BeforeEach(func() {
				mock.ExpectQuery("FROM \"users\"").
					WithArgs(id).
					WillReturnError(fmt.Errorf("no matching record"))

				res, _ = user.Detail(context.TODO(), id)
			})

			It("returns nil", func() {
				Expect(res).To(BeNil())
			})
		})

		Context("error querying event records", func() {
			var e error

			BeforeEach(func() {
				userRow := mock.NewRows([]string{"id", "email"}).
					AddRow(id, email)

				mock.ExpectQuery("FROM \"users\"").
					WithArgs(id).
					WillReturnRows(userRow)

				mock.ExpectQuery("FROM \"events\"").
					WithArgs(id, models.ConsentEmail, models.ConsentSMS).
					WillReturnError(fmt.Errorf("query error"))

				_, e = user.Detail(context.TODO(), id)
			})

			It("returns error", func() {
				Expect(e).NotTo(BeNil())
			})
		})
	})
})
