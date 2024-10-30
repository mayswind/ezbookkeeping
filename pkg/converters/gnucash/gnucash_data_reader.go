package gnucash

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"

	"golang.org/x/net/html/charset"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// gnucashDatabaseReader defines the structure of gnucash database reader
type gnucashDatabaseReader struct {
	xmlDecoder *xml.Decoder
}

// read returns the imported gnucash data
func (r *gnucashDatabaseReader) read(ctx core.Context) (*gnucashDatabase, error) {
	database := &gnucashDatabase{}

	err := r.xmlDecoder.Decode(&database)

	if err != nil {
		return nil, err
	}

	return database, nil
}

func createNewGnuCashDatabaseReader(data []byte) (*gnucashDatabaseReader, error) {
	if len(data) > 2 && data[0] == 0x1F && data[1] == 0x8B { // gzip magic number
		gzipReader, err := gzip.NewReader(bytes.NewReader(data))

		if err != nil {
			return nil, err
		}

		xmlDecoder := xml.NewDecoder(gzipReader)
		xmlDecoder.CharsetReader = charset.NewReaderLabel

		return &gnucashDatabaseReader{
			xmlDecoder: xmlDecoder,
		}, nil
	} else if len(data) > 5 && data[0] == 0x3C && data[1] == 0x3F && data[2] == 0x78 && data[3] == 0x6D && data[4] == 0x6C { // <?xml
		xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
		xmlDecoder.CharsetReader = charset.NewReaderLabel

		return &gnucashDatabaseReader{
			xmlDecoder: xmlDecoder,
		}, nil
	}

	return nil, errs.ErrInvalidGnuCashFile
}
