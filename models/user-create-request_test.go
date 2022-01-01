package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserCreateRequest", func() {
	Describe("Validate", func() {
		Describe("Email", func() {
			Context("empty", func() {
				var err error

				BeforeEach(func() {
					ucr := new(UserCreateRequest)
					err = ucr.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("invalid value", func() {
				var err error

				BeforeEach(func() {
					ucr := UserCreateRequest{Email: "foo-bar"}
					err = ucr.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("valid value", func() {
				var err error

				BeforeEach(func() {
					ucr := UserCreateRequest{Email: "user@example.com"}
					err = ucr.Validate()
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
