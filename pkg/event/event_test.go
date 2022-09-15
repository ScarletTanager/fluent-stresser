package event_test

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"time"
	"unicode/utf8"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/fluent-stresser/pkg/event"
)

var _ = Describe("Event", func() {
	Describe("GenerateFieldNames", func() {
		var (
			numFields int
		)

		BeforeEach(func() {
			numFields = rand.Intn(100)
		})

		It("Generates the correct number of names", func() {
			names := event.GenerateFieldNames(numFields)
			Expect(names).To(HaveLen(numFields))
		})

		It("Generates names of the correct length", func() {
			names := event.GenerateFieldNames(numFields)
			for _, name := range names {
				Expect(utf8.RuneCountInString(name)).To(BeNumerically("<=", event.MAX_FIELDNAME_LEN))
				Expect(utf8.RuneCountInString(name)).To(BeNumerically(">=", event.MIN_FIELDNAME_LEN))
			}
		})
	})

	Describe("EventFactory", func() {
		var (
			numFields, maxFieldLen, minFieldLen int
		)

		BeforeEach(func() {
			numFields = 10
			minFieldLen = 16
			maxFieldLen = 25
		})

		Describe("NewEventFactory", func() {
			It("Creates a new factory", func() {
				Expect(event.NewEventFactory(numFields, minFieldLen, maxFieldLen)).NotTo(BeNil())
			})
		})

		Describe("NewEvent", func() {
			var (
				factory *event.EventFactory
			)

			BeforeEach(func() {
				factory = event.NewEventFactory(numFields, minFieldLen, maxFieldLen)
				Expect(factory).NotTo(BeNil())
			})

			It("Generates an event in json", func() {
				ev := factory.NewEvent()
				Expect(ev).NotTo(BeNil())

				target := reflect.New(factory.EventType).Elem().Addr().Interface()

				Expect(json.Unmarshal(ev, target)).NotTo(HaveOccurred())
			})

			It("Adds a timestamp as the first field", func() {
				ev := factory.NewEvent()
				target := reflect.New(factory.EventType).Elem()
				json.Unmarshal(ev, target.Addr().Interface())

				_, err := time.Parse(time.RFC3339, target.FieldByName("Timestamp").String())
				Expect(err).NotTo(HaveOccurred())

				_, err := time.Parse(time.RFC3339, target.Field(0).String())
				Expect(err).NotTo(HaveOccurred())
			})

			It("Populates each non-Timestamp field with a random string of the correct length", func() {
				ev := factory.NewEvent()

				target := reflect.New(factory.EventType).Elem()
				json.Unmarshal(ev, target.Addr().Interface())

				// Field(0) is the timestamp, skip it
				for i := 1; i < target.NumField(); i++ {
					s := target.Field(i).String()
					Expect(utf8.RuneCountInString(s)).To(BeNumerically(">=", minFieldLen))
					Expect(utf8.RuneCountInString(s)).To(BeNumerically("<=", maxFieldLen))
				}
			})
		})
	})
})
