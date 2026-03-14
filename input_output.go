package goffmpeg

//////////
// FileBase
//////////

type FileBase struct {
	FilePath string
}

func (f *FileBase) GetParameters() []string {
	return []string{f.FilePath}
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

func (i *InputFile) GetParameters() []string {
	params := []string{}
	if i.Format != "" {
		params = append(params, "-f", i.Format)
		params = append(params, "-safe", Ternary(i.Save, "1", "0"))
	}
	params = append(params, "-i")
	params = append(params, i.FileBase.GetParameters()...)
	return params
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

func (n *NullOutput) GetParameters() []string {
	return []string{"-f", "null", "NUL"}
}

func (n *NullOutput) isOutput() {}
