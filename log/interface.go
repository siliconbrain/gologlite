package log

// Target defines a minimal interface for log targets
type Target interface {
	// Record is used to record events to the log target
	// Events comprise a `message` and set of `fields` providing structured details
	Record(message string, fields FieldSet)
}

// FieldSet defines a minimum interface for a collection of structured log details
type FieldSet interface {
	// ForEachField can be used to apply a function to all fields in the collection
	ForEachField(fn func(name string, value interface{}) (stop bool))
}

// IndexableFieldSet defines an interface that allows looking up fields by name
type IndexableFieldSet interface {
	// LookupFieldByName returns the value associated with the specified name in the field set and `true`, or `nil` and `false`, if there's no such value
	LookupFieldByName(name string) (value interface{}, found bool)
}
