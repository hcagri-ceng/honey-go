package defense

import (
	"log"
	"os/exec"
	"runtime"
)

// IsolateNetwork bilgisayarın ağ bağlantısını donanımsal düzeyde keser.
func IsolateNetwork() {
	log.Println("[SAVUNMA] Ağ izolasyon protokolü başlatıldı!")

	if runtime.GOOS == "windows" {
		// Windows'ta IP adresini serbest bırakarak interneti kesen komut
		cmd := exec.Command("ipconfig", "/release")
		err := cmd.Run()
		if err != nil {
			log.Printf("[HATA] Ağ kesilemedi: %v\n", err)
			return
		}
		log.Println("💥 BAĞLANTI KESİLDİ! Cihaz karantinaya alındı, yayılım durduruldu.")
	} else {
		log.Println("[UYARI] Bu işletim sistemi için ağ kesme komutu tanımlanmadı.")
	}
}
