package honeypot

import (
	"os"
)

func CreateBaitFiles() {
	file, err := os.Create("gizli_finans_raporu.txt")
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write([]byte("Bu bir honeypot dosyasıdır"))
	if err != nil {
		return
	}
}
