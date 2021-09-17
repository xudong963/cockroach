// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package rowenc

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/require"
)

func genColumnType() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		columnType := RandColumnType(genParams.Rng)
		return gopter.NewGenResult(columnType, gopter.NoShrinker)
	}
}

func genRandomArrayType() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		arrType := RandArrayType(genParams.Rng)
		return gopter.NewGenResult(arrType, gopter.NoShrinker)
	}
}

func genDatum() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		return gopter.NewGenResult(RandDatum(genParams.Rng, RandColumnType(genParams.Rng),
			false), gopter.NoShrinker)
	}
}

func genDatumWithType(columnType interface{}) gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		datum := RandDatum(genParams.Rng, columnType.(*types.T), false)
		return gopter.NewGenResult(datum, gopter.NoShrinker)
	}
}

func genArrayDatumWithType(arrTyp interface{}) gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		// Mark the array contents to have a 1 in 10 chance of being null.
		datum := RandArray(genParams.Rng, arrTyp.(*types.T), 10)
		return gopter.NewGenResult(datum, gopter.NoShrinker)
	}
}

func genEncodingDirection() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		return gopter.NewGenResult(
			encoding.Direction((genParams.Rng.Int()%int(encoding.Descending))+1),
			gopter.NoShrinker)
	}
}

func hasKeyEncoding(typ *types.T) bool {
	// Only some types are round-trip key encodable.
	switch typ.Family() {
	case types.JsonFamily, types.CollatedStringFamily, types.TupleFamily, types.DecimalFamily,
		types.GeographyFamily, types.GeometryFamily:
		return false
	case types.ArrayFamily:
		return hasKeyEncoding(typ.ArrayContents())
	}
	return true
}

func TestEncodeTableValue(t *testing.T) {
	a := &DatumAlloc{}
	ctx := tree.NewTestingEvalContext(cluster.MakeTestingClusterSettings())
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(parameters)
	var scratch []byte
	properties.Property("roundtrip", prop.ForAll(
		func(d tree.Datum) string {
			b, err := EncodeTableValue(nil, 0, d, scratch)
			if err != nil {
				return "error: " + err.Error()
			}
			newD, leftoverBytes, err := DecodeTableValue(a, d.ResolvedType(), b)
			if len(leftoverBytes) > 0 {
				return "Leftover bytes"
			}
			if err != nil {
				return "error: " + err.Error()
			}
			if newD.Compare(ctx, d) != 0 {
				return "unequal"
			}
			return ""
		},
		genDatum(),
	))
	properties.TestingRun(t)
}

func TestEncodeTableKey(t *testing.T) {
	a := &DatumAlloc{}
	ctx := tree.NewTestingEvalContext(cluster.MakeTestingClusterSettings())
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(parameters)
	roundtripDatum := func(d tree.Datum, dir encoding.Direction) string {
		b, err := EncodeTableKey(nil, d, dir)
		if err != nil {
			return "error: " + err.Error()
		}
		newD, leftoverBytes, err := DecodeTableKey(a, d.ResolvedType(), b, dir)
		if len(leftoverBytes) > 0 {
			return "Leftover bytes"
		}
		if err != nil {
			return "error: " + err.Error()
		}
		if newD.Compare(ctx, d) != 0 {
			return "unequal"
		}
		return ""
	}
	properties.Property("roundtrip", prop.ForAll(
		roundtripDatum,
		genColumnType().
			SuchThat(hasKeyEncoding).
			FlatMap(genDatumWithType, reflect.TypeOf((*tree.Datum)(nil)).Elem()),
		genEncodingDirection(),
	))

	// Also run the property on arrays possibly containing NULL values.
	// The random generator in the property above does not generate NULLs.
	properties.Property("roundtrip-arrays", prop.ForAll(
		roundtripDatum,
		genRandomArrayType().
			SuchThat(hasKeyEncoding).
			FlatMap(genArrayDatumWithType, reflect.TypeOf((*tree.Datum)(nil)).Elem()),
		genEncodingDirection(),
	))

	generateAndCompareDatums := func(datums []tree.Datum, dir encoding.Direction) string {
		d1 := datums[0]
		d2 := datums[1]
		b1, err := EncodeTableKey(nil, d1, dir)
		if err != nil {
			return "error: " + err.Error()
		}
		b2, err := EncodeTableKey(nil, d2, dir)
		if err != nil {
			return "error: " + err.Error()
		}

		expectedCmp := d1.Compare(ctx, d2)
		cmp := bytes.Compare(b1, b2)

		if expectedCmp == 0 {
			if cmp != 0 {
				return fmt.Sprintf("equal inputs produced inequal outputs: \n%v\n%v", b1, b2)
			}
			// If the inputs are equal and so are the outputs, no more checking to do.
			return ""
		}

		cmpsMatch := expectedCmp == cmp
		dirIsAscending := dir == encoding.Ascending

		if cmpsMatch != dirIsAscending {
			return fmt.Sprintf("non-order preserving encoding: \n%v\n%v", b1, b2)
		}
		return ""
	}

	properties.Property("order-preserving", prop.ForAll(
		generateAndCompareDatums,
		// For each column type, generate two datums of that type.
		genColumnType().
			SuchThat(hasKeyEncoding).
			FlatMap(
				func(t interface{}) gopter.Gen {
					colTyp := t.(*types.T)
					return gopter.CombineGens(
						genDatumWithType(colTyp),
						genDatumWithType(colTyp))
				}, reflect.TypeOf([]interface{}{})).
			Map(func(datums []interface{}) []tree.Datum {
				ret := make([]tree.Datum, len(datums))
				for i, d := range datums {
					ret[i] = d.(tree.Datum)
				}
				return ret
			}),
		genEncodingDirection(),
	))

	// Also run the property on arrays possibly containing NULL values.
	// The random generator in the property above does not generate NULLs.
	properties.Property("order-preserving-arrays", prop.ForAll(
		generateAndCompareDatums,
		// For each column type, generate two datums of that type.
		genRandomArrayType().
			SuchThat(hasKeyEncoding).
			FlatMap(
				func(t interface{}) gopter.Gen {
					colTyp := t.(*types.T)
					return gopter.CombineGens(
						genArrayDatumWithType(colTyp),
						genArrayDatumWithType(colTyp))
				}, reflect.TypeOf([]interface{}{})).
			Map(func(datums []interface{}) []tree.Datum {
				ret := make([]tree.Datum, len(datums))
				for i, d := range datums {
					ret[i] = d.(tree.Datum)
				}
				return ret
			}),
		genEncodingDirection(),
	))

	properties.TestingRun(t)
}

func TestSkipTableKey(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(parameters)
	properties.Property("correctness", prop.ForAll(
		func(d tree.Datum, dir encoding.Direction) string {
			b, err := EncodeTableKey(nil, d, dir)
			if err != nil {
				return "error: " + err.Error()
			}
			res, err := SkipTableKey(b)
			if err != nil {
				return "error: " + err.Error()
			}
			if len(res) != 0 {
				fmt.Println(res, len(res), d.ResolvedType(), d.ResolvedType().Family())
				return "expected 0 bytes remaining"
			}
			return ""
		},
		genColumnType().
			SuchThat(hasKeyEncoding).FlatMap(genDatumWithType, reflect.TypeOf((*tree.Datum)(nil)).Elem()),
		genEncodingDirection(),
	))
	properties.TestingRun(t)
}

func TestMarshalColumnValueRoundtrip(t *testing.T) {
	a := &DatumAlloc{}
	ctx := tree.NewTestingEvalContext(cluster.MakeTestingClusterSettings())
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(parameters)

	properties.Property("roundtrip",
		prop.ForAll(
			func(typ *types.T) string {
				d, ok := genDatumWithType(typ).Sample()
				if !ok {
					return "error generating datum"
				}
				datum := d.(tree.Datum)
				desc := descpb.ColumnDescriptor{
					Type: typ,
				}
				value, err := MarshalColumnValue(&desc, datum)
				if err != nil {
					return "error marshaling: " + err.Error()
				}
				outDatum, err := UnmarshalColumnValue(a, typ, value)
				if err != nil {
					return "error unmarshaling: " + err.Error()
				}
				if datum.Compare(ctx, outDatum) != 0 {
					return fmt.Sprintf("datum didn't roundtrip.\ninput: %v\noutput: %v", datum, outDatum)
				}
				return ""
			},
			genColumnType(),
		),
	)
	properties.TestingRun(t)
}

// TestDecodeTableKeyOutOfRangeTimestamp deliberately tests out of range timestamps
// can still be decoded from disk. See #46973.
func TestDecodeTableKeyOutOfRangeTimestamp(t *testing.T) {
	for _, d := range []tree.Datum{
		&tree.DTimestamp{Time: timeutil.Unix(-9223372036854775808, 0).In(time.UTC)},
		&tree.DTimestampTZ{Time: timeutil.Unix(-9223372036854775808, 0).In(time.UTC)},
	} {
		for _, dir := range []encoding.Direction{encoding.Ascending, encoding.Descending} {
			t.Run(fmt.Sprintf("%s/direction:%d", d.String(), dir), func(t *testing.T) {
				encoded, err := EncodeTableKey([]byte{}, d, dir)
				require.NoError(t, err)
				a := &DatumAlloc{}
				decoded, _, err := DecodeTableKey(a, d.ResolvedType(), encoded, dir)
				require.NoError(t, err)
				require.Equal(t, d, decoded)
			})
		}
	}
}

// TestDecodeTableValueOutOfRangeTimestamp deliberately tests out of range timestamps
// can still be decoded from disk. See #46973.
func TestDecodeTableValueOutOfRangeTimestamp(t *testing.T) {
	for _, d := range []tree.Datum{
		&tree.DTimestamp{Time: timeutil.Unix(-9223372036854775808, 0).In(time.UTC)},
		&tree.DTimestampTZ{Time: timeutil.Unix(-9223372036854775808, 0).In(time.UTC)},
	} {
		t.Run(d.String(), func(t *testing.T) {
			var b []byte
			colID := descpb.ColumnID(1)
			encoded, err := EncodeTableValue(b, colID, d, []byte{})
			require.NoError(t, err)
			a := &DatumAlloc{}
			decoded, _, err := DecodeTableValue(a, d.ResolvedType(), encoded)
			require.NoError(t, err)
			require.Equal(t, d, decoded)
		})
	}
}

// This test ensures that decoding a tuple value with a specific, labeled tuple
// type preserves the labels.
func TestDecodeTupleValueWithType(t *testing.T) {
	tupleType := types.MakeLabeledTuple([]*types.T{types.Int, types.String}, []string{"a", "b"})
	datum := tree.NewDTuple(tupleType, tree.NewDInt(tree.DInt(1)), tree.NewDString("foo"))
	buf, err := rowenc.EncodeTableValue(nil, descpb.ColumnID(encoding.NoColumnID), datum, nil)
	if err != nil {
		t.Fatal(err)
	}
	da := rowenc.DatumAlloc{}
	var decoded tree.Datum
	decoded, _, err = rowenc.DecodeTableValue(&da, tupleType, buf)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, decoded, datum)
}
