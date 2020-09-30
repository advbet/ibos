package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"bitbucket.org/advbet/ibos"

	"github.com/sirupsen/logrus"
)

func main() {
	var baseURL string
	var lastFilename string
	var interval time.Duration

	flag.StringVar(&baseURL, "url", "ftp://user:pass@azftp.phumelela.com", "Base URL for FTP-pull delivery method")
	flag.StringVar(&lastFilename, "last", "", "Name of last dowloaded document")
	flag.DurationVar(&interval, "interval", time.Minute, "Recheck interval for detecting new documents")
	flag.Parse()

	c, err := ibos.NewFTPClient(baseURL)
	if err != nil {
		logrus.WithError(err).Fatal("creating FTP-pull client")
	}

	stream := c.Stream(context.Background(), lastFilename, interval)
	for msg := range stream {
		if msg.Error != nil {
			logrus.WithError(msg.Error).Error("stream error")
			continue
		}
		fmt.Printf("==== %s ====\n", msg.Filename)
		fmt.Println(string(msg.Data))
	}
}
