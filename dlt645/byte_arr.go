package dlt645

var _ Analyser = (*ByteArrayAnalyzer)(nil)

type ByteArrayAnalyzer struct {
	length int
	Value  []byte
}

func (b *ByteArrayAnalyzer) Clean() {
	b.length = 0
	b.Value = nil
}

func (b *ByteArrayAnalyzer) Decode(_ []byte) error {
	return nil
}

func (b *ByteArrayAnalyzer) Encode() ([]byte, error) {
	return b.Value, nil
}

func (b *ByteArrayAnalyzer) GetValue() interface{} {
	return b.Value
}
