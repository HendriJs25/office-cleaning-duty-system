package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	*d = Date(t)

	return nil
}

func (d *Date) TimePtr() *time.Time {
	if d == nil {
		return nil
	}

	t := time.Time(*d)
	return &t
}
