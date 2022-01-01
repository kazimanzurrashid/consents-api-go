package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consent", func() {
	Describe("Validate", func() {
		Describe("ID", func() {
			Context("empty", func() {
				var err error

				BeforeEach(func() {
					c := new(Consent)
					err = c.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("invalid value", func() {
				var err error

				BeforeEach(func() {
					c := Consent{ID: "foo-bar"}
					err = c.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("valid value", func() {
				var err error

				BeforeEach(func() {
					c := Consent{ID: ConsentEmail}
					err = c.Validate()
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
