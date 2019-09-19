/*
@Time : 2019/9/19 18:06
@Author : mp
@File : daemon
@Software: GoLand
*/
package Daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func startDaemon() {
	if syscall.Getppid() == 1 { // already daemon
		f, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if err != nil {
			fmt.Println("open /dev/null failed")
			os.Exit(-1)
		}

		fd := f.Fd()
		syscall.Dup2(int(fd), int(os.Stdin.Fd()))
		syscall.Dup2(int(fd), int(os.Stdout.Fd()))
		syscall.Dup2(int(fd), int(os.Stderr.Fd()))

		return
	}

	// 监听系统信号
	go func() {
		_c := make(chan os.Signal, 1)
		signal.Notify(_c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTSTP)
		msg := <- _c
		log.Println(msg)
		os.Exit(0)
	}()

	args := append([]string{os.Args[0]}, os.Args[1:]...)
	_, err := os.StartProcess(os.Args[0], args,
		&os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})

	if err != nil {
		log.Println("Daemon start failed:", err)
	}

	os.Exit(0)
}
