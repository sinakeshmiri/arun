package http

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/sinakeshmiri/arun/internal/ports"
)

type AddReq struct {
	name string `form:"name"`
	src  string `form:"src"`
}

func (httpa Adapter) AddFunction(c *gin.Context) {
	var b AddReq
	c.Bind(&b)
	err := httpa.api.AddFunction(b.name, b.src)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(`{ERR : "%s"}`, err.Error()))
	}

}

func (httpa Adapter) RunFunction(c *gin.Context) {
	// step 1: resolve proxy address, change scheme and host in requets
	req := c.Request
	fName := req.Header["X-Func"][0]
	pod, err := httpa.api.RunFunction(fName)
	if err != nil {
		log.Printf("couldn't run the pod", err)
		c.String(500, "error")
		return
	}
	podUrl := "http://" + pod
	proxy, err := url.Parse(podUrl)
	if err != nil {
		log.Printf("error in parse addr: %v", err)
		c.String(500, "error")
		return
	}
	req.URL.Scheme = proxy.Scheme
	req.URL.Host = proxy.Host
	t1 := time.Now()
	// step 2: use http.Transport to do request to real server.
	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Printf("error in roundtrip: %v", err)
		c.String(500, "error")
		return
	}

	// step 3: return real server response to upstream.
	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Header(k, v)
		}
	}
	defer resp.Body.Close()
	bufio.NewReader(resp.Body).WriteTo(c.Writer)
	t2 := time.Now()
	diff := t2.Sub(t1)
	err = httpa.api.UpdateFunction(fName,diff)
	if err != nil {
		log.Printf("error in updating the record: %v", err)
		c.String(500, "error")
		return
	}
	return
}
