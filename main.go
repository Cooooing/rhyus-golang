package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
	"rhyus-golang/api"
	"rhyus-golang/common"
	"rhyus-golang/conf"
	"rhyus-golang/middlewares"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(conf.Conf.PprofPort), mux)
		if err != nil {
			common.Log.Fatal(common.ExitCodeUnavailablePort, "serve start failed: %s", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.New()
	ginServer.UseH2C = true
	ginServer.MaxMultipartMemory = 1024 * 1024 * 32 // 表示处理上传的文件时，最多将32MB的数据保存在内存中，超出部分会保存到临时文件中。这样可以避免大文件上传时占用过多内存。
	ginServer.Use(
		middlewares.Logging,
		middlewares.Recover,
		middlewares.CorsMiddleware,
		middlewares.Authorize,
	)

	api.ServeAPI(ginServer)

	addr := conf.Conf.Host + ":" + strconv.Itoa(conf.Conf.Port)
	common.Log.Info("server start at: %s", addr)
	if conf.Conf.Ssl.Enabled {
		err := ginServer.RunTLS(addr, conf.Conf.Ssl.CertPath, conf.Conf.Ssl.KeyPath)
		if err != nil {
			common.Log.Fatal(common.ExitCodeCertificateErr, "serve start failed: %s", err)
		}
	} else {
		err := ginServer.Run(addr)
		if err != nil {
			common.Log.Fatal(common.ExitCodeUnavailablePort, "serve start failed: %s", err)
		}
	}

}
