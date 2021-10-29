package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)
func startHealthcheck(config Config, ctx context.Context) {
	log.Println("Starting health checker")

	if len(config.Healthcheck.Hosts) == 0{
		log.Printf("Nothing to check. Stopping health checker")
		return
	}

	doAllHealthChecks(config)
	for {
		select {
		case <-time.After(60 * time.Second):
			doAllHealthChecks(config)
		case <-ctx.Done():
			return
		}
	}
}

func doAllHealthChecks(config Config) {
	wg := sync.WaitGroup{}
	for _, host := range config.Healthcheck.Hosts {
		wg.Add(1)
		go func(endpoint string){
			defer wg.Done()
			err := doHealtcheck(endpoint, config.Healthcheck.Debug)
			if err != nil {
				telegramMessage := fmt.Sprintf("Health check for %s is failed.\n %s", host, err)
				log.Println(telegramMessage)
				sendTelegramMessage(telegramMessage)
			}
		}(host)
	}

	// wait for all checks to be completed
	wg.Wait()
}

func doHealtcheck(endpoint string, debug bool) error {
	if debug {
		log.Printf("Sending healthcheck to %s. \n", endpoint)
	}

	// we force health check to be under 10 seconds
	// otherwise we treat it as fail
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()

		var bodyString string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			bodyString = fmt.Sprintf("Failed to load body %s", err)
		} else {
			bodyString = string(bodyBytes)
		}

		if debug {
			log.Printf("Got non success response %v \n", resp)
		}

		return errors.New(fmt.Sprintf("StatusCode: %d. Response: %s", resp.StatusCode, bodyString))
	}

	// if endpoint is https we also validate expiration time
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0{
		for _,cert := range resp.TLS.PeerCertificates {

			// if certificate expires in less than 15 days, then we trigger fail
			maxExpirationTime := time.Now().Add(15 * time.Hour * 24)
			if cert.NotAfter.Before(maxExpirationTime) {
				return errors.New(fmt.Sprintf("Certificate %s expires at %s", cert.Subject.CommonName, cert.NotAfter.Format(time.RFC822Z)))
			}
		}
	}

	return nil
}


