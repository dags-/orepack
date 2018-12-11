package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/dags-/orepack/ore"
	"github.com/valyala/fasthttp"
)

func main() {
	port := flag.String("port", "8080", "server port")

	router := fasthttprouter.New()
	router.GET("/com/orepack/:id/:version/:file", repoHandler)
	router.NotFound = notFoundHandler
	server := fasthttp.Server{
		Handler:            router.Handler,
		GetOnly:            true,
		DisableKeepalive:   true,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 10,
	}

	go handleStop()

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
			os.Exit(0)
		}
	}
}

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("https://ore.spongepowered.org", fasthttp.StatusPermanentRedirect)
}

func repoHandler(ctx *fasthttp.RequestCtx) {
	defer ctx.SetConnectionClose()
	e := handler(ctx)
	if e != nil {
		ctx.Response.Header.SetStatusCode(http.StatusNotFound)
		fmt.Fprintln(ctx.Response.BodyWriter(), e)
	}
}

func handler(ctx *fasthttp.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	version := ctx.UserValue("version").(string)
	file := ctx.UserValue("file").(string)

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
	return en.Encode(ore.NewPom(id, v))
}
