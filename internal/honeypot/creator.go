package honeypot

import (
	"fmt"
	"os"
)

func CreateBaitFiles() error {
	fileName := "!000_gizli_finans_raporu.txt"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("yem dosya olusturulamadi: %w", err)
	}
	defer file.Close()
	_, err = file.Write([]byte("Bu bir honeypot dosyasıdır. Yetkisiz erisim loglanmaktadir."))
	if err != nil {
		return fmt.Errorf("yem dosyaya icerik yazilamadi: %w", err)
	}
	return nil
}
