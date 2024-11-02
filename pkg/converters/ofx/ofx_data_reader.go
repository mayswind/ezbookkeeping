package ofx

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/sgml"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const ofxUnicodeEncoding = "unicode"
const ofxUSAsciiEncoding = "usascii"
const ofx1SGMLDataFormat = "OFXSGML"

var ofx2HeaderPattern = regexp.MustCompile("<\\?OFX( +[A-Z]+=\"[^=]*\")* *\\?>")
var ofx2HeaderAttributePattern = regexp.MustCompile(" +([A-Z]+)=\"([^=]*)\"")

// ofxFileReader defines the structure of open financial exchange (ofx) file reader
type ofxFileReader interface {
	// read returns the imported open financial exchange (ofx) file
	read(ctx core.Context) (*ofxFile, error)
}

// ofxVersion1FileReader defines the structure of open financial exchange (ofx) declaration version 1.x file reader
type ofxVersion1FileReader struct {
	fileHeader  *ofxFileHeader
	sgmlDecoder *sgml.Decoder
}

// ofxVersion2FileReader defines the structure of open financial exchange (ofx) declaration version 2.x file reader
type ofxVersion2FileReader struct {
	fileHeader *ofxFileHeader
	xmlDecoder *xml.Decoder
}

// read returns the imported open financial exchange (ofx) file
func (r *ofxVersion1FileReader) read(ctx core.Context) (*ofxFile, error) {
	file := &ofxFile{}

	err := r.sgmlDecoder.Decode(&file)

	if err != nil {
		log.Errorf(ctx, "[ofxVersion1FileReader.read] cannot read ofx 1.x file, because %s", err.Error())
		return nil, errs.ErrInvalidOFXFile
	}

	file.FileHeader = r.fileHeader

	return file, nil
}

// read returns the imported open financial exchange (ofx) file
func (r *ofxVersion2FileReader) read(ctx core.Context) (*ofxFile, error) {
	file := &ofxFile{}

	err := r.xmlDecoder.Decode(&file)

	if err != nil {
		log.Errorf(ctx, "[ofxVersion2FileReader.read] cannot read ofx 2.x file, because %s", err.Error())
		return nil, errs.ErrInvalidOFXFile
	}

	file.FileHeader = r.fileHeader

	return file, nil
}

func createNewOFXFileReader(ctx core.Context, data []byte) (ofxFileReader, error) {
	firstNonCrLfIndex := 0

	for i := 0; i < len(data); i++ {
		if data[i] != '\n' && data[i] != '\r' {
			firstNonCrLfIndex = i
			break
		}
	}

	if len(data) > 5 && string(data[firstNonCrLfIndex:firstNonCrLfIndex+5]) == "<?xml" { // ofx 2.x starts with <?xml
		return createNewOFX2FileReader(ctx, data, true)
	} else if len(data) > 10 && string(data[firstNonCrLfIndex:firstNonCrLfIndex+10]) == "OFXHEADER:" { // ofx 1.x starts with OFXHEADER:
		return createNewOFX1FileReader(ctx, data)
	} else if len(data) > 5 && string(data[firstNonCrLfIndex:firstNonCrLfIndex+5]) == "<OFX>" { // no ofx header
		return createNewOFX2FileReader(ctx, data, false)
	}

	return nil, errs.ErrInvalidOFXFile
}

func createNewOFX1FileReader(ctx core.Context, data []byte) (ofxFileReader, error) {
	fileHeader, fileData, dataType, enc, err := readOFX1FileHeader(ctx, data)

	if err != nil {
		return nil, err
	}

	if fileHeader.OFXDeclarationVersion != ofxVersion1 {
		log.Errorf(ctx, "[ofx_data_reader.createNewOFX1FileReader] cannot parse ofx 1.x file header, because declaration version is \"%s\"", fileHeader.OFXDeclarationVersion)
		return nil, errs.ErrInvalidOFXFile
	}

	if dataType != ofx1SGMLDataFormat {
		log.Errorf(ctx, "[ofx_data_reader.createNewOFX1FileReader] cannot parse ofx 1.x file header, because data type is \"%s\"", dataType)
		return nil, errs.ErrInvalidOFXFile
	}

	reader := bytes.NewReader(fileData)
	buffer := &bytes.Buffer{}

	if enc != nil {
		transformReader := transform.NewReader(reader, enc.NewDecoder())
		_, err = buffer.ReadFrom(transformReader)
	} else {
		_, err = buffer.ReadFrom(reader)
	}

	if err != nil {
		log.Errorf(ctx, "[ofx_data_reader.createNewOFX1FileReader] cannot read ofx 1.x file content, because %s", err.Error())
		return nil, errs.ErrInvalidOFXFile
	}

	sgmlData := buffer.String()
	stringReader := strings.NewReader(sgmlData)
	sgmlDecoder := sgml.NewDecoder(stringReader)

	return &ofxVersion1FileReader{
		fileHeader:  fileHeader,
		sgmlDecoder: sgmlDecoder,
	}, nil
}

func createNewOFX2FileReader(ctx core.Context, data []byte, withHeader bool) (ofxFileReader, error) {
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

	return &ofxVersion2FileReader{
		fileHeader: fileHeader,
		xmlDecoder: xmlDecoder,
	}, nil
}

func readOFX1FileHeader(ctx core.Context, data []byte) (fileHeader *ofxFileHeader, fileData []byte, dataType string, enc encoding.Encoding, err error) {
	fileHeader = &ofxFileHeader{}
	dataType = ""
	fileEncoding := ""
	fileCharset := ""
	fileDataStartPosition := 0
	lastCrLf := -1

	for i := 0; i < len(data); i++ {
		if data[i] != '\n' && data[i] != '\r' {
			continue
		}

		if lastCrLf == i-1 {
			lastCrLf = i
			continue
		}

		line := string(data[lastCrLf+1 : i])

		if strings.Index(line, "<OFX>") == 0 {
			fileDataStartPosition = lastCrLf + 1
			break
		}

		lastCrLf = i

		if line == "" {
			continue
		}

		items := strings.Split(line, ":")

		if len(items) != 2 {
			log.Warnf(ctx, "[ofx_data_reader.readOFX1FileHeader] cannot parse line in ofx 1.x file header, because line is \"%s\"", line)
			continue
		}

		key := items[0]
		value := items[1]

		if key == "OFXHEADER" {
			fileHeader.OFXDeclarationVersion = oFXDeclarationVersion(value)
		} else if key == "DATA" {
			dataType = value
		} else if key == "VERSION" {
			fileHeader.OFXDataVersion = value
		} else if key == "SECURITY" {
			fileHeader.Security = value
		} else if key == "ENCODING" {
			fileEncoding = strings.ToLower(value)
		} else if key == "CHARSET" {
			fileCharset = strings.ToLower(value)
		} else if key == "COMPRESSION" {
			continue // ignore
		} else if key == "OLDFILEUID" {
			fileHeader.OldFileUid = value
		} else if key == "NEWFILEUID" {
			fileHeader.NewFileUid = value
		} else {
			log.Warnf(ctx, "[ofx_data_reader.readOFX1FileHeader] cannot parse unknown header line in ofx 1.x file header, because line is \"%s\"", line)
			continue
		}
	}

	if fileEncoding == ofxUSAsciiEncoding {
		if utils.IsStringOnlyContainsDigits(fileCharset) {
			fileCharset = "cp" + fileCharset
		}

		enc, _ = charset.Lookup(fileCharset)

		if enc == nil {
			enc, _ = charset.Lookup("us-ascii")
		}

		if enc == nil {
			enc = charmap.Windows1252
		}
	} else if fileEncoding == ofxUnicodeEncoding {
		enc, _ = charset.Lookup(ofxUnicodeEncoding)

		if enc == nil {
			enc = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		}
	} else {
		log.Errorf(ctx, "[ofx_data_reader.readOFX1FileHeader] cannot parse ofx 1.x file, because encoding \"%s\" is unknown", fileEncoding)
		return nil, nil, "", nil, errs.ErrInvalidOFXFile
	}

	return fileHeader, data[fileDataStartPosition:], dataType, enc, nil
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
