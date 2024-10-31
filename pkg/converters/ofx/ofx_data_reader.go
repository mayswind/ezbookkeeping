package ofx

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

var ofx2HeaderPattern = regexp.MustCompile("<\\?OFX( +[A-Z]+=\"[^=]*\")* *\\?>")
var ofx2HeaderAttributePattern = regexp.MustCompile(" +([A-Z]+)=\"([^=]*)\"")

// ofxFileReader defines the structure of open financial exchange (ofx) file reader
type ofxFileReader struct {
	fileHeader *ofxFileHeader
	xmlDecoder *xml.Decoder
}

// read returns the imported open financial exchange (ofx) file
func (r *ofxFileReader) read(ctx core.Context) (*ofxFile, error) {
	file := &ofxFile{}

	err := r.xmlDecoder.Decode(&file)

	if err != nil {
		return nil, err
	}

	file.FileHeader = r.fileHeader

	return file, nil
}

func createNewOFXFileReader(ctx core.Context, data []byte) (*ofxFileReader, error) {
	if len(data) > 5 && string(data[0:5]) == "<?xml" { // ofx 2.x starts with <?xml
		return createNewOFX2FileReader(ctx, data, true)
	} else if len(data) > 10 && string(data[0:10]) == "OFXHEADER:" { // ofx 1.x starts with OFXHEADER:

	} else if len(data) > 5 && string(data[0:5]) == "<OFX>" { // no ofx header
		return createNewOFX2FileReader(ctx, data, false)
	}

	return nil, errs.ErrInvalidOFXFile
}

func createNewOFX2FileReader(ctx core.Context, data []byte, withHeader bool) (*ofxFileReader, error) {
	var fileHeader *ofxFileHeader = nil
	var err error

	if withHeader {
		fileHeader, err = readOFX2FileHeader(ctx, data)

		if err != nil {
			return nil, err
		}

		if fileHeader.OFXDeclarationVersion != ofxVersion2 {
			log.Errorf(ctx, "[ofx_data_reader.createNewOFX2FileReader] cannot parse ofx 2.x file header, because declaration version is \"%s\"", fileHeader.OFXDeclarationVersion)
			return nil, errs.ErrInvalidOFXFile
		}
	}

	xmlDecoder := xml.NewDecoder(bytes.NewReader(data))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	return &ofxFileReader{
		fileHeader: fileHeader,
		xmlDecoder: xmlDecoder,
	}, nil
}

func readOFX2FileHeader(ctx core.Context, data []byte) (fileHeader *ofxFileHeader, err error) {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	fileHeader = &ofxFileHeader{}
	headerLine := ""

	for scanner.Scan() {
		line := scanner.Text()

		ofxHeaderStartIndex := strings.Index(line, "<?OFX ")

		if ofxHeaderStartIndex >= 0 {
			headerLine = ofx2HeaderPattern.FindString(line)
			break
		}
	}

	if headerLine == "" {
		log.Errorf(ctx, "[ofx_data_reader.readOFX2FileHeader] cannot find ofx 2.x file header")
		return nil, errs.ErrInvalidOFXFile
	}

	headerAttributes := ofx2HeaderAttributePattern.FindAllStringSubmatch(headerLine, -1)

	for _, attributeItems := range headerAttributes {
		if len(attributeItems) != 3 {
			log.Warnf(ctx, "[ofx_data_reader.readOFX2FileHeader] cannot parse line in ofx 2.x file header, because item is \"%s\"", attributeItems)
			continue
		}

		name := attributeItems[1]
		value := attributeItems[2]

		if name == "OFXHEADER" {
			fileHeader.OFXDeclarationVersion = oFXDeclarationVersion(value)
		} else if name == "VERSION" {
			fileHeader.OFXDataVersion = value
		} else if name == "SECURITY" {
			fileHeader.Security = value
		} else if name == "OLDFILEUID" {
			fileHeader.OldFileUid = value
		} else if name == "NEWFILEUID" {
			fileHeader.NewFileUid = value
		} else {
			log.Warnf(ctx, "[ofx_data_reader.readOFX2FileHeader] cannot parse unknown header line in ofx 2.x file header, because item is \"%s\"", attributeItems)
			continue
		}
	}

	return fileHeader, nil
}
