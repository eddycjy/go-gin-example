package xlsx

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	errNilInterface     = errors.New("nil pointer is not a valid argument")
	errNotStructPointer = errors.New("argument must be a pointer to struct")
	errInvalidTag       = errors.New(`invalid tag: must have the format xlsx:idx`)
)

//XLSXUnmarshaler is the interface implemented for types that can unmarshal a Row
//as a representation of themselves.
type XLSXUnmarshaler interface {
	Unmarshal(*Row) error
}

//ReadStruct reads a struct from r to ptr. Accepts a ptr
//to struct. This code expects a tag xlsx:"N", where N is the index
//of the cell to be used. Basic types like int,string,float64 and bool
//are supported
func (r *Row) ReadStruct(ptr interface{}) error {
	if ptr == nil {
		return errNilInterface
	}
	//check if the type implements XLSXUnmarshaler. If so,
	//just let it do the work.
	unmarshaller, ok := ptr.(XLSXUnmarshaler)
	if ok {
		return unmarshaller.Unmarshal(r)
	}
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return errNotStructPointer
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errNotStructPointer
	}
	n := v.NumField()
	for i := 0; i < n; i++ {
		field := v.Type().Field(i)
		idx := field.Tag.Get("xlsx")
		//do a recursive check for the field if it is a struct or a pointer
		//even if it doesn't have a tag
		//ignore if it has a - or empty tag
		isTime := false
		switch {
		case idx == "-":
			continue
		case field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Struct:
			var structPtr interface{}
			if !v.Field(i).CanSet() {
				continue
			}
			if field.Type.Kind() == reflect.Struct {
				structPtr = v.Field(i).Addr().Interface()
			} else {
				structPtr = v.Field(i).Interface()
			}
			//check if the container is a time.Time
			_, isTime = structPtr.(*time.Time)
			if isTime {
				break
			}
			err := r.ReadStruct(structPtr)
			if err != nil {
				return err
			}
			continue
		case len(idx) == 0:
			continue
		}
		pos, err := strconv.Atoi(idx)
		if err != nil {
			return errInvalidTag
		}

		//check if desired position is not out of bounds
		if pos > len(r.Cells)-1 {
			continue
		}
		cell := r.Cells[pos]
		fieldV := v.Field(i)
		//continue if the field is not settable
		if !fieldV.CanSet() {
			continue
		}
		if isTime {
			t, err := cell.GetTime(false)
			if err != nil {
				return err
			}
			if field.Type.Kind() == reflect.Ptr {
				fieldV.Set(reflect.ValueOf(&t))
			} else {
				fieldV.Set(reflect.ValueOf(t))
			}
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			value, err := cell.FormattedValue()
			if err != nil {
				return err
			}
			fieldV.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, err := cell.Int64()
			if err != nil {
				return err
			}
			fieldV.SetInt(value)
		case reflect.Float64:
			value, err := cell.Float()
			if err != nil {
				return err
			}
			fieldV.SetFloat(value)
		case reflect.Bool:
			value := cell.Bool()
			fieldV.SetBool(value)
		}
	}
	value := v.Interface()
	ptr = &value
	return nil
}
