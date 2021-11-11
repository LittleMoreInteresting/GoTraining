package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/**
1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
*/
const addr = ":8808"

func httpServerStart(ctx context.Context, eg *errgroup.Group, svr *http.Server) {
	eg.Go(func() error {
		<-ctx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		fmt.Println("Done received", svr.Addr)
		err := svr.Shutdown(shutdownCtx)
		fmt.Println("Shutdown success", svr.Addr)
		return err
	})
	eg.Go(func() error {
		fmt.Println("ListenAndServe", svr.Addr)
		return svr.ListenAndServe()
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	group, errCtx := errgroup.WithContext(ctx)

	svr1 := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	svr2 := &http.Server{
		Addr: "127.0.0.1:8081",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "Go01: Hello world\n")
	})
	httpServerStart(errCtx, group, svr1)
	httpServerStart(errCtx, group, svr2)

	group.Go(func() error {

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT)
		fmt.Println("[G02]:start...")
		select {
		case <-errCtx.Done():
			fmt.Println("[G02]: 异常...")
			return errCtx.Err()
		case s := <-c:
			fmt.Printf("[G02]退出%v...\n", s)
			cancel()
			return errors.New("[G02] :退出")
		}
	})

	err := group.Wait()
	if err == nil {
		fmt.Println("都完成了")
	} else {
		fmt.Printf("get error:%v", err)
	}

}
