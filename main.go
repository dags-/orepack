package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jackwhelpton/fasthttp-routing"
	"github.com/jackwhelpton/fasthttp-routing/file"
	"github.com/valyala/fasthttp"

	"github.com/dags-/orepack/ore"
)

func main() {
	port := flag.Int("port", 8082, "server port")
	flag.Parse()

	router := routing.New()
	router.Get("/", file.Content("assets/index.html"))
	router.Get("/index.html", file.Content("assets/index.html"))
	router.Get("/script.js", file.Content("assets/script.js"))
	router.Get("/style.css", file.Content("assets/style.css"))
	router.Get("/com/orepack/<owner>/<project>/<version>/<filename>", repoHandlerWrapper)

	server := fasthttp.Server{
		Handler:            router.HandleRequest,
		DisableKeepalive:   true,
		GetOnly:            true,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 10,
	}

	log.Println("serving on port", *port)
	e := server.ListenAndServe(fmt.Sprintf(":%v", *port))
	if e != nil {
		panic(e)
	}
}

func repoHandlerWrapper(ctx *routing.Context) error {
	owner := ctx.Param("owner")
	project := ctx.Param("project")
	version := ctx.Param("version")
	filename := ctx.Param("filename")

	if !strings.HasPrefix(filename, project+"-"+version) {
		return http.ErrNotSupported
	}

	switch filepath.Ext(filename) {
	case ".pom":
		return pom(ctx, owner, project, version)
	case ".md5":
		return md5(ctx, owner, project, version)
	case ".jar":
		return jar(ctx, owner, project, version)
	case ".sha1":
		return nil
	default:
		return http.ErrNotSupported
	}
}

func jar(ctx *routing.Context, owner, project, version string) error {
	p, e := ore.GetProject(owner, project)
	if e != nil {
		return e
	}
	r, e := ore.GetJar(p.ID, version)
	if e != nil {
		return e
	}
	defer r.Close()
	_, e = io.Copy(ctx.Response.BodyWriter(), r)
	return e
}

func md5(ctx *routing.Context, owner, project, version string) error {
	p, e := ore.GetProject(owner, project)
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

func pom(ctx *routing.Context, owner, project, version string) error {
	p, e := ore.GetProject(owner, project)
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
	return en.Encode(ore.NewPom(owner, project, v.Name))
}
