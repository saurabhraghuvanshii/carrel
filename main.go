// Command carrel starts a local practice server and opens it in your browser.
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
	"strings"
	"time"

	"github.com/saurabhraghuvanshii/carrel/internal/config"
	"github.com/saurabhraghuvanshii/carrel/internal/problems"
	"github.com/saurabhraghuvanshii/carrel/internal/runner"
	"github.com/saurabhraghuvanshii/carrel/internal/server"
	"github.com/saurabhraghuvanshii/carrel/internal/store"
)

// version is set at release time with -ldflags "-X main.version=v0.1.0".
var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version":
			fmt.Println("carrel", version)
			return
		case "doctor":
			runDoctor()
			return
		case "export":
			if err := runExport(os.Args[2:]); err != nil {
				log.Fatal(err)
			}
			return
		case "import":
			if err := runImport(os.Args[2:]); err != nil {
				log.Fatal(err)
			}
			return
		}
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
	log.Printf("Carrel is running at %s (press Ctrl+C to stop)", url)
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
	fmt.Printf("carrel doctor: checking the tools needed to run your code (carrel %s)\n", version)
	for _, t := range runner.Doctor() {
		if t.Found {
			fmt.Printf("  ok       %-6s %s  (%s)\n", t.Name, t.Version, t.Path)
		} else {
			fmt.Printf("  missing  %-6s install it and make sure it is on your PATH\n", t.Name)
		}
	}
}

// runExport writes all solutions and progress to a zip: carrel export [file.zip].
func runExport(args []string) error {
	if len(args) > 1 {
		return errors.New("usage: carrel export [file.zip]")
	}
	name := "carrel-solutions-" + time.Now().Format("2006-01-02") + ".zip"
	if len(args) == 1 {
		name = args[0]
	}
	st, err := openStore()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", name, err)
	}
	if err := st.ExportZip(f); err != nil {
		f.Close()
		os.Remove(name)
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	fmt.Printf("saved your solutions and progress to %s\n", name)
	return nil
}

// runImport reads a zip made by export: carrel import <file.zip> [--overwrite].
func runImport(args []string) error {
	overwrite, name := false, ""
	for _, a := range args {
		switch {
		case a == "--overwrite" || a == "-overwrite":
			overwrite = true
		case name == "" && !strings.HasPrefix(a, "-"):
			name = a
		default:
			return errors.New("usage: carrel import <file.zip> [--overwrite]")
		}
	}
	if name == "" {
		return errors.New("usage: carrel import <file.zip> [--overwrite]")
	}
	st, err := openStore()
	if err != nil {
		return err
	}
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	rep, err := st.ImportZip(f, info.Size(), overwrite)
	if err != nil {
		return err
	}
	fmt.Printf("imported %d, skipped %d, rejected %d\n", len(rep.Imported), len(rep.Skipped), len(rep.Rejected))
	for _, n := range rep.Imported {
		fmt.Println("  imported  ", n)
	}
	for _, n := range rep.Skipped {
		fmt.Println("  skipped   ", n, "(already there; use --overwrite to replace)")
	}
	for _, r := range rep.Rejected {
		fmt.Printf("  rejected   %s (%s)\n", r.Name, r.Reason)
	}
	if rep.ProgressMerged > 0 {
		fmt.Printf("  progress   %d entries added or updated\n", rep.ProgressMerged)
	}
	return nil
}

func openStore() (*store.Store, error) {
	home, err := config.Home()
	if err != nil {
		return nil, fmt.Errorf("cannot find your home folder: %w", err)
	}
	return store.New(home)
}
