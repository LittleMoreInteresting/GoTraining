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
)

/**
1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
*/
const addr = ":8808"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	group, errCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		srv := &http.Server{Addr: addr}
		fmt.Println("[G01]:start...")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, "Go01: Hello world\n")
		})
		http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
			_ = srv.Shutdown(context.Background())
		})
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				fmt.Printf("[G01]异常:%v...\n", err)
				cancel()
				return
			}
		}()
		select {
		case <-errCtx.Done():
			fmt.Println("[G01]: 退出...")
			return errCtx.Err()
		}
	})
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
