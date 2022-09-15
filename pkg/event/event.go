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

// NewEventFactory returns a pointer to a new and fully configured factory instance,
// and the instance can be used for event generation.
// The arguments are the number of fields (other than the timestamp, which is added
// automatically, and which should therefore not be included in this count),
// the minimum length of field values, and the maximum length of field values.
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

// Init initializes the factory with the actual struct type used for event generation
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

// NewEvent returns a new event as a byte slice
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

// GenerateFieldNames generates a slice of random strings to be used as the names of
// fields.  Each name is at least MIN_FIELDNAME_LEN characters long and at most
// MAX_FIELDNAME_LEN characters long.
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
