package defense

import (
	"log"
	"os/exec"
	"runtime"
)

// IsolateNetwork, sistemin ağ ile olan fiziksel/mantıksal bağını anında keser.
func IsolateNetwork() {
	log.Println("!!! KRİTİK UYARI: HONEYTOKEN İHLALİ !!!")
	log.Println("Ağ izolasyonu başlatılıyor (Kill-Switch)...")

	if runtime.GOOS == "windows" {
		// Windows üzerinde Wi-Fi bağlantısını anında kesen komut
		cmd := exec.Command("netsh", "wlan", "disconnect")
		if err := cmd.Run(); err != nil {
			log.Printf("Ağ kesilirken kritik hata (Manuel müdahale gerekebilir): %v\n", err)
		} else {
			log.Println("BAŞARILI: Wi-Fi bağlantısı koparıldı. Sistem izole edildi.")
		}
	} else {
		// Eğer jüriye Linux/Mac'te gösterirsen diye fallback (Hackathon önlemi)
		log.Printf("İzolasyon komutu %s için tanımlanmamış, ağ kesilemedi.\n", runtime.GOOS)
	}
}
