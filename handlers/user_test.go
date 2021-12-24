package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kazimanzurrashid/consents-api-go/models"
)

var _ = Describe("User", func() {
	Describe("Create", func() {
		Context("success", func() {
			const id = "7b5a3155-7a73-42de-b87e-23f50a10180a"
			const email = "user@example.com"

			var statusCode int
			var res models.User

			BeforeEach(func() {
				var payload bytes.Buffer

				err := json.NewEncoder(&payload).Encode(models.UserCreateRequest{
					Email: email,
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					"/users",
					&payload)

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{
					user: &models.User{
						ID:       id,
						Email:    email,
						Consents: make([]models.Consent, 0),
					},
					err: nil,
				})

				handler := http.HandlerFunc(user.Create)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns newly created user", func() {
				Expect(res.ID).To(Equal(id))
				Expect(res.Email).To(Equal(email))
				Expect(res.Consents).To(BeEmpty())
			})

			It("returns http status code Created", func() {
				Expect(statusCode).To(Equal(http.StatusCreated))
			})
		})

		Context("error reading request body", func() {
			var statusCode int
			var res errorResult

			BeforeEach(func() {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users",
					strings.NewReader("malformed json"))

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{err: nil, user: nil})

				handler := http.HandlerFunc(user.Create)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns malformed request body in errors", func() {
				Expect(res.Errors[0]).To(MatchRegexp("Malformed request"))
			})

			It("returns http status code UnprocessableEntity", func() {
				Expect(statusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})

		Context("invalid request body", func() {
			var statusCode int
			var res errorResult

			BeforeEach(func() {
				var payload bytes.Buffer

				err := json.NewEncoder(&payload).Encode(models.UserCreateRequest{
					Email: "foo-bar",
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					"/users",
					&payload)

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{err: nil, user: nil})

				handler := http.HandlerFunc(user.Create)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns validation errors", func() {
				Expect(res.Errors).NotTo(BeEmpty())
			})

			It("returns http status code UnprocessableEntity", func() {
				Expect(statusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})

		Context("error in service call", func() {
			var statusCode int
			var res errorResult

			BeforeEach(func() {
				var payload bytes.Buffer

				err := json.NewEncoder(&payload).Encode(models.UserCreateRequest{
					Email: "user@example.com",
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					"/users",
					&payload)

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{err: fmt.Errorf("error")})

				handler := http.HandlerFunc(user.Create)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns email exists in errors", func() {
				Expect(res.Errors[0]).To(MatchRegexp("Email already exists"))
			})

			It("returns http status code UnprocessableEntity", func() {
				Expect(statusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})
	})

	Describe("Delete", func() {
		const id = "7b5a3155-7a73-42de-b87e-23f50a10180a"

		Context("success", func() {
			var statusCode int

			BeforeEach(func() {
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/users/%v", id),
					nil)

				if err != nil {
					panic(err)
				}

				req = mux.SetURLVars(req, map[string]string{
					"id": id,
				})

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{err: nil, user: nil})

				handler := http.HandlerFunc(user.Delete)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code
			})

			It("returns http status code NoContent", func() {
				Expect(statusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("error in service call", func() {
			var statusCode int
			var res errorResult

			BeforeEach(func() {
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("/users/%v", id),
					nil)

				if err != nil {
					panic(err)
				}

				req = mux.SetURLVars(req, map[string]string{
					"id": id,
				})

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{
					err:  fmt.Errorf("error"),
					user: nil,
				})

				handler := http.HandlerFunc(user.Delete)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns internal server error in errors", func() {
				Expect(res.Errors[0]).To(MatchRegexp("Internal server error"))
			})

			It("returns http status code InternalServerError", func() {
				Expect(statusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Detail", func() {
		const id = "7b5a3155-7a73-42de-b87e-23f50a10180a"
		const email = "user@example.com"

		Context("existent", func() {
			var statusCode int
			var res models.User

			BeforeEach(func() {
				req, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/users/%v", id),
					nil)

				if err != nil {
					panic(err)
				}

				req = mux.SetURLVars(req, map[string]string{
					"id": id,
				})

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{
					user: &models.User{
						ID:    id,
						Email: email,
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
					},
					err: nil,
				})

				handler := http.HandlerFunc(user.Detail)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns matching user", func() {
				Expect(res.ID).To(Equal(id))
				Expect(res.Email).To(Equal(email))
				Expect(res.Consents[0].ID).To(Equal(models.ConsentEmail))
				Expect(res.Consents[0].Enabled).To(BeTrue())
				Expect(res.Consents[1].ID).To(Equal(models.ConsentSMS))
				Expect(res.Consents[1].Enabled).To(BeFalse())
			})

			It("returns http status code Ok", func() {
				Expect(statusCode).To(Equal(http.StatusOK))
			})
		})

		Context("non-existent", func() {
			var statusCode int
			var res errorResult

			BeforeEach(func() {
				req, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/users/%v", id),
					nil)

				if err != nil {
					panic(err)
				}

				req = mux.SetURLVars(req, map[string]string{
					"id": id,
				})

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{
					user: nil,
					err:  nil,
				})

				handler := http.HandlerFunc(user.Detail)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns user not found in errors", func() {
				Expect(res.Errors[0]).To(MatchRegexp("User not found"))
			})

			It("returns http status code NotFound", func() {
				Expect(statusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("error in service call", func() {
			var statusCode int
			var res errorResult

			BeforeEach(func() {
				req, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/users/%v", id),
					nil)

				if err != nil {
					panic(err)
				}

				req = mux.SetURLVars(req, map[string]string{
					"id": id,
				})

				recorder := httptest.NewRecorder()
				user := NewUser(&fakeUserService{
					user: nil,
					err:  fmt.Errorf("error"),
				})

				handler := http.HandlerFunc(user.Detail)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns internal server error in errors", func() {
				Expect(res.Errors[0]).To(MatchRegexp("Internal server error"))
			})

			It("returns http status code InternalServerError", func() {
				Expect(statusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

type fakeUserService struct {
	user *models.User
	err  error
}

//goland:noinspection GoUnusedParameter
func (srv fakeUserService) Create(
	ctx context.Context,
	request *models.UserCreateRequest) (*models.User, error) {
	return srv.user, srv.err
}

//goland:noinspection GoUnusedParameter
func (srv fakeUserService) Delete(ctx context.Context, id string) error {
	return srv.err
}

//goland:noinspection GoUnusedParameter
func (srv fakeUserService) Detail(
	ctx context.Context,
	id string) (*models.User, error) {
	return srv.user, srv.err
}
