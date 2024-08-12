package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type forwardProxy struct {
}

func (p *forwardProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  // The "Host:" header is promoted to Request.Host and is removed from
  // request.Header by net/http, so we print it out explicitly.
  log.Println(req.RemoteAddr, "\t\t", req.Method, "\t\t", req.URL, "\t\t Host:", req.Host)
  log.Println("\t\t\t\t\t", req.Header)

  if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
    msg := "unsupported protocal scheme " + req.URL.Scheme
    http.Error(w, msg, http.StatusBadRequest)
    log.Println(msg)
    return
  }

  client := &http.Client{}
  // When a http.Request is sent through an http.Client, RequestURI should not
  // be set (see documentation of this field).
  req.RequestURI = ""

  removeHopHeaders(req.Header)
  removeConnectionHeaders(req.Header)

  if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
    appendHostToXForwardHeader(req.Header, clientIP)
  }

  resp, err := client.Do(req)
  if err != nil {
    http.Error(w, "Server Error", http.StatusInternalServerError)
    log.Fatal("ServeHTTP:", err)
  }
  defer resp.Body.Close()

  log.Println(req.RemoteAddr, " ", resp.Status)

  removeHopHeaders(resp.Header)
  removeConnectionHeaders(resp.Header)

  copyHeader(w.Header(), resp.Header)
  w.WriteHeader(resp.StatusCode)
  io.Copy(w, resp.Body)
}

func main() {
  var addr = flag.String("addr", "127.0.0.1:9999", "proxy address")
  flag.Parse()

  proxy := &forwardProxy{}

  log.Println("Starting proxy server on", *addr)
  if err := http.ListenAndServe(*addr, proxy); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
