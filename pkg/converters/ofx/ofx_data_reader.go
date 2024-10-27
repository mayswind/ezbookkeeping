package ofx

import (
	"bytes"
	"encoding/xml"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ofxFileReader defines the structure of open financial exchange (ofx) file reader
type ofxFileReader struct {
	xmlDecoder *xml.Decoder
}

// read returns the imported open financial exchange (ofx) file
func (r *ofxFileReader) read(ctx core.Context) (*ofxFile, error) {
	file := &ofxFile{}

	err := r.xmlDecoder.Decode(&file)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func createNewOFXFileReader(data []byte) (*ofxFileReader, error) {
	if len(data) > 5 && data[0] == 0x3C && data[1] == 0x3F && data[2] == 0x78 && data[3] == 0x6D && data[4] == 0x6C { // ofx 2.x starts with <?xml
		xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
		xmlDecoder.CharsetReader = utils.IdentReader

		return &ofxFileReader{
			xmlDecoder: xmlDecoder,
		}, nil
	} else if len(data) > 13 && string(data[0:13]) == "OFXHEADER:100" { // ofx 1.x starts with OFXHEADER:100

	} else if len(data) > 5 && string(data[0:5]) == "<OFX>" { // no ofx header
		xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
		xmlDecoder.CharsetReader = utils.IdentReader

		return &ofxFileReader{
			xmlDecoder: xmlDecoder,
		}, nil
	}

	return nil, errs.ErrInvalidOFXFile
}
