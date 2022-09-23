package ogx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"gitee.com/chentanyang/ogx/schema"
)

var errNilModel = errors.New("ogx: Model(nil)")

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

type Model = schema.Model

type rowScanner interface {
	ScanRow(ctx context.Context, rows *sql.Rows) error
}

type TableModel interface {
	Model

	schema.BeforeAppendModelHook
	schema.BeforeScanRowHook
	schema.AfterScanRowHook
	ScanColumn(column string, src interface{}) error

	Table() *schema.Table
	Relation() *schema.Relation

	join(string) *relationJoin
	getJoin(string) *relationJoin
	getJoins() []relationJoin
	addJoin(relationJoin) *relationJoin

	rootValue() reflect.Value
	parentIndex() []int
	mount(reflect.Value)

	updateSoftDeleteField(time.Time) error
}

func newModel(db *DB, dest []interface{}) (Model, error) {
	if len(dest) == 1 {
		return _newModel(db, dest[0], true)
	}

	values := make([]reflect.Value, len(dest))

	for i, el := range dest {
		v := reflect.ValueOf(el)
		if v.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("ogx: Scan(non-pointer %T)", dest)
		}

		v = v.Elem()
		if v.Kind() != reflect.Slice {
			return newScanModel(db, dest), nil
		}

		values[i] = v
	}

	return newSliceModel(db, dest, values), nil
}

func newSingleModel(db *DB, dest interface{}) (Model, error) {
	return _newModel(db, dest, false)
}

func _newModel(db *DB, dest interface{}, scan bool) (Model, error) {
	switch dest := dest.(type) {
	case nil:
		return nil, errNilModel
	case Model:
		return dest, nil
	case sql.Scanner:
		if !scan {
			return nil, fmt.Errorf("ogx: Model(unsupported %T)", dest)
		}
		return newScanModel(db, []interface{}{dest}), nil
	}

	v := reflect.ValueOf(dest)
	if !v.IsValid() {
		return nil, errNilModel
	}
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("ogx: Model(non-pointer %T)", dest)
	}

	if v.IsNil() {
		typ := v.Type().Elem()
		if typ.Kind() == reflect.Struct {
			return newStructTableModel(db, dest, db.Table(typ)), nil
		}
		return nil, fmt.Errorf("ogx: Model(nil %T)", dest)
	}

	v = v.Elem()

	switch v.Kind() {
	case reflect.Map:
		typ := v.Type()
		if err := validMap(typ); err != nil {
			return nil, err
		}
		mapPtr := v.Addr().Interface().(*map[string]interface{})
		return newMapModel(db, mapPtr), nil
	case reflect.Struct:
		if v.Type() != timeType {
			return newStructTableModelValue(db, dest, v), nil
		}
	case reflect.Slice:
		switch elemType := sliceElemType(v); elemType.Kind() {
		case reflect.Struct:
			if elemType != timeType {
				return newSliceTableModel(db, dest, v, elemType), nil
			}
		case reflect.Map:
			if err := validMap(elemType); err != nil {
				return nil, err
			}
			slicePtr := v.Addr().Interface().(*[]map[string]interface{})
			return newMapSliceModel(db, slicePtr), nil
		}
		return newSliceModel(db, []interface{}{dest}, []reflect.Value{v}), nil
	}

	if scan {
		return newScanModel(db, []interface{}{dest}), nil
	}

	return nil, fmt.Errorf("ogx: Model(unsupported %T)", dest)
}

func newTableModelIndex(
	db *DB,
	table *schema.Table,
	root reflect.Value,
	index []int,
	rel *schema.Relation,
) (TableModel, error) {
	typ := typeByIndex(table.Type, index)

	if typ.Kind() == reflect.Struct {
		return &structTableModel{
			db:    db,
			table: table.Dialect().Tables().Get(typ),
			rel:   rel,

			root:  root,
			index: index,
		}, nil
	}

	if typ.Kind() == reflect.Slice {
		structType := indirectType(typ.Elem())
		if structType.Kind() == reflect.Struct {
			m := sliceTableModel{
				structTableModel: structTableModel{
					db:    db,
					table: table.Dialect().Tables().Get(structType),
					rel:   rel,

					root:  root,
					index: index,
				},
			}
			m.init(typ)
			return &m, nil
		}
	}

	return nil, fmt.Errorf("ogx: NewModel(%s)", typ)
}

func validMap(typ reflect.Type) error {
	if typ.Key().Kind() != reflect.String || typ.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("ogx: Model(unsupported %s) (expected *map[string]interface{})",
			typ)
	}
	return nil
}

//------------------------------------------------------------------------------

func isSingleRowModel(m Model) bool {
	switch m.(type) {
	case *mapModel,
		*structTableModel,
		*scanModel:
		return true
	default:
		return false
	}
}
