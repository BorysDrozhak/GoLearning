//package main
//
//import ( "net/http/httptest"
//	"testing"
//	"net/http"
//	"fmt"
//)
//
//func Test_v0_ping (t *testing.T) {
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintln(w, "Hello, client")
//	}))
//	defer ts.Close()
//
//}