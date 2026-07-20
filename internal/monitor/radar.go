package monitor

import (
	"context"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/hcagri-ceng/honey-go/internal/defense"
	"github.com/hcagri-ceng/honey-go/internal/models"
	"github.com/hcagri-ceng/honey-go/internal/storage"
)

// SetupHoneytoken artık repo parametresi de alıyor (Dependency Injection)
func SetupHoneytoken(filePath string, repo *storage.SQLiteRepo) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dummyData := []byte("GİZLİ: 2026 Q3 Finansal Tablolar ve Yönetim Kurulu Kararları...\n")
		if err := os.WriteFile(filePath, dummyData, 0644); err != nil {
			return err
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					log.Printf("!!! MÜDAHALE TESPİT EDİLDİ !!! Dosya: %s | İşlem: %s", event.Name, event.Op.String())

					// 1. Olayı veritabanına uyacak şekilde (Event modeliyle) paketle
					dbEvent := models.Event{
						SourceIP:   "127.0.0.1", // Dosya işlemi lokalde olduğu için
						SourcePort: 0,
						TargetPort: 0,
						Protocol:   "FS_ALERT", // Dashboard'da kırmızı görünmesi için özel bir protokol adı
						Payload:    []byte("HONEYTOKEN İHLALİ: Dosya " + event.Op.String() + " işlemine maruz kaldı!"),
					}

					// 2. Ağı kesmeden hemen önce veritabanına (Dashboard'a) yaz
					if err := repo.SaveEvent(context.Background(), dbEvent); err != nil {
						log.Printf("Honeytoken logu DB'ye yazılamadı: %v", err)
					}

					// 3. Silahı ateşle (Ağı kes)
					defense.IsolateNetwork()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher hatası:", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		return err
	}

	log.Printf("Yem dosya (Honeytoken) izlemeye alındı: %s", filePath)
	return nil
}
