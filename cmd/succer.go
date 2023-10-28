package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"encoding/json"

	"golang.org/x/net/proxy"
)

type Succer struct {
	domain       string
	output       string
	outputJson   bool
	socksProxy   string
	FoundDomains []string
}

type Envelope struct {
	Body struct {
		GetFederationInformationResponseMessage struct {
			Response struct {
				Domains struct {
					Domain []string `xml:"Domain"`
				} `xml:"Domains"`
			} `xml:"Response"`
		} `xml:"GetFederationInformationResponseMessage"`
	} `xml:"Body"`
}

func (s *Succer) Run() error {

	// Get the tenant domains
	s.enumerateTenantDomains("")

	if s.outputJson {
		// Construct the json data
		jsonData := make(map[string]interface{})
		jsonData["domains"] = s.FoundDomains
		jsonData["source"] = s.domain

		jsonOutput, err := json.Marshal(jsonData)
		if err != nil {
			log.Fatalf("Failed to marshal data: %v", err)
		}
		if s.output != "" {
			// Write to file
			file, err := os.Create(s.output)
			if err != nil {
				log.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()
			_, err = file.Write(jsonOutput)
			if err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
		} else {
			// Write to stdout
			fmt.Println(string(jsonOutput))
		}
	} else {
		if s.output != "" {
			// Write to file
			file, err := os.Create(s.output)
			if err != nil {
				log.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()
			for _, domain := range s.FoundDomains {
				file.WriteString(domain + "\n")
			}
		} else {
			// Write to stdout
			for _, domain := range s.FoundDomains {
				fmt.Println(domain)
			}
		}
	}
	return nil
}

func (s *Succer) constructXML() string {
	endpoint := "https://autodiscover-s.outlook.com/autodiscover/autodiscover.svc"
	xmlTemplate := `
	<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:exm="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:ext="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
		<soap:Header>
			<a:Action soap:mustUnderstand="1">http://schemas.microsoft.com/exchange/2010/Autodiscover/Autodiscover/GetFederationInformation</a:Action>
			<a:To soap:mustUnderstand="1">%s</a:To>
			<a:ReplyTo>
				<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
			</a:ReplyTo>
		</soap:Header>
		<soap:Body>
			<GetFederationInformationRequestMessage xmlns="http://schemas.microsoft.com/exchange/2010/Autodiscover">
				<Request>
					<Domain>%s</Domain>
				</Request>
			</GetFederationInformationRequestMessage>
		</soap:Body>
	</soap:Envelope>
	`
	return fmt.Sprintf(xmlTemplate, endpoint, s.domain)

}

func (s *Succer) enumerateTenantDomains(userAgent string) []string {

	var httpClient *http.Client

	if userAgent == "" {
		userAgent = "AutodiscoverClient"
	}

	headers := map[string]string{
		"Content-Type": "text/xml; charset=utf-8",
		"SOAPAction":   `"http://schemas.microsoft.com/exchange/2010/Autodiscover/Autodiscover/GetFederationInformation"`,
		"User-Agent":   userAgent, // This should be set to the provided userAgent
	}

	xmlBody := s.constructXML()
	xmlData := []byte(strings.TrimSpace(xmlBody))
	endpoint := "https://autodiscover-s.outlook.com/autodiscover/autodiscover.svc"

	if s.socksProxy != "" {
		dialer, err := proxy.SOCKS5("tcp", s.socksProxy, nil, proxy.Direct)
		if err != nil {
			log.Fatalf("Failed to create proxy dialer: %v", err)
			return nil // Added a return to gracefully exit the function
		}
		httpTransport := &http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		httpClient = &http.Client{Transport: httpTransport}
	} else {
		httpClient = &http.Client{}
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(xmlData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		return nil // Added a return to gracefully exit the function
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
		return nil // Added a return to gracefully exit the function
	}
	defer resp.Body.Close()

	if resp.StatusCode == 421 {
		log.Println("No tenant domains found.")
		return nil
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return nil
	}

	var envelope Envelope
	err = xml.Unmarshal(bodyBytes, &envelope)
	if err != nil {
		log.Printf("Error decoding XML: %v", err)
		return nil // Added a return to gracefully exit the function
	}

	domainList := envelope.Body.GetFederationInformationResponseMessage.Response.Domains.Domain

	// Storing the list of domains
	// Remove the onmicrosoft.com domains and lowercase contents
	for _, domain := range domainList {
		if !strings.Contains(domain, "onmicrosoft") {
			// lowercase the domain
			s.FoundDomains = append(s.FoundDomains, strings.ToLower(domain))

		}
	}

	return nil

}
