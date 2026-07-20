### 🍯 Honeygo: Real-Time Honeypot & Active Defense System
Honeygo is a highly concurrent Honeypot and active defense system developed with Go (Golang) to detect, delay, and isolate cyber attackers from critical systems.

### 🚀 Key Features
High-Performance Network Layer: A self-protecting TCP listener built against Resource Exhaustion (DoS) attacks, utilizing Go's goroutine architecture and a buffered channel (semaphore) based fail-fast design.

Fake Service Simulation (Tarpit): Keeps the attacker engaged by mimicking an outdated and vulnerable Apache server to deceive automated scanning tools (Nmap, ZDI, etc.).

FS-Alert & Honeytoken (Bait File): Creates fake documents with critical names (e.g., !000_gizli_finans_raporu.txt) on the file system and performs millisecond-level monitoring using fsnotify.

Physical Network Isolation (Kill-Switch): The moment an attacker interacts with the Honeytoken (read, write, delete), the system instantly terminates the local network (Wi-Fi) connection via an emergency procedure, preventing lateral movement.

Real-Time Dashboard: Collected threat data (IP, Port, Raw Payload) is logged on SQLite and visualized in real-time on a TailwindCSS interface powered by the Fiber framework.

###🛠️ Tech Stack
Backend: Go (Golang), Fiber (REST API), fsnotify

Database: SQLite (CGO-free, modernc.org/sqlite)

Frontend: HTML5, JavaScript (Fetch API), TailwindCSS (CDN)

Architecture: Clean Architecture, Dependency Injection, Worker-Pool/Semaphore

### 📦 Installation & Execution
The project is designed as "Standalone" and requires no external dependencies (Docker, external database setup, etc.).

1. Clone the repository and navigate into the directory:
    git clone https://github.com/hcagri-ceng/honey-go.git
        cd honey-go
2. Download modules / Install dependencies:
    go mod tidy
3. Run the system:
    go run cmd/honeygo/main.go

Dashboard: You can access the real-time monitoring panel at http://localhost:3000.

Honeypot Port: The system lies in wait on port :8080 by default.

### 🎯 Test Scenarios (For Jury Demos)

1. Reconnaissance Lure:
    Send an HTTP request to the system from another terminal. The system will write the log to the database and return a fake response pretending to be an old Apache server.
        curl -v http://localhost:8080
2. Kill-Switch and Isolation: 
    When the program runs, a file named !000_gizli_finans_raporu.txt is generated in the root directory of the project. Try deleting this file or writing a character in it and saving it.

⚠️ WARNING: When you perform this action, your Wi-Fi connection will be immediately terminated for security purposes!