package pdf

import (
	"errors"

	"github.com/one-eyed-willy/pkg/utils"
)

func ValidatePdf(data []byte) error {
	isInitialPdfBytes := utils.IsSubSlice(data[:5], []byte{0x25, 0x50, 0x44, 0x46, 0x2D})

	if !isInitialPdfBytes {
		return errors.New("This is not a pdf file")
	}

	if isVersion1dot3(data) {
		eof1dot3 := []byte{
			0x25, // %
			0x25, // %
			0x45, // E
			0x4F, // O
			0x46, // F
			0x20, // SPACE
			0x0A, // EOL
		}

		if utils.IsSubSlice(data[len(data)-7:], eof1dot3) {
			return nil
		}
		return errors.New("Invalid file terminator pdf v1.3")
	}

	if isVersion1dot4(data) {
		eof1dot4 := []byte{
			0x25, // %
			0x25, // %
			0x45, // E
			0x4F, // O
			0x46, // F
			0x0A, // EOL
		}

		if utils.IsSubSlice(data[len(data)-6:], eof1dot4) {
			return nil
		}
		return errors.New("Invalid file terminator for pdf v1.4")
	}

	return nil
}

func isVersion1dot3(data []byte) bool {
	return utils.IsSubSlice(data[5:8], []byte{0x31, 0x2E, 0x33})
}

func isVersion1dot4(data []byte) bool {
	return utils.IsSubSlice(data[5:8], []byte{0x31, 0x2E, 0x34})
}
