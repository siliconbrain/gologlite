package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	t.Run("with absent target", func(t *testing.T) {
		require.NotPanics(t, func() {
			Event(nil, "the quick brown fox...", FieldMap{"jumped over": "the lazy dog"})
		})
	})
	t.Run("records event", func(t *testing.T) {
		target := &logTarget{}
		message := "the quick brown fox..."
		fields := FieldMap{"jumped over": "the lazy dog"}
		expected := []logRecord{
			{
				message: message,
				fields:  fields,
			},
		}
		Event(target, message, fields)
		require.Equal(t, len(expected), len(target.records))
		for i := range expected {
			expectedRecord := expected[i]
			record := target.records[i]
			require.Equal(t, expectedRecord.message, record.message)
			require.Equal(t, CollapseFieldSets(expectedRecord.fields), CollapseFieldSets(record.fields))
		}
	})
}

func TestLookupFieldByName(t *testing.T) {
	t.Run("with absent field set", func(t *testing.T) {
		value, found := LookupFieldByName(nil, "field1")
		require.Nil(t, value)
		require.False(t, found)
	})
	t.Run("with indexable field set", func(t *testing.T) {
		fieldName, fieldValue := "fieldName", 42
		value, found := LookupFieldByName(indexableOnlyFieldSet{fieldName: fieldValue}, fieldName)
		require.Equal(t, fieldValue, value)
		require.True(t, found)
	})
	t.Run("with simple field set", func(t *testing.T) {
		fieldName, fieldValue := "fieldName", 42
		value, found := LookupFieldByName(iterableOnlyFieldSet{fieldName: fieldValue}, fieldName)
		require.Equal(t, fieldValue, value)
		require.True(t, found)
	})
}

func TestLittleEndianFieldSetList_ForEachField(t *testing.T) {
	testCases := map[string]struct {
		fieldSetList   []FieldSet
		expectedFields FieldMap
	}{
		"with no field sets": {
			fieldSetList:   nil,
			expectedFields: FieldMap{},
		},
		"with single field set": {
			fieldSetList: []FieldSet{
				FieldMap{
					"a": 1,
					"b": 1,
				},
			},
			expectedFields: FieldMap{
				"a": 1,
				"b": 1,
			},
		},
		"with multiple field sets": {
			fieldSetList: []FieldSet{
				FieldMap{
					"a": 1,
					"b": 1,
				},
				FieldMap{
					"b": 2,
					"c": 2,
				},
				FieldMap{
					"c": 3,
					"d": 3,
				},
			},
			expectedFields: FieldMap{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 3,
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			fields := FieldMap{}
			LittleEndianFieldSetList(testCase.fieldSetList).ForEachField(func(name string, value interface{}) bool {
				fields[name] = value
				return false
			})
			require.Equal(t, testCase.expectedFields, fields)
		})
	}
}

func TestLittleEndianFieldSetList_LookupFieldByName(t *testing.T) {
	fieldName, fieldValue := "fieldName", 42
	testCases := map[string]struct {
		fieldSetList  []FieldSet
		expectedValue interface{}
		expectedFound bool
	}{
		"with no field sets": {
			fieldSetList:  nil,
			expectedValue: nil,
			expectedFound: false,
		},
		"with single field set": {
			fieldSetList: []FieldSet{
				FieldMap{
					fieldName: fieldValue,
				},
			},
			expectedValue: fieldValue,
			expectedFound: true,
		},
		"with multiple field sets, field in multiple sets": {
			fieldSetList: []FieldSet{
				FieldMap{
					fieldName: fieldValue - 1,
				},
				FieldMap{
					fieldName: fieldValue,
				},
			},
			expectedValue: fieldValue,
			expectedFound: true,
		},
		"with multiple field sets, field in first set": {
			fieldSetList: []FieldSet{
				FieldMap{
					fieldName: fieldValue,
				},
				FieldMap{},
			},
			expectedValue: fieldValue,
			expectedFound: true,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			value, found := LittleEndianFieldSetList(testCase.fieldSetList).LookupFieldByName(fieldName)
			require.Equal(t, testCase.expectedValue, value)
			require.Equal(t, testCase.expectedFound, found)
		})
	}
}

func TestCollapseFieldSets(t *testing.T) {
	fields := CollapseFieldSets(
		FieldMap{
			"a": 1,
			"b": 1,
		},
		FieldMap{
			"b": 2,
			"c": 2,
		},
		FieldMap{
			"c": 3,
			"d": 3,
		},
	)
	expected := FieldMap{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 3,
	}
	require.Equal(t, expected, fields)
}

type logTarget struct {
	records []logRecord
}

func (t *logTarget) Record(message string, fields FieldSet) {
	t.records = append(t.records, logRecord{
		message: message,
		fields:  fields,
	})
}

type logRecord struct {
	message string
	fields  FieldSet
}

type indexableOnlyFieldSet FieldMap

func (fs indexableOnlyFieldSet) LookupFieldByName(name string) (interface{}, bool) {
	return FieldMap(fs).LookupFieldByName(name)
}

func (indexableOnlyFieldSet) ForEachField(func(string, interface{}) bool) {
	panic("not implemented for testing purposes")
}

type iterableOnlyFieldSet FieldMap

func (fs iterableOnlyFieldSet) ForEachField(fn func(string, interface{}) bool) {
	FieldMap(fs).ForEachField(fn)
}
