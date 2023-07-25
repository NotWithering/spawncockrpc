package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var (
	cooldown int = 0
)

func runCommand(command string) (*exec.Cmd, error) {
	cmd := exec.Command("bash", "-c", command)

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	return cmd, nil
}

func main() {
	go func() {
		for {
			if cooldown != 0 {
				time.Sleep(time.Second * 1)
				cooldown -= 1
			}
		}
	}()

	http.HandleFunc("/cockroach", func(w http.ResponseWriter, r *http.Request) {
		if cooldown == 0 {
			cmd, err := runCommand("wine 'C:\\Program Files (x86)\\Cockroach on Desktop\\CockroachOnDesktop.exe'")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				serveHTMLFile(w, "failed.html")
				return
			}

			go func() {
				err := cmd.Wait()
				if err != nil {
					fmt.Println("Command execution failed:", err)
				}
			}()

			cooldown = 5
			w.WriteHeader(http.StatusOK)
			serveHTMLFile(w, "accepted.html")
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			serveHTMLFile(w, "failed.html")
		}
	})

	http.ListenAndServe(":8080", nil)
}

func serveHTMLFile(w http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open %s: %v", filename, err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write %s: %v", filename, err), http.StatusInternalServerError)
	}
}
