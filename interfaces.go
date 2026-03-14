package goffmpeg

type ISingleParameter interface {
	GetParameters() []string
}

type IMultiParameter interface {
	GetParameters() []string
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
