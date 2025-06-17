package camt

import (
	"bytes"
	"encoding/xml"

	"golang.org/x/net/html/charset"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// camt053FileReader defines the structure of camt.053 file reader
type camt053FileReader struct {
	xmlDecoder *xml.Decoder
}

// read returns the imported camt.053 data
func (r *camt053FileReader) read(ctx core.Context) (*camt053File, error) {
	file := &camt053File{}

	err := r.xmlDecoder.Decode(&file)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func createNewCamt053FileReader(data []byte) (*camt053FileReader, error) {
	if len(data) > 5 && data[0] == 0x3C && data[1] == 0x3F && data[2] == 0x78 && data[3] == 0x6D && data[4] == 0x6C { // <?xml
		xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
		xmlDecoder.CharsetReader = charset.NewReaderLabel

		return &camt053FileReader{
			xmlDecoder: xmlDecoder,
		}, nil
	}

	return nil, errs.ErrInvalidXmlFile
}
