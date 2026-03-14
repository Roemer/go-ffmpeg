package goffmpeg

import (
	"fmt"
	"strings"
)

//////////
// FileBase
//////////

type FileBase struct {
	FilePath string
}

func (f *FileBase) GetParameterString() string {
	return fmt.Sprintf("%q", f.FilePath)
}

//////////
// InputFile
//////////

type InputFile struct {
	FileBase
	Index  int
	Format string
	Save   bool
}

func NewInputFile(filePath string) *InputFile {
	return &InputFile{
		FileBase: FileBase{FilePath: filePath},
		Index:    -1,
	}
}

func (i *InputFile) UseConcatDemuxer(save bool) *InputFile {
	i.Format = "concat"
	i.Save = save
	return i
}

func (i *InputFile) GetParameterString() string {
	sb := strings.Builder{}
	if i.Format != "" {
		fmt.Fprintf(&sb, "-f %s ", i.Format)
		fmt.Fprintf(&sb, "-safe %s ", Ternary(i.Save, "1", "0"))
	}
	fmt.Fprintf(&sb, "-i %s", i.FileBase.GetParameterString())
	return sb.String()
}

//////////
// OutputFile
//////////

type OutputFile struct {
	FileBase
}

func NewOutputFile(filePath string) *OutputFile {
	return &OutputFile{FileBase: FileBase{FilePath: filePath}}
}

func (o *OutputFile) isOutput() {}

//////////
// NullOutput
//////////

type NullOutput struct{}

func (n *NullOutput) GetParameterString() string {
	return "-f null NUL"
}

func (n *NullOutput) isOutput() {}
