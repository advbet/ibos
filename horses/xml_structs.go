package horses

import (
	"encoding/xml"
	"fmt"
	"time"

	"bitbucket.org/advbet/decimal"
)

type MarketStatus string
type MarketOptionStatus string

const (
	MarketOpen       MarketStatus = "OP"
	MarketNew        MarketStatus = "NE"
	MarketSuspended  MarketStatus = "SU"
	MarketClosed     MarketStatus = "CL"
	MarketProvResult MarketStatus = "PR"
	MarketFinalized  MarketStatus = "FI"
	MarketAbandoned  MarketStatus = "AB"

	MarketOptionOpen      MarketOptionStatus = "OP"
	MarketOptionDisabled  MarketOptionStatus = "DI"
	MarketOptionScratched MarketOptionStatus = "SC"
)

type Events []Event

type Event struct {
	Name           string
	GameID         int           // Unique Identifier for game
	EventDate      time.Time     // Date of the Event
	Places         int           // Indicated the number places that would pay in the market
	IsExact3       bool          // Indicates market included in Trifecta exotic bet type
	IsExact2       bool          // Indicates market included in Exacta exotic bet type
	InRunningDelay time.Duration // Describes the delay applied to in-running bet execution
	ResultExpected time.Time     // Date time at which the result is expected
	Status         MarketStatus  // Indicates the status of the event (market)
	MarketName     string        // Descriptive name of a market
	MarketID       int           // Unique Identifier for a market

	Refund       *Refund
	Participants []Participant
}

type Participant struct {
	Name           string             // Descriptive name of a participant
	OpenWinPrice   decimal.Number     // Price for win bet when the market was opened
	Status         MarketOptionStatus // Indicates the status of a marked
	DisplayNumber  int                // Runner number (saddle number in most instances – in some American races the runner number and saddle number may differ)
	ResultPosition int
}

type Refund struct {
	IsBefore        bool      // Whether the refund is affective before or after the applied date
	AppliedDateTime time.Time // Date refund was applied
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *Events) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	data := struct {
		Games struct {
			Games []struct {
				GameID int    `xml:"GameID,attr"` // Unique Identifier for game
				Name   string `xml:"Name"`        // Descriptive name of a game
			} `xml:"Game"`
		} `xml:"Games"`
		Participants struct {
			Participants []struct {
				ID     int    `xml:"ParticipantID,attr"` // Unique Identifier for a participant
				GameID int    `xml:"GameID,attr"`        // Maps back to the Game
				Name   string `xml:"Name"`               // Descriptive name of a participant
			} `xml:"Participant"`
		} `xml:"Participants"`
		Markets struct {
			Markets []struct {
				ID             int          `xml:"MarketID,attr"`       // Unique Identifier for a market
				GameID         int          `xml:"GameID,attr"`         // The game this market belongs to
				EventDate      string       `xml:"EventDate,attr"`      // Date of the Event
				Places         int          `xml:"NoOfPlaces,attr"`     // Indicated the number places that would pay in the market
				IsExact3       int          `xml:"IsExact3,attr"`       // Indicates market included in Trifecta exotic bet type
				IsExact2       int          `xml:"IsExact2,attr"`       // Indicates market included in Exacta exotic bet type
				InRunningDelay int          `xml:"InRunningDelay,attr"` // Describes the delay applied to in-running bet execution in seconds
				ResultExpected string       `xml:"ResultExpected,attr"` // Date time at which the result is expected
				Status         MarketStatus `xml:"EventStatus,attr"`    // Indicates the status of the event (market)
				Name           string       `xml:"Name"`                // Descriptive name of a market
			} `xml:"Market"`
		} `xml:"Markets"`
		MarketOptions struct {
			Options []struct {
				ID             int                `xml:"MarketOptionID,attr"` // Unique Identifier for a market option
				ParticipantID  int                `xml:"ParticipantID,attr"`  // Refers to a participant forming this market option from
				ResultPosition int                `xml:"ResultPosition,attr"`
				OpenWinPrice   decimal.Number     `xml:"OpenWinPrice,attr"`           // Price for win bet when the market was opened
				MarketID       int                `xml:"MarketID,attr"`               // The event this market belongs to
				Status         MarketOptionStatus `xml:"MarketOptionStatusCode,attr"` // Indicates the status of a marked
				DisplayNumber  int                `xml:"DisplayNumber,attr"`          // Runner number (saddle number in most instances – in some American races the runner number and saddle number may differ)
			} `xml:"MarketOption"`
		} `xml:"MarketOptions"`
		MarketRefunds struct {
			Refunds []struct {
				MarketID        int    `xml:"MarketID,attr"`        // ID of refunded market
				IsBefore        int    `xml:"IsBefore,attr"`        // Whether the refund is affective before or after the applied date
				AppliedDateTime string `xml:"AppliedDateTime,attr"` // Date refund was applied
			} `xml:"MarketRefund"`
		} `xml:"MarketRefunds"`
	}{}
	if err := d.DecodeElement(&data, &start); err != nil {
		return err
	}
	games := make(map[int]string)
	for _, game := range data.Games.Games {
		games[game.GameID] = game.Name
	}

	participantsMap := make(map[int]string)
	for _, p := range data.Participants.Participants {
		participantsMap[p.ID] = p.Name
	}

	*e = make([]Event, 0)
	for _, market := range data.Markets.Markets {
		var refund *Refund
		for _, r := range data.MarketRefunds.Refunds {
			if r.MarketID == market.ID {
				appliedDateTime, err := time.Parse("2006-01-02T15:04:05", r.AppliedDateTime)
				if err != nil {
					return err
				}
				refund = &Refund{
					IsBefore:        r.IsBefore == 1,
					AppliedDateTime: appliedDateTime,
				}
				break
			}
		}

		participants := make([]Participant, 0)
		for _, op := range data.MarketOptions.Options {
			if market.ID != op.MarketID {
				continue
			}

			participant, ok := participantsMap[op.ParticipantID]
			if !ok {
				return fmt.Errorf("participant ID %d not found for market option", op.ParticipantID)
			}
			participants = append(participants, Participant{
				Name:           participant,
				OpenWinPrice:   op.OpenWinPrice,
				Status:         op.Status,
				DisplayNumber:  op.DisplayNumber,
				ResultPosition: op.ResultPosition,
			})
		}

		eventDate, err := time.Parse("2006-01-02T15:04:05", market.EventDate)
		if err != nil {
			return err
		}
		resultExpected, err := time.Parse("2006-01-02T15:04:05", market.ResultExpected)
		if err != nil {
			return err
		}

		*e = append(*e, Event{
			Name:           games[market.GameID],
			GameID:         market.GameID,
			EventDate:      eventDate,
			Places:         market.Places,
			IsExact3:       market.IsExact3 == 1,
			IsExact2:       market.IsExact2 == 1,
			InRunningDelay: time.Duration(market.InRunningDelay) * time.Second,
			ResultExpected: resultExpected,
			Status:         market.Status,
			MarketName:     market.Name,
			MarketID:       market.ID,
			Refund:         refund,
			Participants:   participants,
		})
	}

	return nil
}
