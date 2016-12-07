package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elemc/msikeyboard/gomsikeyboard"
)

// Server is a main struct for http server instance
type Server struct {
	Addr string
}

// Start functions for start http server
func (s *Server) Start() (err error) {
	http.HandleFunc("/set", s.handlerSet)
	http.HandleFunc("/test", s.handlerTest)
	err = http.ListenAndServe(s.Addr, nil)
	return
}

// Close function for close resource when server stopped
func (s *Server) Close() {

}

func (s *Server) handlerSet(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		defer r.Body.Close()
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	led := gomsikeyboard.LEDSetting{}

	mode := r.Form.Get("mode")
	if mode == "" {
		mode = "normal"
	}
	led.Mode = mode

	led.Left.Color = r.Form.Get("left-color")
	led.Left.Intensity = r.Form.Get("left-intensity")
	led.Middle.Color = r.Form.Get("middle-color")
	led.Middle.Intensity = r.Form.Get("middle-intensity")
	led.Right.Color = r.Form.Get("right-color")
	led.Right.Intensity = r.Form.Get("right-intensity")

	if r.Form.Get("all-color") != "" {
		led.Left.Color = r.Form.Get("all-color")
		led.Middle.Color = r.Form.Get("all-color")
		led.Right.Color = r.Form.Get("all-color")
	}
	if r.Form.Get("all-intensity") != "" {
		led.Left.Intensity = r.Form.Get("all-intensity")
		led.Middle.Intensity = r.Form.Get("all-intensity")
		led.Right.Intensity = r.Form.Get("all-intensity")
	}
	if r.Form.Get("theme") != "" {
		led, err = gomsikeyboard.GetTheme(r.Form.Get("theme"))
	}

	err = led.Check()
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = led.Set()
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := s.getLEDData(&led)
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) handlerTest(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		defer r.Body.Close()
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
	log.Printf("DEBUG: test handler running for %s on %s", r.RemoteAddr, r.RequestURI)
}

func (s *Server) getLEDData(led *gomsikeyboard.LEDSetting) (data []byte, err error) {
	data, err = json.MarshalIndent(led, "", "\t")
	return
}