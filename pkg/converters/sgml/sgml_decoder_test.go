package sgml

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

type TestSimpleStruct struct {
	SGMLName string `sgml:"Root"`
	Text1    string `sgml:"Text1"`
	Text2    string `sgml:"Text2"`
}

type TestNestedStruct1 struct {
	SGMLName string           `sgml:"Root"`
	Child    TestSimpleStruct `sgml:"Child"`
	Text3    string           `sgml:"Text3"`
	Text4    string           `sgml:"Text4"`
}

type TestNestedStruct2 struct {
	SGMLName string            `sgml:"Root"`
	Child    *TestSimpleStruct `sgml:"Child"`
	Text3    string            `sgml:"Text3"`
	Text4    string            `sgml:"Text4"`
}

type TestEmbeddedStruct struct {
	TestSimpleStruct
	Text5 string `sgml:"Text5"`
	Text6 string `sgml:"Text6"`
}

type TestSliceStruct1 struct {
	SGMLName string             `sgml:"Root"`
	Children []TestSimpleStruct `sgml:"Child"`
	Text7    string             `sgml:"Text7"`
}

type TestSliceStruct2 struct {
	SGMLName string              `sgml:"Root"`
	Children []*TestSimpleStruct `sgml:"Child"`
	Text7    string              `sgml:"Text7"`
}

type TestSimpleStructWithXMLTag struct {
	XMLName xml.Name `xml:"Root"`
	Text1   string   `xml:"Text1"`
	Text2   string   `xml:"Text2"`
}

type TestStructWithXMLTag struct {
	XMLName xml.Name                   `xml:"Root"`
	Child   TestSimpleStructWithXMLTag `xml:"Child"`
	Text3   string                     `xml:"Text3"`
	Text4   string                     `xml:"Text4"`
}

type TestNotExportedFieldStruct struct {
	SGMLName string `sgml:"Root"`
	Text1    string `sgml:"Text1"`
	Text2    string
	text3    string `sgml:"Text3"`
}

type TestUnsupportedStruct struct {
	SGMLName string `sgml:"Root"`
	Number   int    `sgml:"Number"`
}

type TestEmbeddedUnsupportedStruct struct {
	TestUnsupportedStruct
	Text1 string `sgml:"Text1"`
}

func TestDecoderDecode(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Text1>Foo\n" +
			"<Text2>Bar\n" +
			"</Root>\n"))

	testStruct := &TestSimpleStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, "Foo", testStruct.Text1)
	assert.Equal(t, "Bar", testStruct.Text2)
}

func TestDecoderDecode_WithRedundantFields(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Text1>Foo\n" +
			"<Text2>Bar\n" +
			"<Text3>Hello\n" +
			"<Child>\n" +
			"<Text4>World\n" +
			"</Child>\n" +
			"</Root>\n"))

	testStruct := &TestSimpleStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, "Foo", testStruct.Text1)
	assert.Equal(t, "Bar", testStruct.Text2)
}

func TestDecoderDecode_WithEndElement(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Text1>Foo</Text1>\n" +
			"<Text2>Bar</Text2>\n" +
			"</Root>\n"))

	testStruct := &TestSimpleStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, "Foo", testStruct.Text1)
	assert.Equal(t, "Bar", testStruct.Text2)
}

func TestDecoderDecode_WithoutBreakLine(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>" +
			"<Text1>Foo" +
			"<Text2>Bar" +
			"</Root>"))

	testStruct := &TestSimpleStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, "Foo", testStruct.Text1)
	assert.Equal(t, "Bar", testStruct.Text2)
}

func TestDecoderDecode_NestedStruct(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Child>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"</Child>\n" +
			"<Text3>Foo\n" +
			"<Text4>Bar\n" +
			"</Root>\n"))

	testStruct := &TestNestedStruct1{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.NotNil(t, testStruct.Child)
	assert.Equal(t, "Hello", testStruct.Child.Text1)
	assert.Equal(t, "World", testStruct.Child.Text2)
	assert.Equal(t, "Foo", testStruct.Text3)
	assert.Equal(t, "Bar", testStruct.Text4)
}

func TestDecoderDecode_NestedStructUsingPointer(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Child>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"</Child>\n" +
			"<Text3>Foo\n" +
			"<Text4>Bar\n" +
			"</Root>\n"))

	testStruct := &TestNestedStruct2{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.NotNil(t, testStruct.Child)
	assert.Equal(t, "Hello", testStruct.Child.Text1)
	assert.Equal(t, "World", testStruct.Child.Text2)
	assert.Equal(t, "Foo", testStruct.Text3)
	assert.Equal(t, "Bar", testStruct.Text4)
}

func TestDecoderDecode_EmbeddedStruct(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"<Text5>Foo\n" +
			"<Text6>Bar\n" +
			"</Root>\n"))

	testStruct := &TestEmbeddedStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, "Hello", testStruct.Text1)
	assert.Equal(t, "World", testStruct.Text2)
	assert.Equal(t, "Foo", testStruct.Text5)
	assert.Equal(t, "Bar", testStruct.Text6)
}

func TestDecoderDecode_StructSlice(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Child>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"</Child>\n" +
			"<Child>\n" +
			"<Text1>Hello2\n" +
			"<Text2>World2\n" +
			"</Child>\n" +
			"<Text7>Foo\n" +
			"</Root>\n"))

	testStruct := &TestSliceStruct1{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, 2, len(testStruct.Children))
	assert.Equal(t, "Hello", testStruct.Children[0].Text1)
	assert.Equal(t, "World", testStruct.Children[0].Text2)
	assert.Equal(t, "Hello2", testStruct.Children[1].Text1)
	assert.Equal(t, "World2", testStruct.Children[1].Text2)
	assert.Equal(t, "Foo", testStruct.Text7)
}

func TestDecoderDecode_StructSliceUsingPointer(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Child>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"</Child>\n" +
			"<Child>\n" +
			"<Text1>Hello2\n" +
			"<Text2>World2\n" +
			"</Child>\n" +
			"<Text7>Foo\n" +
			"</Root>\n"))

	testStruct := &TestSliceStruct2{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, 2, len(testStruct.Children))
	assert.Equal(t, "Hello", testStruct.Children[0].Text1)
	assert.Equal(t, "World", testStruct.Children[0].Text2)
	assert.Equal(t, "Hello2", testStruct.Children[1].Text1)
	assert.Equal(t, "World2", testStruct.Children[1].Text2)
	assert.Equal(t, "Foo", testStruct.Text7)
}

func TestDecoderDecode_UsingXMLTag(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Child>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"</Child>\n" +
			"<Text3>Foo\n" +
			"<Text4>Bar\n" +
			"</Root>\n"))

	testStruct := &TestStructWithXMLTag{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.NotNil(t, testStruct.Child)
	assert.Equal(t, "Hello", testStruct.Child.Text1)
	assert.Equal(t, "World", testStruct.Child.Text2)
	assert.Equal(t, "Foo", testStruct.Text3)
	assert.Equal(t, "Bar", testStruct.Text4)
}

func TestDecoderDecode_WithNotExportedFields(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Text1>Foo\n" +
			"<Text2>Bar\n" +
			"<Text3>Hello World\n" +
			"</Root>\n"))

	testStruct := &TestNotExportedFieldStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.Nil(t, err)
	assert.NotNil(t, testStruct)
	assert.Equal(t, "Foo", testStruct.Text1)
	assert.Equal(t, "", testStruct.Text2)
	assert.Equal(t, "", testStruct.text3)
}

func TestDecoderDecode_StructWithoutEndElement(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Text1>Foo\n" +
			"<Text2>Bar\n"))

	testStruct := &TestSimpleStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.EqualError(t, err, errs.ErrInvalidSGMLFile.Message)

	sgmlDecoder = NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Child>\n" +
			"<Text1>Hello\n" +
			"<Text2>World\n" +
			"<Text3>Foo\n" +
			"<Text4>Bar\n" +
			"</Root>\n"))

	testStruct2 := &TestNestedStruct2{}
	err = sgmlDecoder.Decode(&testStruct2)

	assert.EqualError(t, err, errs.ErrInvalidSGMLFile.Message)
}

func TestDecoderDecode_WithNotSupportedField(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Number>1234\n" +
			"</Root>\n"))

	testStruct := &TestUnsupportedStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.EqualError(t, err, errs.ErrInvalidSGMLFile.Message)
}

func TestDecoderDecode_WithEmbeddedNotSupportedField(t *testing.T) {
	sgmlDecoder := NewDecoder(strings.NewReader(
		"<Root>\n" +
			"<Number>1234\n" +
			"<Text1>Foo\n" +
			"</Root>\n"))

	testStruct := &TestEmbeddedUnsupportedStruct{}
	err := sgmlDecoder.Decode(&testStruct)

	assert.EqualError(t, err, errs.ErrInvalidSGMLFile.Message)
}
