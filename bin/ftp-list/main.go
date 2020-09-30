package main

import (
	"flag"
	"fmt"
	"time"

	"bitbucket.org/advbet/ibos"
	"github.com/sirupsen/logrus"
)

func main() {
	var baseURL string

	flag.StringVar(&baseURL, "url", "ftp://user:pass@azftp.phumelela.com", "Base URL for FTP-pull delivery method")
	flag.Parse()

	started := time.Now()
	c, err := ibos.NewFTPClient(baseURL)
	if err != nil {
		logrus.WithError(err).Fatal("creating FTP-pull client")
	}

	files, err := c.List()
	if err != nil {
		logrus.WithError(err).Fatal("listing remote files")
	}

	for _, f := range files {
		fmt.Printf("\t%s\n", f)
	}
	fmt.Printf("Duration: %s\n", time.Since(started))
	fmt.Printf("Total: %d\n", len(files))
}
