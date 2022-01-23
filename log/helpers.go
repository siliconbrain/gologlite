package log

// LookupFieldByName returns the value associated with the specified name in the specified field set and `true`, or `nil` and `false` if there's no such value
// It treats absent (i.e. `nil`) field sets as empty and tries to use the IndexableFieldSet interface when available before searching through the set using ForEach
func LookupFieldByName(fs FieldSet, name string) (value interface{}, found bool) {
	if fs == nil {
		return
	}
	if indexable, ok := fs.(IndexableFieldSet); ok {
		return indexable.LookupFieldByName(name)
	}

	fs.ForEachField(func(n string, v interface{}) (stop bool) {
		if name == n {
			value, found = v, true
		}
		return found
	})

	return
}

// Event records an event to `target` with `message` and `fields`
// It handles an absent (i.e. `nil`) target by discarding the event
func Event(target Target, message string, fields ...FieldSet) {
	if target != nil {
		target.Record(message, LittleEndianFieldSetList(fields))
	}
}

// LittleEndianFieldSetList stores its field sets from least significant to most significant
// This means that if fields with the same name are present in multiple sets, only the associated value in the most significant set is returned by ForEach and Lookup
type LittleEndianFieldSetList []FieldSet

func (fsl LittleEndianFieldSetList) ForEachField(fn func(name string, value interface{}) bool) {
	names := map[string]bool{}
	for i := len(fsl); i > 0; i-- {
		fs := fsl[i-1]
		var stop bool
		fs.ForEachField(func(name string, value interface{}) bool {
			if !names[name] {
				stop = fn(name, value)
				names[name] = true
			}
			return stop
		})
		if stop {
			return
		}
	}
}

func (fsl LittleEndianFieldSetList) LookupFieldByName(name string) (value interface{}, found bool) {
	for i := len(fsl); i > 0; i-- {
		fs := fsl[i-1]
		if val, fnd := LookupFieldByName(fs, name); fnd {
			value, found = val, fnd
			break
		}
	}
	return
}

// FieldMap implements the FieldSet and IndexableFieldSet interfaces using a `map`
type FieldMap map[string]interface{}

func (f FieldMap) ForEachField(fn func(name string, value interface{}) (stop bool)) {
	for name, value := range f {
		if fn(name, value) {
			return
		}
	}
}

func (f FieldMap) LookupFieldByName(name string) (value interface{}, found bool) {
	value, found = f[name]
	return
}

// CollapseFieldSets returns a union of all fields in the sets with the latest values kept for duplicate fields
func CollapseFieldSets(fss ...FieldSet) (res FieldMap) {
	res = make(FieldMap)
	for _, fs := range fss {
		fs.ForEachField(func(name string, value interface{}) (stop bool) {
			res[name] = value
			return
		})
	}
	return
}
