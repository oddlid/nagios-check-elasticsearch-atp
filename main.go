package main

/*
Rewrite of vgt_check_elasticsearch_atp.pl
Odd, 2016-07-05 14:36:19
*/


import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli" // renamed from codegansta
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	VERSION    string  = "2016-07-06"
	UA         string  = "VGT MnM ElasticCheck/1.0"
	DEF_INDEX  string  = "monitoring"
	DEF_QUERY  string  = "ATP"
	DEF_HPORT  uint16  = 3302
	DEF_EPORT  uint16  = 9200
	DEF_WARN   int     = 60
	DEF_CRIT   int     = 120
	DEF_TMOUT  float64 = 10.0
	DEF_PROT   string  = "http"
	URL_TMPL   string  = "%s://%s:%d/%s-*/_search?pretty=false&amp;sort=@timestamp:desc&amp;size=1&amp;q=%s&amp;limit=1"
	S_OK       string  = "OK"
	S_WARNING  string  = "WARNING"
	S_CRITICAL string  = "CRITICAL"
	S_UNKNOWN  string  = "UNKNOWN"
	E_OK       int     = 0
	E_WARNING  int     = 1
	E_CRITICAL int     = 2
	E_UNKNOWN  int     = 3
)

func Log2HAProxy(adr, data string, tmout time.Duration) error {
	c, err := net.DialTimeout("tcp", adr, tmout)
	if err != nil {
		return err
	}
	fmt.Fprint(c, data)
	return c.Close()
}

// geturl() fetches a URL and returns the HTTP response
func geturl(url string, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", UA)

	tr := &http.Transport{
		DisableKeepAlives: true, // we're not reusing the connection, so don't let it hang open
	}
	if strings.Index(url, "https") >= 0 {
		// Verifying certs is not the job of this plugin,
		// so we save ourselves a lot of grief by skipping any SSL verification
		// Could be a good idea for later to set this at runtime instead
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Transport: tr,
		Timeout: timeout,
	}

	return client.Do(req)
}

// run_check() takes the CLI params and glue together all logic in the program
func run_check(c *cli.Context) error {
	hah := c.String("ha-host")
	elh := c.String("el-host")
	idx := c.String("index")
	qry := c.String("query")
	hap := c.Uint("ha-port")
	elp := c.Uint("el-port")
	tmout := c.Float64("timeout")

	url := fmt.Sprintf(URL_TMPL, DEF_PROT, elh, elp, idx, qry)
	conn_timeout := time.Second * time.Duration(tmout)
	log.Debugf("Elasticsearch URL: %q", url)

	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err)
		hostname = "localhost"
	}
	ha_log_ts := time.Now().Unix()
	ha_log_entry := fmt.Sprintf("%d MONITORING ATP - ELK - running from %s\n", ha_log_ts, hostname)
	ha_adr := fmt.Sprintf("%s:%d", hah, hap)
	log.Debugf("HAProxy adr: %q", ha_adr)
	err = Log2HAProxy(ha_adr, ha_log_entry, conn_timeout)
	if err != nil {
		log.Error(err)
	}

	resp, err := geturl(url, conn_timeout)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", data)

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "check_elasticsearch-atp"
	app.Version = VERSION
	app.Author = "Odd E. Ebbesen"
	app.Email = "odd.ebbesen@wirelesscar.com"
	app.Usage = "Check Elasticsearch All-Through-Ping and alert in Nagios/op5"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ha-host",
			Usage: "HAProxy hostname or IP",
		},
		cli.UintFlag{
			Name:  "ha-port",
			Value: uint(DEF_HPORT),
			Usage: "HAProxy port",
		},
		cli.StringFlag{
			Name:  "el-host",
			Usage: "Elasticsearch hostname or IP",
		},
		cli.UintFlag{
			Name:  "el-port",
			Value: uint(DEF_EPORT),
			Usage: "Elasticsearch port",
		},
		cli.StringFlag{
			Name:  "index, i",
			Value: DEF_INDEX,
			Usage: "Elasticsearch index",
		},
		cli.StringFlag{
			Name:  "query, q",
			Value: DEF_QUERY,
			Usage: "Elasticsearch query",
		},
		cli.IntFlag{
			Name:  "warning, w",
			Value: DEF_WARN,
			Usage: "Warning threshold",
		},
		cli.IntFlag{
			Name:  "critical, c",
			Value: DEF_CRIT,
			Usage: "Critical threshold",
		},
		cli.Float64Flag{
			Name:  "timeout, t",
			Value: DEF_TMOUT,
			Usage: "Number of seconds before connection times out",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "fatal",
			Usage: "Log level (options: debug, info, warn, error, fatal, panic)",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "Run in debug mode",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stdout)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatal(err.Error())
		}
		log.SetLevel(level)
		if !c.IsSet("log-level") && !c.IsSet("l") && c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Action = run_check
	app.Run(os.Args)
}
