package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kazimanzurrashid/consents-api-go/models"
)

var _ = Describe("Event", func() {
	Describe("Create", func() {
		Context("success", func() {
			var statusCode int

			BeforeEach(func() {
				var payload bytes.Buffer

				err := json.NewEncoder(&payload).Encode(models.EventCreateRequest{
					User: &models.EventCreateUser{
						ID: "7b5a3155-7a73-42de-b87e-23f50a10180a",
					},
					Consents: &[]models.Consent{
						{
							ID:      models.ConsentEmail,
							Enabled: true,
						},
						{
							ID:      models.ConsentSMS,
							Enabled: false,
						},
					},
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					"/events",
					&payload)

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				event := NewEvent(&fakeEventService{err: nil})

				handler := http.HandlerFunc(event.Create)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code
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
					"/events",
					strings.NewReader("malformed json"))

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				event := NewEvent(&fakeEventService{err: nil})

				handler := http.HandlerFunc(event.Create)
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

				err := json.NewEncoder(&payload).Encode(models.EventCreateRequest{
					User: &models.EventCreateUser{
						ID: "foo",
					},
					Consents: &[]models.Consent{
						{
							ID:      "bar",
							Enabled: true,
						},
						{
							ID:      "baz",
							Enabled: false,
						},
					},
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					"/events",
					&payload)

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				event := NewEvent(&fakeEventService{err: nil})

				handler := http.HandlerFunc(event.Create)
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

				err := json.NewEncoder(&payload).Encode(models.EventCreateRequest{
					User: &models.EventCreateUser{
						ID: "7b5a3155-7a73-42de-b87e-23f50a10180a",
					},
					Consents: &[]models.Consent{
						{
							ID:      models.ConsentEmail,
							Enabled: true,
						},
						{
							ID:      models.ConsentSMS,
							Enabled: false,
						},
					},
				})

				if err != nil {
					panic(err)
				}

				req, err := http.NewRequest(
					http.MethodPost,
					"/events",
					&payload)

				if err != nil {
					panic(err)
				}

				recorder := httptest.NewRecorder()
				event := NewEvent(&fakeEventService{err: fmt.Errorf("error")})

				handler := http.HandlerFunc(event.Create)
				handler.ServeHTTP(recorder, req)

				statusCode = recorder.Code

				err = json.NewDecoder(recorder.Body).Decode(&res)

				if err != nil {
					panic(err)
				}
			})

			It("returns invalid request in errors", func() {
				Expect(res.Errors[0]).To(MatchRegexp("Invalid request"))
			})

			It("returns http status code UnprocessableEntity", func() {
				Expect(statusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})
	})
})

type fakeEventService struct {
	err error
}

func (srv *fakeEventService) Create(
	_ context.Context,
	_ *models.EventCreateRequest) error {
	return srv.err
}
