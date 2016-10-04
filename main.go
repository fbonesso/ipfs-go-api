package main

import (
  "os"
  "os/exec"
  "io"
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/dchest/uniuri"
)

func main() {
  route := mux.NewRouter()

  /*
  * Add file to IPFS
  */
  route.HandleFunc("/ipfs/add", func(w http.ResponseWriter, r *http.Request) {
    // CORS
    w.Header().Set("Access-Control-Allow-Origin", "localhost")

    // Form
    r.ParseMultipartForm(32 << 20)

    file, _, err := r.FormFile("file")
    if err != nil {
      fmt.Println(err)
      return
    }
    defer file.Close()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    // UUID
    uuid := uniuri.NewLen(10)

    // Save
    dir, err := os.Create("tmp/" + uuid)
    defer dir.Close()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    if _, err := io.Copy(dir, file); err != nil {
      return
    }

    // IPFS
    var ipfs []byte

    ipfs, _ = exec.Command("curl", []string{"-F", "file=@tmp/" + uuid, "http://localhost:5001/api/v0/add"}...).Output()

    // Remove
    os.Remove("tmp/" + uuid)

    // Response
    w.Write(ipfs)
  })

  /*
  * Execute command on IPFS
  */
  route.HandleFunc("/ipfs/execute/{command}", func(w http.ResponseWriter, r *http.Request) {
    // CORS
    w.Header().Set("Access-Control-Allow-Origin", "localhost")

    // Params
    params := mux.Vars(r)

    // Command
    command := params["command"]

    // IPFS
    var ipfs []byte

    ipfs, _ = exec.Command("ipfs", []string{command}...).Output()

    // Response
    w.Write(ipfs)
  })

  log.Fatal(http.ListenAndServe(":8080", route))
}
