package event

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"time"
)

const (
	MIN_FIELDNAME_LEN = 5
	MAX_FIELDNAME_LEN = 15
)

func NewEventFactory(numFields, minFieldLen, maxFieldLen int) *EventFactory {
	factory := &EventFactory{
		NumFields:   numFields,
		MinFieldLen: minFieldLen,
		MaxFieldLen: maxFieldLen,
	}

	factory.Init()

	return factory
}

// EventFactory instances are used to create new events.
type EventFactory struct {
	NumFields   int
	MinFieldLen int
	MaxFieldLen int
	EventType   reflect.Type
}

func (f *EventFactory) Init() {
	names := GenerateFieldNames(f.NumFields)
	fields := make([]reflect.StructField, len(names)+1)
	for i := range fields {
		if i == 0 {
			fields[i] = reflect.StructField{
				Name: "Timestamp",
				Type: reflect.TypeOf(string("")),
			}
		} else {
			fields[i] = reflect.StructField{
				Name: names[i-1],
				Type: reflect.TypeOf(string("")),
			}
		}
	}

	f.EventType = reflect.StructOf(fields)
}

func (f *EventFactory) NewEvent() []byte {
	alphaNum := []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	event := reflect.New(f.EventType).Elem()
	ts, _ := time.Now().UTC().MarshalText()
	event.FieldByName("Timestamp").SetString(string(ts))

	fc := event.NumField()
	for i := 1; i < fc; i++ {
		l := rand.Intn(f.MaxFieldLen-f.MinFieldLen) + f.MinFieldLen
		s := make([]rune, l)
		for pos := 0; pos < l; pos++ {
			s[pos] = alphaNum[rand.Intn(len(alphaNum))]
		}
		event.Field(i).SetString(string(s))
	}

	thing := event.Addr().Interface()

	b, _ := json.Marshal(thing)

	return b
}

func GenerateFieldNames(numFields int) []string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyz")
	upper := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	names := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		s := make([]rune, rand.Intn(MAX_FIELDNAME_LEN-MIN_FIELDNAME_LEN)+MIN_FIELDNAME_LEN)
		s[0] = upper[rand.Intn(len(upper))]
		for pos := range s[:len(s)-1] {
			s[pos+1] = alphabet[rand.Intn(len(alphabet))]
		}
		names[i] = string(s)
	}

	return names
}
