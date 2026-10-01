// Command dsa starts a local practice server and opens it in your browser.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"dsa/internal/config"
	"dsa/internal/problems"
	"dsa/internal/runner"
	"dsa/internal/server"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		runDoctor()
		return
	}

	port := flag.Int("port", 7777, "preferred port (a free one is used if it is taken)")
	noOpen := flag.Bool("no-open", false, "do not open the browser")
	flag.Parse()

	home, err := config.Home()
	if err != nil {
		log.Fatalf("cannot find your home folder: %v", err)
	}

	lib, err := problems.Load(problems.Builtin(), os.DirFS(filepath.Join(home, "packs")))
	if err != nil {
		log.Fatalf("cannot load problems: %v", err)
	}

	srv, err := server.New(home, lib)
	if err != nil {
		log.Fatalf("cannot start: %v", err)
	}

	ln, err := listen(*port)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	srv.SetAddr(ln.Addr().String())

	url := "http://" + ln.Addr().String()
	log.Printf("dsa is running at %s (press Ctrl+C to stop)", url)
	log.Printf("your solutions are saved in %s", filepath.Join(home, "solutions"))
	if !*noOpen {
		go func() {
			time.Sleep(50 * time.Millisecond)
			openBrowser(url)
		}()
	}

	httpSrv := &http.Server{Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second}
	if err := httpSrv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// listen binds to loopback only, so nothing outside this computer can connect.
func listen(port int) (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err == nil {
		return ln, nil
	}
	return net.Listen("tcp", "127.0.0.1:0")
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

func runDoctor() {
	fmt.Println("dsa doctor: checking the tools needed to run your code")
	for _, t := range runner.Doctor() {
		if t.Found {
			fmt.Printf("  ok       %-6s %s  (%s)\n", t.Name, t.Version, t.Path)
		} else {
			fmt.Printf("  missing  %-6s install it and make sure it is on your PATH\n", t.Name)
		}
	}
}
