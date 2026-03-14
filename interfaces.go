package goffmpeg

type ISingleParameter interface {
	GetParameterString() string
}

type IMultiParameter interface {
	GetParameterStrings() []string
}

type IOutput interface {
	ISingleParameter
	// Unexported method to identify output types without relying on type assertions
	isOutput()
}

type IVideoEncoding interface {
	IMultiParameter
}

type IAudioEncoding interface {
	IMultiParameter
}
