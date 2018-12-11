package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/dags-/orepack/ore"
	"github.com/valyala/fasthttp"
)

const (
	path  = "/com/orepack/:id/:version/:file"
	group = "com.orepack"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	flag.Parse()

	router := fasthttprouter.New()
	router.GET(path, repoHandlerWrapper)
	router.NotFound = notFoundHandler
	server := fasthttp.Server{
		Handler:            router.Handler,
		GetOnly:            true,
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

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("https://ore.spongepowered.org", fasthttp.StatusPermanentRedirect)
}

func repoHandlerWrapper(ctx *fasthttp.RequestCtx) {
	id, version, file := value(ctx, "id"), value(ctx, "version"), value(ctx, "file")

	e := repoHandler(id, version, file, ctx)
	if e != nil {
		ctx.Response.Header.SetStatusCode(http.StatusNotFound)
		fmt.Fprintln(ctx.Response.BodyWriter(), e)
		log.Printf("error (id:%s,version:%s,file:%s): %s\n", id, version, file, e.Error())
	}
}

func repoHandler(id, version, file string, ctx *fasthttp.RequestCtx) error {
	if !strings.HasPrefix(file, id+"-"+version) {
		return http.ErrNoLocation
	}

	switch filepath.Ext(file) {
	case ".pom":
		return pom(ctx, id, version)
	case ".md5":
		return md5(ctx, id, version)
	case ".jar":
		return jar(ctx, id, version)
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

func jar(ctx *fasthttp.RequestCtx, id, version string) error {
	j, e := ore.GetJar(id, version)
	if e != nil {
		return e
	}
	_, e = io.Copy(ctx.Response.BodyWriter(), j)
	return e
}

func md5(ctx *fasthttp.RequestCtx, id, version string) error {
	v, e := ore.GetVersion(id, version)
	if e != nil {
		return e
	}
	_, e = fmt.Fprintln(ctx.Response.BodyWriter(), v.MD5)
	return e
}

func pom(ctx *fasthttp.RequestCtx, id, version string) error {
	v, e := ore.GetVersion(id, version)
	if e != nil {
		return e
	}
	ctx.Response.Header.SetContentType("application/xml")
	en := xml.NewEncoder(ctx.Response.BodyWriter())
	en.Indent("", "  ")
	return en.Encode(ore.NewPom(id, group, v.Name))
}
