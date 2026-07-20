package network

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/hcagri-ceng/honey-go/internal/models"
	"github.com/hcagri-ceng/honey-go/internal/storage"
)

// Server, TCP sunucusunun yapılandırmasını tutar.
type Server struct {
	address     string
	readTimeout time.Duration
	maxConn     int
	sem         chan struct{}
	repo        storage.Repository // <-- Ekledik
}

// NewServer artık repo da kabul ediyor
func NewServer(address string, maxConcurrent int, repo storage.Repository) *Server {
	return &Server{
		address:     address,
		readTimeout: 5 * time.Second,
		maxConn:     maxConcurrent,
		sem:         make(chan struct{}, maxConcurrent),
		repo:        repo, // <-- Ekledik
	}
}

// Start, sunucuyu ayağa kaldırır ve gelen bağlantıları dinler.
func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	// Uygulama kapandığında portun açık kalmasını (zombie port) engelleriz.
	defer listener.Close()

	log.Printf("Honeygo %s adresinde dinleniyor...", s.address)

	// Graceful shutdown mekanizması: Context iptal edildiğinde listener'ı kapatır.
	go func() {
		<-ctx.Done()
		log.Println("Kapanma sinyali alındı, listener durduruluyor...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil // Normal kapanma durumu
			default:
				log.Printf("Bağlantı kabul hatası: %v", err)
				continue
			}
		}

		// Her bağlantıyı asenkron olarak işliyoruz.
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(s.readTimeout))

	remoteAddr := conn.RemoteAddr().String()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	// 1. IP ve Port'u birbirinden ayır
	host, portStr, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Printf("Adres parçalanırken hata: %v", err)
		return
	}

	sourcePort, _ := strconv.Atoi(portStr)

	// 2. Event modelimizi oluştur
	event := models.Event{
		SourceIP:   host,
		SourcePort: sourcePort,
		TargetPort: 8080, // Dinlediğimiz port
		Protocol:   "TCP",
		Payload:    buffer[:n],
	}

	// 3. Veritabanına kaydet
	// (Gerçek projelerde bunu ayrı bir goroutine'e atarız ki I/O işlemi socketi bloklamasın ama şimdilik doğrudan yazıyoruz)
	if err := s.repo.SaveEvent(context.Background(), event); err != nil {
		log.Printf("DB Kayıt Hatası: %v", err)
	} else {
		log.Printf("[%s] Payload DB'ye kaydedildi (%d byte)", host, n)
	}
}
