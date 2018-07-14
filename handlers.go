package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

// Hit database composition
type Hit struct {
	Hitokoto string `json:"hitokoto"` // Hitokoto sentence
	Source   string `json:"source"`   // Hitokoto source
}

// Hitokoto handle function
func Hitokoto(w http.ResponseWriter, r *http.Request) {
	// query
	var hito string
	var source string
	var content string

	// get params
	encode := r.URL.Query().Get("encode")
	length := r.URL.Query().Get("length")
	// if url param have callback then will ignore encode
	callback := r.URL.Query().Get("callback")
	if length != "" {
		err1 := db.QueryRow("SELECT hitokoto, source FROM main WHERE LENGTH(hitokoto) < ? ORDER BY RAND() LiMIT 1;", length).Scan(&hito, &source)
		if err1 != nil {
			hito = ""
			source = ""
		}
	} else {
		nBig, err := rand.Int(rand.Reader, big.NewInt(AMOUNT))
		n := nBig.Int64()
		err = db.QueryRow("SELECT hitokoto, source FROM main LIMIT ?, 1;", n).Scan(&hito, &source)
		checkErr(err)
	}

	if callback != "" {
		hs := &Hit{hito, source}
		fmtJSON, _ := json.Marshal(hs)
		fmt.Fprintf(w, "%s(%s);", callback, fmtJSON)
	} else {
		if encode == "js" {
			content = fmt.Sprintf("%s——「%s」", hito, source)
			fmt.Fprintf(w, "var hitokoto=%s;var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;", content)
		} else if encode == "json" {
			hs := &Hit{
				hito,
				source,
			}
			fmtJSON, _ := json.Marshal(hs)
			fmt.Fprintf(w, "%s", string(fmtJSON))
		} else if encode == "text" {
			fmt.Fprintf(w, "%s", hito)
		} else {
			content = fmt.Sprintf("%s——「%s」", hito, source)
			fmt.Fprintf(w, "%s", content)
		}
	}
}
