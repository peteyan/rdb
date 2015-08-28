// Copyright 2014 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

package ms

import (
	"math/big"
	"reflect"
	"testing"

	"bitbucket.org/kardianos/rdb"
)

func TestNumber(t *testing.T) {
	defer recoverTest(t)

	cmd := &rdb.Command{
		Sql: `
			select
				bt = @bt, bf = @bf,
				i8 = @i8, i16 = @i16, i64 = @i64,
				bb = @bb,
				dec = @dec,
				fl32 = @fl32,
				fl64 = @fl64
		`,
		Arity: rdb.OneMust,
	}

	var bt, bf bool
	var i8 byte
	var i16 int16
	var bb []byte
	var dec *big.Rat
	var fl32 float32
	var fl64 float64

	var i64 int64 = 1234567

	params := []rdb.Param{
		{Name: "bt", Type: rdb.TypeBool, Value: true},
		{Name: "bf", Type: rdb.TypeBool, Value: false},
		{Name: "i8", Type: rdb.TypeInt8, Value: byte(55)},
		{Name: "i16", Type: rdb.TypeInt16, Value: int16(1234)},
		{Name: "i64", Type: rdb.Integer, Value: i64},
		{Name: "bb", Type: rdb.Binary, Length: 0, Value: []byte{23, 24, 25, 26, 27}},
		{Name: "dec", Type: rdb.TypeDecimal, Precision: 38, Scale: 4, Value: big.NewRat(1234, 100)},
		{Name: "fl32", Type: rdb.TypeFloat32, Value: float32(45.67)},
		{Name: "fl64", Type: rdb.TypeFloat64, Value: float64(89.1011)},
	}

	res := db.Query(cmd, params...)
	defer res.Close()

	res.Scan(&bt, &bf, &i8, &i16, &i64, &bb, &dec, &fl32, &fl64)

	compare := []interface{}{bt, bf, i8, i16, i64, bb, dec, fl32, fl64}

	for i := range compare {
		in := params[i]
		if !reflect.DeepEqual(compare[i], in.Value) {
			t.Errorf("Param %s did not round trip: Want (%v) got (%v)", in.Name, in.Value, compare[i])
		}
	}
}

func TestDecimal(t *testing.T) {
	t.Skip("TODO: Fix large decimals")
	defer recoverTest(t)

	cmd := &rdb.Command{
		Sql: `
			declare @ld decimal(38,6) = @d;

			select d = @ld, s = cast(@ld as varchar(100));
		`,
		Arity: rdb.OneMust,
	}

	var dec *big.Rat
	var sdec string

	dIn := &big.Rat{}
	// dIn.SetString("1.035")
	dIn = big.NewRat(4661225614328463, 4503599627370496)
	params := []rdb.Param{
		{Name: "d", Type: rdb.TypeDecimal, Precision: 38, Scale: 6, Value: dIn},
	}

	res := db.Query(cmd, params...)
	defer res.Close()

	res.Scan(&dec, &sdec)

	if dIn.FloatString(10) != dec.FloatString(10) {
		t.Errorf("D: %v, S: %v, In: %v", dec.FloatString(10), sdec, dIn.FloatString(10))
	}
}

func TestBytesValue(t *testing.T) {
	defer recoverTest(t)

	cmd := &rdb.Command{
		Sql: `
			select @bytesEmpty, @bytesOne
		`,
		Arity: rdb.OneMust,
	}

	bytesEmpty, bytesOne := []byte{}, []byte{01}
	var bytesEmptyOut, bytesOneOut []byte

	params := []rdb.Param{
		{Name: "bytesEmpty", Type: rdb.Binary, Value: bytesEmpty},
		{Name: "bytesOne", Type: rdb.Binary, Value: bytesOne},
	}

	res := db.Query(cmd, params...)
	defer res.Close()

	res.Scan(&bytesEmptyOut, &bytesOneOut)

	compare := []interface{}{bytesEmptyOut, bytesOneOut}

	for i := range compare {
		in := params[i]
		if !reflect.DeepEqual(compare[i], in.Value) {
			t.Errorf("Param %s did not round trip: Want (%v) got (%v)", in.Name, in.Value, compare[i])
		}
	}
}

func TestNullNumbers(t *testing.T) {
	defer recoverTest(t)

	cmd := &rdb.Command{
		Sql: `
			select @decimal
		`,
		Arity: rdb.OneMust,
	}

	params := []rdb.Param{
		{Name: "decimal", Type: rdb.Decimal, Value: nil, Precision: 38, Scale: 6, Null: true},
	}

	res := db.Query(cmd, params...)
	defer res.Close()

	res.Scan()
	val := res.Getx(0)

	if val != nil {
		t.Fatalf("Rat should be nil: %v", val)
	}
}
