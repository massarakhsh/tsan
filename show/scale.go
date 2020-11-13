package show

import (
	"github.com/massarakhsh/lik"
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
)

//	Масштабирование изображения в файле
func MakeScaleJpg(filesrc string, width int, height int, filetrg string) bool {
	ok := false
	if file, err := os.Open(filesrc); err != nil {
		lik.SayError("Error 1 scale JPG")
	} else if img, err := jpeg.Decode(file); err != nil {
		lik.SayError("Error 2 scale JPG")
		file.Close()
	} else {
		file.Close()
		m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
		if out, err := os.Create(filetrg); err != nil {
			lik.SayError("Error 3 scale JPG")
		} else {
			jpeg.Encode(out, m, nil)
			out.Close()
			ok = true
		}
	}

	return ok
}

