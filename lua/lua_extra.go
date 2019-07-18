package lua

import (
	"fmt"
	"github.com/camsiabor/qcom/util"
	"reflect"
)

func (L *State) TableSetValue(tableIndex int, key string, val interface{}) (err error) {

	defer func() {
		if err == nil {
			L.SetField(tableIndex, key)
		}
	}()

	if val == nil {
		L.PushNil()
		return
	}

	if gofunc, ok := val.(LuaGoFunction); ok {
		L.PushGoFunction(gofunc)
		return
	}

	if str, ok := val.(string); ok {
		L.PushString(str)
		return
	}

	var vref = reflect.ValueOf(val)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var n = util.AsInt64(val, 0)
		L.PushInteger(n)
	case reflect.Float32, reflect.Float64:
		var n = util.AsFloat64(val, 0)
		L.PushNumber(n)
	case reflect.Bool:
		var b = val.(bool)
		L.PushBoolean(b)
	case reflect.Map:
		err = fmt.Errorf("map not support")
	case reflect.Ptr, reflect.Chan, reflect.Array, reflect.Slice:
		var uptr = vref.Pointer()
		L.PushInteger(int64(uptr))
	default:
		L.PushGoStruct(val)
	}
	return nil
}

func (L *State) TableRegister(table string, name string, val interface{}) error {
	var _, err = L.GetTableByName(table, true)
	if err != nil {
		return err
	}
	err = L.TableSetValue(-2, name, val)
	L.Pop(1)
	return nil
}

func (L *State) TableRegisters(table string, funcs map[string]interface{}) error {
	if funcs == nil {
		return fmt.Errorf("no function is set")
	}
	var _, err = L.GetTableByName(table, true)
	if err != nil {
		return err
	}
	for key, val := range funcs {
		err = L.TableSetValue(-2, key, val)
		if err != nil {
			break
		}
	}
	L.Pop(1)
	return err
}
