package monitor

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/hcagri-ceng/honey-go/internal/defense"
)

// WatchBait, verilen hedef dosyayı işletim sistemi seviyesinde izlemeye alır.
func WatchBait(filePath string) error {
	// 1. Yeni bir izleyici (watcher) nesnesi oluştur
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher motoru baslatilamadi: %w", err)
	}
	// Fonksiyon ne olursa olsun bittiğinde bellekte sızıntı yapmaması için watcher'ı kapat
	defer watcher.Close()

	// 2. Hedef dosyayı izleme listesine ekle
	err = watcher.Add(filePath)
	if err != nil {
		return fmt.Errorf("dosya izlemeye alinamadi (%s): %w", filePath, err)
	}

	log.Printf("📡 RADAR AKTİF: %s dosyasi 7/24 izleniyor...", filePath)

	// 3. Sonsuz Döngü (Event Loop) - Go'nun gücü burada başlar
	for {
		// select bloğu, Go'da kanalları (channels) dinlemek için kullanılır.
		// İşletim sisteminden bir sinyal gelene kadar kod burada bloklanır, CPU harcamaz.
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Fidye virüsleri (Ransomware) ne yapar?
			// Ya dosyanın adını değiştirir (.encrypted yapar), ya içini değiştirir (Write) ya da siler (Remove).
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
				log.Printf("[🚨 ALARM] Yem dosyaya müdahale tespit edildi! Islem: %s", event.Op)
				log.Println("Siber güvenlik protokolü devreye sokuluyor...")

				defense.IsolateNetwork()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("[HATA] Radar dinleme hatası: %v", err)
		}
	}
}
