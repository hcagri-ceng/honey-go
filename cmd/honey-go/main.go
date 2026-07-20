package main

import (
	"log"

	"github.com/hcagri-ceng/honey-go/internal/honeypot"
	"github.com/hcagri-ceng/honey-go/internal/monitor"
)

func main() {
	log.Println("🛡️ Honey-Go Siber Savunma Sistemi Başlatılıyor...")

	fileName := "!000_gizli_finans_raporu.txt"

	err := honeypot.CreateBaitFiles()
	if err != nil {
		log.Fatalf("Sistem başlatılamadı (Yem hatası): %v", err)
	}

	log.Println("✅ Yem dosya başarıyla oluşturuldu.")
	err = monitor.WatchBait(fileName)
	if err != nil {
		log.Fatalf("Sistem çöktü (Radar hatası): %v", err)
	}
}
