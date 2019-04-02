package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/dags-/orepack/ore"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	flag.Parse()

	router := fasthttprouter.New()
	router.GET("/com/orepack/:owner/:name/:version/:file", repoHandlerWrapper)

	server := fasthttp.Server{
		Handler:            router.Handler,
		DisableKeepalive:   true,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 10,
	}

	go handleStop()

	log.Println("serving on port", *port)
	e := server.ListenAndServe(fmt.Sprintf(":%v", *port))
	if e != nil {
		panic(e)
	}
}

func handleStop() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		cmd := strings.ToLower(strings.TrimSpace(s.Text()))
		if cmd == "stop" {
			log.Println("stopping")
			os.Exit(0)
		}
	}
}

func repoHandlerWrapper(ctx *fasthttp.RequestCtx) {
	owner := value(ctx, "owner")
	name := value(ctx, "name")
	version := value(ctx, "version")
	file := value(ctx, "file")
	log.Println(owner, name, version, file)
	e := repoHandler(owner, name, version, file, ctx)
	if e != nil {
		ctx.Response.Header.SetStatusCode(http.StatusNotFound)
		fmt.Fprintln(ctx.Response.BodyWriter(), e)
		log.Printf("error (id:%s,version:%s,file:%s): %s\n", name, version, file, e.Error())
	}
}

func repoHandler(owner, id, version, file string, ctx *fasthttp.RequestCtx) error {
	if !strings.HasPrefix(file, id+"-"+version) {
		return http.ErrNoLocation
	}

	switch filepath.Ext(file) {
	case ".pom":
		return pom(ctx, owner, id, version)
	case ".md5":
		return md5(ctx, owner, id, version)
	case ".jar":
		return jar(ctx, owner, id, version)
	default:
		return http.ErrNoLocation
	}
}

func value(ctx *fasthttp.RequestCtx, name string) string {
	if str, ok := ctx.UserValue(name).(string); ok {
		return str
	}
	return ""
}

func jar(ctx *fasthttp.RequestCtx, owner, name, version string) error {
	p, e := ore.GetProject(owner, name)
	if e != nil {
		return e
	}
	url, e := ore.GetJarURL(p.ID, version)
	if e != nil {
		return e
	}
	ctx.Redirect(url, fasthttp.StatusOK)
	return e
}

func md5(ctx *fasthttp.RequestCtx, owner, name, version string) error {
	p, e := ore.GetProject(owner, name)
	if e != nil {
		return e
	}
	v, e := ore.GetVersion(p.ID, version)
	if e != nil {
		return e
	}
	_, e = fmt.Fprintln(ctx.Response.BodyWriter(), v.MD5)
	return e
}

func pom(ctx *fasthttp.RequestCtx, owner, name, version string) error {
	p, e := ore.GetProject(owner, name)
	if e != nil {
		return e
	}
	v, e := ore.GetVersion(p.ID, version)
	if e != nil {
		return e
	}
	ctx.Response.Header.SetContentType("application/xml")
	en := xml.NewEncoder(ctx.Response.BodyWriter())
	en.Indent("", "  ")
	return en.Encode(ore.NewPom(owner, name, v.Name))
}
