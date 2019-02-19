package geodesy_test

import (
	"github.com/recombinant/go-geodesy"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  */
/*  Geodesy Test Harness - dms                                        (c) Chris Veness 2014-2017  */
/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -  */

type resultLookup struct {
	s string
	f float64
}

func TestPassParseDMS(t *testing.T) {
	dataSliceZero := []resultLookup{
		{"0", 0.0},
		{"0°", 0.0},
		{"0.0", 0.0},
		{"0.0°", 0.0},
		{"000", 0.0},
		{"000°", 0.0},
		{"000.0", 0.0},
		{"000.0°", 0.0},

		{"0 0", 0.0},
		{"0°0′", 0.0},
		{"0°0'", 0.0},
		{"0° 0′", 0.0}, // minute sign
		{"0° 0'", 0.0}, // single quote
		{"0 0.0", 0.0},
		{"0°0.0′", 0.0},
		{"0°0.0'", 0.0},
		{"0° 0.0′", 0.0},
		{"0° 0.0'", 0.0},
		{"000 00", 0.0},
		{"000°00′", 0.0},
		{"000°00'", 0.0},
		{"000° 00′", 0.0}, // minute sign
		{"000° 00'", 0.0}, // single quote
		{"000 00.0", 0.0},
		{"000°00.0′", 0.0},
		{"000°00.0'", 0.0},
		{"000° 00.0′", 0.0},
		{"000° 00.0'", 0.0},

		{"0 0 0", 0.0},
		{"0°0′0″", 0.0},
		{"0°0'0\"", 0.0},
		{"0° 0′ 0″", 0.0},  // minute sign
		{"0° 0' 0\"", 0.0}, // single quote
		{"0 0 0.0", 0.0},
		{"0°0′0.0″", 0.0},
		{"0°0'0.0\"", 0.0},
		{"0° 0′ 0″", 0.0},
		{"0° 0' 0\"", 0.0},
		{"000 00 00", 0.0},
		{"000°00′00″", 0.0},
		{"000°00'00", 0.0},
		{"000° 00′ 00″", 0.0},  // minute sign
		{"000° 00' 00\"", 0.0}, // single quote
		{"000 00 00.0", 0.0},
		{"000°00′00.0″", 0.0},
		{"000°00'00.0\"", 0.0},
		{"000° 00′ 00.0″", 0.0},
		{"000° 00' 00.0\"", 0.0},
	}
	dataSliceValue := []resultLookup{
		{"45.76260", 45.76260},
		{"45.76260°", 45.76260},
		{"45°45.756′", 45.76260},
		{"45° 45.756′", 45.76260},
		{"45 45.756", 45.76260},
		{"45°45′45.36″", 45.76260},
		{`45º45'45.36"`, 45.76260},
		{"45°\u202f45’\u202f45.36”", 45.76260}, // U+202F &#8239; NARROW NO-BREAK SPACE (as escape)
		{"45°\u202f45′\u202f45.36″", 45.76260},
		{`45º 45' 45.36"`, 45.76260}, // U+202F &#8239; NARROW NO-BREAK SPACE (as character)
		{"45° 45’ 45.36”", 45.76260},
		{"45° 45′ 45.36″", 45.76260},
		{`45º 45' 45.36"`, 45.76260}, // U+00BA, &#186; MASCULINE ORDINAL INDICATOR
		{"45° 45’ 45.36”", 45.76260},
	}
	dataSliceOutOfRange := []resultLookup{
		// Out of range (is tested both positive and negative)
		{"185", 185.0},
		{"365", 365.0},
	}

	t.Run("Parse zero",
		func(t *testing.T) {
			variations(&dataSliceZero, t)
		})
	t.Run("Parse number",
		func(t *testing.T) {
			variations(&dataSliceValue, t)
		})
	t.Run("Parse out of range",
		func(t *testing.T) {
			variations(&dataSliceOutOfRange, t)
		})

}

func variations(dataSlice *[]resultLookup, t *testing.T) {
	decoratorSlice := []struct {
		prefix, postfix string
		multiplier      float64
	}{
		{"+", "", 1.0},
		{"-", "", -1.0},
		{"−", "", -1.0}, // math minus sign 0x2212
		{"", "N", 1.0},
		{"", "S", -1.0},
		{"", "E", 1.0},
		{"", "W", -1.0},
		{"", " N", 1.0},
		{"", " S", -1.0},
		{"", " E", 1.0},
		{"", " W", -1.0},
	}
	for _, data := range *dataSlice {
		// parse original
		assert.Equal(t, geodesy.ParseDMS(data.s), data.f, data.s)

		// add some decoration -, N, S, E, W
		for _, decorator := range decoratorSlice {
			dataString := decorator.prefix + data.s + decorator.postfix
			result := data.f
			result *= decorator.multiplier

			assert.Equal(t, geodesy.ParseDMS(dataString), result, dataString)

			// leading & trailing, space & spaces
			for _, spaces := range []string{" ", "  "} {
				s1 := spaces + dataString
				assert.Equal(t, geodesy.ParseDMS(s1), result, s1)
				s2 := dataString + spaces
				assert.Equal(t, geodesy.ParseDMS(s2), result, s2)
				s3 := spaces + dataString + spaces
				assert.Equal(t, geodesy.ParseDMS(s3), result, s3)
			}
		}
	}
}

func TestFailParseDMS(t *testing.T) {
	assert.True(t, math.IsNaN(geodesy.ParseDMS("0 0 0 0")))
	assert.True(t, math.IsNaN(geodesy.ParseDMS("xxx")))
}

func TestToDMS(t *testing.T) {
	t.Run("toDMS zero",
		func(t *testing.T) {
			assert.Equal(t, *geodesy.ToDMS1(0), "000°00′00″")
			assert.Equal(t, *geodesy.ToDMS2(0, geodesy.FmtD), "000.0000°")
			assert.Equal(t, *geodesy.ToDMS2(0, geodesy.FmtDM), "000°00.00′")
			assert.Equal(t, *geodesy.ToDMS2(0, geodesy.FmtDMS), "000°00′00″")
			assert.Equal(t, *geodesy.ToDMS3(0, geodesy.FmtD, 0), "000°")
			assert.Equal(t, *geodesy.ToDMS3(0, geodesy.FmtDM, 0), "000°00′")
			assert.Equal(t, *geodesy.ToDMS3(0, geodesy.FmtDMS, 0), "000°00′00″")
			assert.Equal(t, *geodesy.ToDMS3(0, geodesy.FmtD, 2), "000.00°")
			assert.Equal(t, *geodesy.ToDMS3(0, geodesy.FmtDM, 2), "000°00.00′")
			assert.Equal(t, *geodesy.ToDMS3(0, geodesy.FmtDMS, 2), "000°00′00.00″")
		})
	t.Run("toDMS value",
		func(t *testing.T) {
			assert.Equal(t, *geodesy.ToDMS1(45.76260), "045°45′45″")
			assert.Equal(t, *geodesy.ToDMS2(45.76260, geodesy.FmtD), "045.7626°")
			assert.Equal(t, *geodesy.ToDMS2(45.76260, geodesy.FmtDM), "045°45.76′")
			assert.Equal(t, *geodesy.ToDMS2(45.76260, geodesy.FmtDMS), "045°45′45″")
			assert.Equal(t, *geodesy.ToDMS3(45.76260, geodesy.FmtD, 0), "046°")
			assert.Equal(t, *geodesy.ToDMS3(45.76260, geodesy.FmtDM, 0), "045°46′")
			assert.Equal(t, *geodesy.ToDMS3(45.76260, geodesy.FmtDMS, 0), "045°45′45″")
			assert.Equal(t, *geodesy.ToDMS3(45.76260, geodesy.FmtD, 6), "045.762600°")
			assert.Equal(t, *geodesy.ToDMS3(45.76260, geodesy.FmtDM, 4), "045°45.7560′")
			assert.Equal(t, *geodesy.ToDMS3(45.76260, geodesy.FmtDMS, 2), "045°45′45.36″")
		})
	t.Run("toDMS round up",
		func(t *testing.T) {
			assert.Equal(t, *geodesy.ToDMS2(1.99999999999999, geodesy.FmtDM), "002°00.00′")
			assert.Equal(t, *geodesy.ToDMS2(51.19999999999999, geodesy.FmtD), "051.2000°")
			assert.Equal(t, *geodesy.ToDMS2(51.19999999999999, geodesy.FmtDM), "051°12.00′")
			assert.Equal(t, *geodesy.ToDMS2(51.19999999999999, geodesy.FmtDMS), "051°12′00″")
			assert.Equal(t, *geodesy.ToDMS2(51.99999999999999, geodesy.FmtDMS), "052°00′00″")

		})
	t.Run("toDMS NaN",
		func(t *testing.T) {
			assert.Nil(t, geodesy.ToDMS1(math.NaN()))
			assert.Nil(t, geodesy.ToDMS2(math.NaN(), geodesy.FmtD))
			assert.Nil(t, geodesy.ToDMS2(math.NaN(), geodesy.FmtDM))
			assert.Nil(t, geodesy.ToDMS2(math.NaN(), geodesy.FmtDMS))
			assert.Nil(t, geodesy.ToDMS3(math.NaN(), geodesy.FmtD, 0))
			assert.Nil(t, geodesy.ToDMS3(math.NaN(), geodesy.FmtDM, 0))
			assert.Nil(t, geodesy.ToDMS3(math.NaN(), geodesy.FmtDMS, 0))
			assert.Nil(t, geodesy.ToDMS3(math.NaN(), geodesy.FmtD, 6))
			assert.Nil(t, geodesy.ToDMS3(math.NaN(), geodesy.FmtDM, 4))
			assert.Nil(t, geodesy.ToDMS3(math.NaN(), geodesy.FmtDMS, 2))
		})
}

func TestCompass(t *testing.T) {
	assert.Equal(t, geodesy.CompassPoint1(1.0), "N")
	assert.Equal(t, geodesy.CompassPoint1(0), "N")
	assert.Equal(t, geodesy.CompassPoint1(-1), "N")
	assert.Equal(t, geodesy.CompassPoint1(359), "N")
	assert.Equal(t, geodesy.CompassPoint1(24), "NNE")
	assert.Equal(t, geodesy.CompassPoint2(24, geodesy.CardinalPrecision), "N")
	assert.Equal(t, geodesy.CompassPoint2(24, geodesy.InterCardinalPrecision), "NE")
	assert.Equal(t, geodesy.CompassPoint2(24, geodesy.SecondaryInterCardinalPrecision), "NNE")
	assert.Equal(t, geodesy.CompassPoint1(226), "SW")
	assert.Equal(t, geodesy.CompassPoint2(226, geodesy.CardinalPrecision), "W")
	assert.Equal(t, geodesy.CompassPoint2(226, geodesy.InterCardinalPrecision), "SW")
	assert.Equal(t, geodesy.CompassPoint2(226, geodesy.SecondaryInterCardinalPrecision), "SW")
	assert.Equal(t, geodesy.CompassPoint1(237), "WSW")
	assert.Equal(t, geodesy.CompassPoint2(237, geodesy.CardinalPrecision), "W")
	assert.Equal(t, geodesy.CompassPoint2(237, geodesy.InterCardinalPrecision), "SW")
	assert.Equal(t, geodesy.CompassPoint2(237, geodesy.SecondaryInterCardinalPrecision), "WSW")
}

func TestToLatLon(t *testing.T) {
	t.Run("toLat",
		func(t *testing.T) {
			assert.Equal(t, geodesy.ToLat2(51.2, geodesy.FmtDMS), "51°12′00″N")
			assert.Equal(t, geodesy.ToLat3(51.2, geodesy.FmtDMS, 0), "51°12′00″N")
		})
	t.Run("toLon",
		func(t *testing.T) {
			assert.Equal(t, geodesy.ToLon2(0.33, geodesy.FmtDMS), "000°19′48″E")
			assert.Equal(t, geodesy.ToLon3(0.33, geodesy.FmtDMS, 0), "000°19′48″E")
		})
	t.Run("toLon",
		func(t *testing.T) {
			assert.Equal(t, geodesy.ToBrng(1.0, geodesy.FmtDMS, 0), "001°00′00″")
			assert.Equal(t, geodesy.ToBrng(359.9999999999999, geodesy.FmtDMS, 0), "000°00′00″")
			assert.Equal(t, geodesy.ToBrng(math.NaN(), geodesy.FmtDMS, 2), "-")
		})
}
