package sgml

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const sgmlTagName = "sgml"
const sgmlNameFieldName = "SGMLName"
const xmlTagName = "xml"           // reuse xml tag
const xmlNameFieldName = "XMLName" // reuse xml tag

// sgmlFieldType represents SGML field type
type sgmlFieldType byte

// Transaction template types
const (
	sgmlNotSupportedField sgmlFieldType = 0
	sgmlTextualField      sgmlFieldType = 1
	sgmlStructField       sgmlFieldType = 2
	sgmlStructSliceField  sgmlFieldType = 3
)

// sgmlTypeInfo represents the struct of SGML type reflection info
type sgmlTypeInfo struct {
	supportedFields map[string]*sgmlFieldInfo
}

// sgmlFieldInfo represents the struct of SGML field info
type sgmlFieldInfo struct {
	sgmlFieldName   string
	sgmlFieldType   sgmlFieldType
	structFieldName string
}

type Decoder struct {
	xmlDecoder *xml.Decoder
}

var sgmlTypeInfoMap sync.Map // map[reflect.Type]*typeInfo

// Decode unmarshal the specified struct instance and returns whether error occurs
func (d *Decoder) Decode(v any) error {
	value := reflect.ValueOf(v).Elem()
	finalValue := value
	finalType := value.Type()

	for finalValue.Kind() == reflect.Pointer {
		finalValue = value.Elem()
		finalType = finalValue.Type()
	}

	rootNameField, exists := finalType.FieldByName(sgmlNameFieldName)

	if !exists {
		rootNameField, exists = finalType.FieldByName(xmlNameFieldName)
	}

	if !exists {
		return nil
	}

	rootElementName := rootNameField.Tag.Get(sgmlTagName)

	if rootElementName == "" {
		rootElementName = rootNameField.Tag.Get(xmlTagName)
	}

	for {
		token, err := d.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == rootElementName {
				return d.unmarshal(value.Elem(), rootElementName)
			}
		}
	}

	return nil
}

func (d *Decoder) unmarshal(element reflect.Value, elementName string) error {
	typeInfo, err := d.getStructTypeInfo(element.Type())

	if err != nil {
		return err
	}

	if typeInfo == nil {
		return errs.ErrInvalidSGMLFile
	}

	textualFieldWithoutEndElementNames := make(map[string]bool)
	textualFieldValues := make(map[string]string)

	hasEndElement := false
	currentSGMLFieldName := ""

	for {
		token, err := d.xmlDecoder.RawToken()

		if err == io.EOF {
			break
		}

		switch token := token.(type) {
		case xml.StartElement:
			if fieldInfo, exists := typeInfo.supportedFields[token.Name.Local]; exists {
				if fieldInfo.sgmlFieldType == sgmlStructField || fieldInfo.sgmlFieldType == sgmlStructSliceField {
					field := element.FieldByName(fieldInfo.structFieldName)
					childElementType := field.Type()
					childElementKind := field.Kind()
					var childElement reflect.Value

					if fieldInfo.sgmlFieldType == sgmlStructSliceField {
						childElementType = childElementType.Elem()
						childElementKind = childElementType.Kind()
					}

					if childElementKind == reflect.Pointer {
						childElement = reflect.New(childElementType.Elem())
					} else if childElementKind == reflect.Struct {
						childElement = reflect.New(childElementType)
					}

					err := d.unmarshal(childElement.Elem(), fieldInfo.sgmlFieldName)

					if err != nil {
						return err
					}

					if childElementKind == reflect.Struct {
						childElement = childElement.Elem()
					}

					if fieldInfo.sgmlFieldType == sgmlStructField {
						field.Set(childElement)
					} else if fieldInfo.sgmlFieldType == sgmlStructSliceField {
						if field.Len() == 0 {
							slice := reflect.MakeSlice(reflect.SliceOf(childElement.Type()), 0, 0)
							field.Set(reflect.Append(slice, childElement))
						} else {
							field.Set(reflect.Append(field.Addr().Elem(), childElement))
						}
					}
				} else if fieldInfo.sgmlFieldType == sgmlTextualField {
					currentSGMLFieldName = token.Name.Local
					textualFieldWithoutEndElementNames[token.Name.Local] = true
				}
			}
		case xml.EndElement:
			if fieldInfo, exists := typeInfo.supportedFields[token.Name.Local]; exists {
				if fieldInfo.sgmlFieldType == sgmlTextualField {
					delete(textualFieldWithoutEndElementNames, token.Name.Local)
				}
			} else if token.Name.Local == elementName {
				hasEndElement = true
				break
			}
		case xml.CharData:
			if currentSGMLFieldName != "" {
				if fieldInfo, exists := typeInfo.supportedFields[currentSGMLFieldName]; exists {
					if fieldInfo.sgmlFieldType == sgmlTextualField {
						textualFieldValues[currentSGMLFieldName] = string(token)
					}
				}
			}

			currentSGMLFieldName = ""
		}

		if hasEndElement {
			break
		}
	}

	if !hasEndElement {
		return errs.ErrInvalidSGMLFile
	}

	for sgmlFieldName, fieldValue := range textualFieldValues {
		finalValue := d.getActualFieldValue(sgmlFieldName, fieldValue, textualFieldWithoutEndElementNames)
		fieldInfo, exists := typeInfo.supportedFields[sgmlFieldName]

		if !exists {
			continue
		}

		field := element.FieldByName(fieldInfo.structFieldName)
		field.SetString(finalValue)
	}

	return nil
}

func (d *Decoder) getStructTypeInfo(reflectType reflect.Type) (*sgmlTypeInfo, error) {
	if reflectType.Kind() != reflect.Struct {
		return nil, nil
	}

	typeInfo, exists := sgmlTypeInfoMap.Load(reflectType)

	if exists {
		return typeInfo.(*sgmlTypeInfo), nil
	}

	newTypeInfo := &sgmlTypeInfo{
		supportedFields: make(map[string]*sgmlFieldInfo),
	}

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)

		if field.Anonymous {
			fieldType := field.Type

			if fieldType.Kind() == reflect.Struct {
				fieldSgmlTypeInfo, err := d.getStructTypeInfo(fieldType)

				if err != nil {
					return nil, err
				}

				for sgmlFieldName, fieldInfo := range fieldSgmlTypeInfo.supportedFields {
					newTypeInfo.supportedFields[sgmlFieldName] = fieldInfo
				}
			}

			continue
		} else if !field.IsExported() {
			continue
		}

		sgmlFieldName := field.Tag.Get(sgmlTagName)

		if sgmlFieldName == "" {
			sgmlFieldName = field.Tag.Get(xmlTagName)
		}

		if sgmlFieldName == "" || field.Name == sgmlNameFieldName || field.Name == xmlNameFieldName {
			continue
		}

		sgmlFieldType := sgmlNotSupportedField
		finalFieldType := field.Type

		for finalFieldType.Kind() == reflect.Pointer {
			finalFieldType = finalFieldType.Elem()
		}

		switch finalFieldType.Kind() {
		case reflect.String:
			sgmlFieldType = sgmlTextualField
		case reflect.Struct:
			sgmlFieldType = sgmlStructField
		case reflect.Slice:
			childFinalFieldType := finalFieldType.Elem()

			for childFinalFieldType.Kind() == reflect.Pointer {
				childFinalFieldType = childFinalFieldType.Elem()
			}

			if childFinalFieldType.Kind() == reflect.Struct {
				sgmlFieldType = sgmlStructSliceField
			}
		default:
			sgmlFieldType = sgmlNotSupportedField
		}

		if sgmlFieldType == sgmlNotSupportedField {
			return nil, errs.ErrInvalidSGMLFile
		}

		newTypeInfo.supportedFields[sgmlFieldName] = &sgmlFieldInfo{
			sgmlFieldName:   sgmlFieldName,
			sgmlFieldType:   sgmlFieldType,
			structFieldName: field.Name,
		}
	}

	typeInfo, _ = sgmlTypeInfoMap.LoadOrStore(reflectType, newTypeInfo)

	return typeInfo.(*sgmlTypeInfo), nil
}

func (d *Decoder) getActualFieldValue(fieldName string, fieldValue string, textualFieldWithoutEndElementNames map[string]bool) string {
	_, notHasEndElement := textualFieldWithoutEndElementNames[fieldName]

	if !notHasEndElement {
		return fieldValue
	}

	for i := 0; i < len(fieldValue); i++ {
		if fieldValue[i] == '\r' || fieldValue[i] == '\n' {
			return fieldValue[0:i]
		}
	}

	return fieldValue
}

// NewDecoder creates a new SGML parser reading from specified io reader
func NewDecoder(reader io.Reader) *Decoder {
	xmlDecoder := xml.NewDecoder(reader)
	xmlDecoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	return &Decoder{
		xmlDecoder: xmlDecoder,
	}
}
