package modelsDataEncryptor

const AesGcm128 = 1
const AesGcm192 = 2
const AesGcm256 = 3

var Methods = map[int]MethodParams{
	AesGcm128: {
		KeyLength: 16,
		IvSize:    12,
		BlockSize: 128,
	},
	AesGcm192: {
		KeyLength: 24,
		IvSize:    12,
		BlockSize: 192,
	},
	AesGcm256: {
		KeyLength: 32,
		IvSize:    12,
		BlockSize: 256,
	},
}

type MethodParams struct {
	KeyLength int
	IvSize    int
	BlockSize int
}
