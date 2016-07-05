package main

/*
Rewrite of vgt_check_elasticsearch_atp.pl
Odd, 2016-07-05 14:36:19
*/


import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli" // renamed from codegansta
	"os"
)

const (
	VERSION    string  = "2016-07-05"
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

// run_check() takes the CLI params and glue together all logic in the program
func run_check(c *cli.Context) {
	//hah := c.String("ha-hostname")
	//elh := c.String("el-hostname")
	//idx := c.String("index")
	//qry := c.String("query")
	//hap := c.Int("ha-port")
	//elp := c.Int("e-port")
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
			Name:  "ha-hostname",
			Usage: "HAProxy hostname or IP",
		},
		cli.IntFlag{
			Name:  "ha-port",
			Value: int(DEF_HPORT),
			Usage: "HAProxy port",
		},
		cli.StringFlag{
			Name:  "el-hostname",
			Usage: "Elasticsearch hostname or IP",
		},
		cli.IntFlag{
			Name:  "el-port",
			Value: int(DEF_EPORT),
			Usage: "Elasticsearch port",
		},
		cli.StringFlag{
			Name:  "index, i",
			Usage: "Elasticsearch index",
		},
		cli.StringFlag{
			Name:  "query, q",
			Usage: "Elasticsearch query",
		},
		cli.IntFlag{
			Name:  "warning, w",
			Usage: "Warning threshold",
		},
		cli.IntFlag{
			Name:  "critical, c",
			Usage: "Critical threshold",
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
