package exporter

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/and-period/furumaru/api/internal/codes"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Exporter interface {
	WriteHeader(data Receipt) error
	WriteBody(data Receipt) error
	Flush() error
}

type Receipt interface {
	Header() []string
	Record() []string
}

type exporter struct {
	w *csv.Writer
}

func NewExporter(w io.Writer, encodingType codes.CharacterEncodingType) Exporter {
	if encodingType == codes.CharacterEncodingTypeShiftJIS {
		w = transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())
	}
	return &exporter{
		w: csv.NewWriter(w),
	}
}

func (e *exporter) WriteHeader(receipt Receipt) error {
	if err := e.w.Write(receipt.Header()); err != nil {
		return fmt.Errorf("exporter: failed to write csv header: %w", err)
	}
	return nil
}

func (e *exporter) WriteBody(receipt Receipt) error {
	if err := e.w.Write(receipt.Record()); err != nil {
		return fmt.Errorf("exporter: failed to write csv record: %w", err)
	}
	return nil
}

func (e *exporter) Flush() error {
	e.w.Flush()
	if err := e.w.Error(); err != nil {
		return fmt.Errorf("csv: failed to write csv: %w", err)
	}
	return nil
}
