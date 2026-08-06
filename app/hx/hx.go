package hx

import "net/http"

func Redirect(w http.ResponseWriter, target string) {
	w.Header().Add("HX-Redirect", target)
}

func Retarget(w http.ResponseWriter, target string) {
	w.Header().Add("HX-Retarget", target)
}

func Reswap(w http.ResponseWriter, swap string) {
	w.Header().Add("HX-Reswap", swap)
}
