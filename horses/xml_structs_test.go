package horses

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"

	"github.com/advbet/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDecimal(t *testing.T, s string) decimal.Number {
	t.Helper()
	d, err := decimal.FromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func makeTime(t *testing.T, s string) time.Time {
	tm, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

func TestParseEvents(t *testing.T) {
	tests := []struct {
		file   string
		events Events
	}{
		{
			file: "testdata/igaming_20200828183127.xml",
			events: Events{
				{
					Name:           "Newmarket",
					EventDate:      makeTime(t, "2020-08-28T18:20:00"),
					Places:         4,
					IsExact3:       false,
					IsExact2:       false,
					ResultExpected: makeTime(t, "2020-08-28T18:30:00"),
					Status:         MarketProvResult,
					MarketName:     "Race 5",
					Refund:         nil,
					MarketID:       6060380,
					GameID:         519452,
					Participants: []Participant{
						{
							Name:           "Han Solo Berger",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  1,
							ResultPosition: 3,
						},
						{
							Name:           "Zig Zag Zyggy",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  3,
							ResultPosition: 0,
						},
						{
							Name:           "Live In The Moment",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  4,
							ResultPosition: 1,
						},
						{
							Name:           "Thegreatestshowman",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  5,
							ResultPosition: 4,
						},
						{
							Name:           "Rio Ronaldo",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  7,
							ResultPosition: 0,
						},
						{
							Name:           "Grandfather Tom",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  8,
							ResultPosition: 0,
						},
						{
							Name:           "Colonel Frank",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  9,
							ResultPosition: 2,
						},
						{
							Name:           "Jumira Bridge",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  10,
							ResultPosition: 0,
						},
						{
							Name:           "Haveoneyerself",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  11,
							ResultPosition: 0,
						},
						{
							Name:           "Green Door",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  12,
							ResultPosition: 0,
						},
						{
							Name:           "Shamshon",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  14,
							ResultPosition: 0,
						},
					},
				},
				{
					Name:           "Newmarket",
					EventDate:      makeTime(t, "2020-08-28T19:30:00"),
					Places:         3,
					IsExact3:       false,
					IsExact2:       false,
					ResultExpected: makeTime(t, "2020-08-28T19:40:00"),
					Status:         MarketOpen,
					MarketName:     "Race 7",
					Refund:         nil,
					MarketID:       6060385,
					GameID:         519452,
					Participants: []Participant{
						{
							Name:           "Alpine Mistral",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  2,
							ResultPosition: 0,
						},
						{
							Name:           "Katniss Everdeen",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  6,
							ResultPosition: 0,
						},
						{
							Name:           "Licit",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  8,
							ResultPosition: 0,
						},
						{
							Name:           "Corofin",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  9,
							ResultPosition: 0,
						},
					},
				},
				{
					Name:           "Newmarket",
					EventDate:      makeTime(t, "2020-08-28T20:00:00"),
					Places:         3,
					IsExact3:       false,
					IsExact2:       false,
					ResultExpected: makeTime(t, "2020-08-28T20:10:00"),
					Status:         MarketOpen,
					MarketName:     "Race 8",
					Refund:         nil,
					MarketID:       6060386,
					GameID:         519452,
					Participants: []Participant{
						{
							Name:           "Filles De Fleur",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  1,
							ResultPosition: 0,
						},
						{
							Name:           "Zephyrina",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  3,
							ResultPosition: 0,
						},
						{
							Name:           "Mrs Meader",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  5,
							ResultPosition: 0,
						},
					},
				},
				{
					Name:           "Newmarket",
					EventDate:      makeTime(t, "2020-08-28T20:35:00"),
					Places:         3,
					IsExact3:       false,
					IsExact2:       false,
					ResultExpected: makeTime(t, "2020-08-28T20:45:00"),
					Status:         MarketOpen,
					MarketName:     "Race 9",
					Refund:         nil,
					MarketID:       6060387,
					GameID:         519452,
					Participants: []Participant{
						{
							Name:           "Fantastic Blue",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  1,
							ResultPosition: 0,
						},
						{
							Name:           "Berrahri",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  5,
							ResultPosition: 0,
						},
						{
							Name:           "Overpriced Mixer",
							OpenWinPrice:   newDecimal(t, "0.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  10,
							ResultPosition: 0,
						},
					},
				},
				{
					Name:           "Curragh",
					EventDate:      makeTime(t, "2020-08-28T18:10:00"),
					Places:         2,
					IsExact3:       false,
					IsExact2:       false,
					ResultExpected: makeTime(t, "2020-08-28T18:20:00"),
					Status:         MarketFinalized,
					MarketName:     "Race 7",
					Refund:         nil,
					MarketID:       6064409,
					GameID:         524074,
					Participants: []Participant{
						{
							Name:           "Cursory Exam",
							OpenWinPrice:   newDecimal(t, "4.5000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  1,
							ResultPosition: 0,
						},
						{
							Name:           "Get Noticed",
							OpenWinPrice:   newDecimal(t, "12.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  2,
							ResultPosition: 0,
						},
						{
							Name:           "Shambara",
							OpenWinPrice:   newDecimal(t, "3.5000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  3,
							ResultPosition: 0,
						},
						{
							Name:           "Star Image",
							OpenWinPrice:   newDecimal(t, "9.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  4,
							ResultPosition: 1,
						},
						{
							Name:           "Abogados",
							OpenWinPrice:   newDecimal(t, "4.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  5,
							ResultPosition: 0,
						},
						{
							Name:           "Omakase",
							OpenWinPrice:   newDecimal(t, "1.7500"),
							Status:         MarketOptionOpen,
							DisplayNumber:  6,
							ResultPosition: 2,
						},
						{
							Name:           "Verbal Fencing",
							OpenWinPrice:   newDecimal(t, "8.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  7,
							ResultPosition: 0,
						},
					},
				},
				{
					Name:           "Curragh",
					EventDate:      makeTime(t, "2020-08-28T18:45:00"),
					Places:         4,
					IsExact3:       false,
					IsExact2:       false,
					ResultExpected: makeTime(t, "2020-08-28T18:55:00"),
					Status:         MarketOpen,
					MarketName:     "Race 8",
					Refund:         nil,
					MarketID:       6064412,
					GameID:         524074,
					Participants: []Participant{
						{
							Name:           "Ramiro",
							OpenWinPrice:   newDecimal(t, "6.5000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  8,
							ResultPosition: 0,
						},
						{
							Name:           "Scherzando",
							OpenWinPrice:   newDecimal(t, "4.5000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  10,
							ResultPosition: 0,
						},
						{
							Name:           "Puddle Of Luck",
							OpenWinPrice:   newDecimal(t, "12.0000"),
							Status:         MarketOptionOpen,
							DisplayNumber:  12,
							ResultPosition: 0,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		blob, err := ioutil.ReadFile(test.file)
		require.NoError(t, err)
		var obj Events
		err = xml.Unmarshal(blob, &obj)
		require.NoError(t, err)
		assert.Equal(t, test.events, obj, test.file)
	}
}
