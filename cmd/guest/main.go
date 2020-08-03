// !build +windows

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"github.com/urfave/cli/v2"

	"github.com/ttys3/uap"
)

var Version = "dev"
var CommitSHA = "dev"
var BuildDate = "unkown"

const appName = "uap"

var machine = "unkown"

var netClient *http.Client

func init() {
	m, err := os.Hostname()
	if err == nil {
		machine = m
	}
	netClient = &http.Client{
		Timeout: time.Second * 2,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 1 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 1 * time.Second,
		},
	}
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Value:   "",
				Usage:   "host API url",
				EnvVars: []string{"UAP_HOST"},
			},
			&cli.StringFlag{
				Name:    "auth",
				Aliases: []string{"p"},
				Value:   uap.DftPasswd,
				Usage:   "auth password",
				EnvVars: []string{"UAP_AUTH"},
			},
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Value:   uap.RepoURL,
				Usage:   "url to open",
				EnvVars: []string{"UAP_URL"},
			},
		},
		Action: sendUrl,
	}

	err := app.Run(os.Args)
	fmt.Println(err)
	if err != nil {
		walk.MsgBox(nil, fmt.Sprintf("%s - %s - %s %s %s", "open URL failed",
			appName, Version, CommitSHA, BuildDate), err.Error(),
			walk.MsgBoxIconError)
	}
}

func sendUrl(c *cli.Context) error {
	hostURL := c.String("host")
	if hostURL == "" {
		return fmt.Errorf("you need to set the UAP_HOST env")
	}
	urlToOpen := c.Args().First()
	urlParam := c.String("url")
	if urlToOpen == "" && urlParam == "" {
		return fmt.Errorf("url can not be empty")
	}
	if urlToOpen == "" {
		urlToOpen = urlParam
	}

	auth := c.String("auth")
	reqBody, err := json.Marshal(uap.UaProxyReq{
		Auth:        auth,
		FromMachine: machine,
		Url:         urlToOpen,
		ReqTs:       time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	apiURl := fmt.Sprintf("%s/open", hostURL)
	resp, err := netClient.Post(apiURl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	var rsp uap.UaProxyRsp
	if err := json.Unmarshal(body, &rsp); err != nil {
		return err
	}
	if rsp.RetCode != uap.RetCodeOK {
		return fmt.Errorf("%s", rsp.Msg)
	}
	return nil
}
