package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventCreateUser", func() {
	Describe("Validate", func() {
		Describe("ID", func() {
			Context("empty", func() {
				var err error

				BeforeEach(func() {
					ecu := new(EventCreateUser)
					err = ecu.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("invalid value", func() {
				var err error

				BeforeEach(func() {
					ecu := EventCreateUser{ID: "foo-bar"}
					err = ecu.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("valid value", func() {
				var err error

				BeforeEach(func() {
					ecu := EventCreateUser{ID: "7b5a3155-7a73-42de-b87e-23f50a10180a"}
					err = ecu.Validate()
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})

var _ = Describe("EventCreateRequest", func() {
	Describe("Validate", func() {
		Describe("User", func() {
			Context("nil", func() {
				var err error

				BeforeEach(func() {
					ecr := EventCreateRequest{
						Consents: &[]Consent{
							{ID: ConsentEmail},
						},
					}
					err = ecr.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("not nil", func() {
				var err error

				BeforeEach(func() {
					ecr := EventCreateRequest{
						User: &EventCreateUser{ID: "7b5a3155-7a73-42de-b87e-23f50a10180a"},
						Consents: &[]Consent{
							{ID: ConsentSMS},
						},
					}
					err = ecr.Validate()
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})
			})
		})

		Describe("Consents", func() {
			Context("nil", func() {
				var err error

				BeforeEach(func() {
					ecr := EventCreateRequest{
						User: &EventCreateUser{ID: "7b5a3155-7a73-42de-b87e-23f50a10180a"},
					}
					err = ecr.Validate()
				})

				It("returns error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("not nil", func() {
				var err error

				BeforeEach(func() {
					ecr := EventCreateRequest{
						User: &EventCreateUser{ID: "7b5a3155-7a73-42de-b87e-23f50a10180a"},
						Consents: &[]Consent{
							{ID: ConsentEmail},
						},
					}
					err = ecr.Validate()
				})

				It("does not return error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
