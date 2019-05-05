package qualomat

import "time"

type date struct{ time.Time }

func (d *date) UnmarshalJSON(data []byte) error {
	var err error
	d.Time, err = time.Parse("\"2006-01-02\"", string(data))
	return err
}
