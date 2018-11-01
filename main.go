package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	serviceMode = flag.Bool("s", false, "enable route service mode")
)

func goToSleep() {
	var sleepLen int64
	if os.Getenv("SLEEP_INTERVAL") != "" {
		s, err := strconv.Atoi(os.Getenv("SLEEP_INTERVAL"))
		if err != nil {
			fmt.Printf("error Parsing $SLEEP_INTERVAL: %s\n", err)
			sleepLen = int64(rand.Intn(10))
		} else {
			sleepLen = int64(s)
		}
	}
	time.Sleep(time.Duration(sleepLen) * time.Second)
	fmt.Printf("Slept for %d seconds\n", sleepLen)
}

func routeServiceHandler(w http.ResponseWriter, r *http.Request) {
	origReq := r.Header.Get("X-CF-Forwarded-Url")
	if origReq == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{ \"error\": \"X-CF-Forwarded-Url Header missing\"}")))
		return
	}
	tr := &http.Transport{
		//MaxIdleConns:       10,
		//IdleConnTimeout:    30 * time.Second,
		//DisableCompression: true,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	newReq, err := http.NewRequest(r.Method, origReq, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{ \"error\": \"New Request: %s\"}", err)))
		return
	}
	newReq.Header = r.Header
	newReq.Header.Add("X-DanL-Route-Service", "welcome to route services")
	//goToSleep()
	resp, err := client.Do(newReq)
	//fmt.Printf("%v\n", newReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{ \"ierror\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{ \"jerror\": \"%s\"}", err)))
		return
	}
	goToSleep()
	w.Write(body)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Reply struct {
		Message string `json:"message"`
	}
	jdata, _ := json.Marshal(Reply{"hello"})
	if r.Header.Get("X-DanL-Route-Service") != "" {
		jdata, _ = json.Marshal(Reply{"hello " + r.Header.Get("X-DanL-Route-Service")})
	}
	goToSleep()
	w.Write(jdata)
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if *serviceMode {
		http.HandleFunc("/", routeServiceHandler)
	} else {
		http.HandleFunc("/", rootHandler)
	}
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
