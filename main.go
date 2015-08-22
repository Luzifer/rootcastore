package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Luzifer/rootcastore/cert"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cenkalti/backoff"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

var (
	config       = loadConfig()
	version      = "dev"
	s3Connection *s3.S3
)

const (
	certDataSource = "https://hg.mozilla.org/mozilla-central/raw-file/tip/security/nss/lib/ckfw/builtins/certdata.txt"
)

func init() {
	if config.EnableCache {
		s3Connection = s3.New(&aws.Config{})
	}
}

func main() {
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %dh", config.RefreshInterval), refreshFromSource)
	c.Start()

	refreshFromSource()

	r := mux.NewRouter()
	r.HandleFunc("/v1/store/{version}", handleGetStore)

	http.ListenAndServe(":3000", r)
}

func refreshFromSource() {
	certData := bytes.NewBuffer([]byte{})

	err := backoff.Retry(func() error {
		src, err := http.Get(certDataSource)
		if err != nil {
			return err
		}
		defer src.Body.Close()

		license, cvsID, objects := cert.ParseInput(src.Body)

		fmt.Fprint(certData, license)
		if len(cvsID) > 0 {
			fmt.Fprint(certData, "CVS_ID "+cvsID+"\n")
		}

		cert.OutputTrustedCerts(certData, objects)

		return nil
	}, backoff.NewExponentialBackOff())

	if err != nil {
		log.Fatal(err)
	}

	saveToCache(strconv.FormatInt(time.Now().UTC().Unix(), 10), certData.Bytes())
}

func handleGetStore(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	version = vars["version"]
	if version == "latest" {
		latest := getLatestVersion()
		if latest == "" {
			http.Error(res, fmt.Sprintf("Did not find version '%s'", version), http.StatusNotFound)
			return
		}

		http.Redirect(res, r, fmt.Sprintf("/v1/store/%s", latest), http.StatusFound)
    return
	}

	content, ok := getFromCache(version)

	if !ok {
		http.Error(res, fmt.Sprintf("Did not find version '%s'", version), http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/x-pem-file")
	res.Write(content)
}
