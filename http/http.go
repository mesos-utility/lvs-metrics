package http

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"

	"github.com/golang/glog"
	"github.com/mesos-utility/lvs-metrics/g"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// start http server
func Start() {
	go startHttpServer()
}

func configRoutes() {
	configCommonRoutes()
}

func startHttpServer() {
	if !g.Config().Http.Enable {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	// init url mapping
	configRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	glog.Infoln("http.startHttpServer ok, listening ", addr)
	glog.Fatalln(s.ListenAndServe())
}

// WriteJSON writes the value v to the http response stream as json with standard json encoding.
func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, err string) {
	RenderJson(w, map[string]string{"msg": "failed", "data": err})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}
