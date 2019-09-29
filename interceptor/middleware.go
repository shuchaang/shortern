package interceptor

import (
	"log"
	"net/http"
	"time"
)

type MiddleWare struct {
}


//logging interceptor
func (m MiddleWare)LoggingHandler(next http.Handler)http.Handler{
	fn:= func(w http.ResponseWriter,r *http.Request) {
		t1:=time.Now()
		next.ServeHTTP(w,r)
		t2:=time.Now()
		log.Printf("[%s] %q %v ",r.Method,r.URL.String(),t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

//recover interceptor
func (m MiddleWare)RecoverHandler(next http.Handler)http.Handler{
	fn:= func(w http.ResponseWriter,r *http.Request) {
		defer func() {
			if e:=recover();e!=nil{
				log.Printf("recover from panic %+v",e)
				http.Error(w,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w,r)
	}
	return http.HandlerFunc(fn)
}
