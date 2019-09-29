package main

import (
	"encoding/json"
	"fmt"
	"github.com/justinas/alice"
	"log"
	"net/http"
	error2 "shortern/error"
	"shortern/interceptor"
	"shortern/storage"
)


type shortReq struct {
	URL                   string `json:"url" `
	ExpirationInMinutes int64  `json:"expire_time"`
}

type shortLinkResp struct {
	ShortLink string `json:"shortlink"`
}


var Client *storage.RedisClient

func shorter(writer http.ResponseWriter, request *http.Request){
	var req shortReq
	if err:=json.NewDecoder(request.Body).Decode(&req);err!=nil{
		responseWithError(writer,error2.StatusError{
			http.StatusBadRequest,
			fmt.Errorf("bad param :%v",request.Body),
		})
		request.Body.Close()
	}
	defer request.Body.Close()
	s, e := Client.Shorten(req.URL, req.ExpirationInMinutes)
	if e!=nil{
		responseWithError(writer,e)
	}else{
		responseWithJson(writer,http.StatusCreated,s)
	}
}

func info(writer http.ResponseWriter, request *http.Request){
	query := request.URL.Query()
	s := query.Get("shortlink")
	shortInfo, e := Client.ShortInfo(s)
	if e!=nil{
		responseWithError(writer,e)
	}else{
		responseWithJson(writer,http.StatusCreated,shortInfo)
	}
}

func addr(writer http.ResponseWriter, request *http.Request){
	query := request.URL.Query()
	s := query.Get("addr")
	shorten, e := Client.UnShorten(s)
	if e!=nil{
		responseWithError(writer,e)
	}else{
		responseWithJson(writer,http.StatusCreated,shorten)
	}
}

func main() {
	Client = storage.NewRedisClient()
	m:= interceptor.MiddleWare{}
	chain := alice.New(m.LoggingHandler, m.RecoverHandler)
	http.Handle("/api/shorten",chain.ThenFunc(shorter))
	http.Handle("/api/info",chain.ThenFunc(info))
	http.Handle("/api/addr",chain.ThenFunc(addr))
	http.ListenAndServe(":8081",nil)
}

func responseWithError(writer http.ResponseWriter, err error) {
	switch e := err.(type) {
		case error2.MyError:
			log.Printf("HTTP  %d -- %s",e.Status,e.Error())
			responseWithJson(writer,http.StatusInternalServerError,http.StatusText(http.StatusInternalServerError))
		default:
			responseWithJson(writer,http.StatusInternalServerError,http.StatusText(http.StatusInternalServerError))
	}
}

func responseWithJson(writer http.ResponseWriter, code int, payload interface{}) {
	resp,_:= json.Marshal(payload)
	writer.Header().Set("Content-Type","application/json")
	writer.WriteHeader(code)
	writer.Write(resp)
}
