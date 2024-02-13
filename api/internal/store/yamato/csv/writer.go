package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Writer struct {
	w *csv.Writer
}

func NewWriter(w io.Writer, encodingType CharacterEncodingType) *Writer {
	return &Writer{
		w: newCSVWriter(w, encodingType),
	}
}

func newCSVWriter(w io.Writer, encodingType CharacterEncodingType) *csv.Writer {
	switch encodingType {
	case CharacterEncodingTypeUTF8:
		return csv.NewWriter(w)
	case CharacterEncodingTypeShiftJIS:
		ww := transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
		return csv.NewWriter(ww)
	default:
		return csv.NewWriter(w)
	}
}

func (w *Writer) WriteHeader() error {
	if err := w.w.Write(deliveryReceiptHeaders); err != nil {
		return fmt.Errorf("csv: failed to write csv header: %w", err)
	}
	return nil
}

func (w *Writer) WriteBody(data *DeliveryReceipt) error {
	if err := w.w.Write(data.Record()); err != nil {
		return fmt.Errorf("csv: failed to write csv record: %w", err)
	}
	return nil
}

func (w *Writer) Flush() error {
	w.w.Flush()
	if err := w.w.Error(); err != nil {
		return fmt.Errorf("csv: failed to write csv: %w", err)
	}
	return nil
}
