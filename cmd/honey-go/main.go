package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hcagri-ceng/honey-go/internal/api"
	"github.com/hcagri-ceng/honey-go/internal/monitor"
	"github.com/hcagri-ceng/honey-go/internal/network"
	"github.com/hcagri-ceng/honey-go/internal/storage"
)

func main() {
	log.Println("Honeygo başlatılıyor...")

	// 1. ÖNCE VERİTABANI BAĞLANTISI KURULMALI (repo burada tanımlanıyor)
	repo, err := storage.NewSQLiteRepo("honeygo.db")
	if err != nil {
		log.Fatalf("Veritabanı başlatılamadı: %v", err)
	}
	log.Println("SQLite veritabanı bağlandı.")

	// 2. SONRA HONEYTOKEN SİSTEMİ BAŞLATILMALI (Çünkü artık elimizde repo var)
	tokenFile := "!000_gizli_finans_raporu.txt"
	if err := monitor.SetupHoneytoken(tokenFile, repo); err != nil {
		log.Fatalf("Honeytoken sistemi başlatılamadı: %v", err)
	}

	// 3. KAPANIŞ SİNYALLERİ (Graceful Shutdown)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	// API Sunucusunu hazırlıyoruz
	apiServer := api.NewServer(repo)

	// TCP Honeypot Sunucusunu hazırlıyoruz
	tcpServer := network.NewServer(":8080", 100, repo)

	// 1. API Sunucusunu kendi goroutine'i içinde başlat (3000 portunda çalışsın)
	go func() {
		log.Println("Dashboard API :3000 adresinde başlatılıyor...")
		if err := apiServer.Start(":3000"); err != nil {
			log.Printf("API Sunucusu kapandı: %v", err)
		}
	}()

	// 2. TCP Sunucusunu kendi goroutine'i içinde başlat
	go func() {
		if err := tcpServer.Start(ctx); err != nil {
			log.Printf("TCP Sunucusu hatası: %v", err)
		}
	}()

	// 3. Sistem sinyallerini bekle (Kapanış prosedürü)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Sistem sinyali alındı: %v. Kapanış başlıyor...", sig)

	cancel()             // TCP sunucusunu durdurur
	apiServer.Shutdown() // Fiber sunucusunu durdurur

	log.Println("Honeygo temiz bir şekilde kapatıldı.")
}
