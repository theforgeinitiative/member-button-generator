package dymo

import (
	"crypto/tls"
	_ "embed"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/theforgeinitiative/member-button-generator/sfdc"
)

//go:embed label.xml
var labelTemplate string

type Client struct {
	URL     string
	Printer string

	HTTPClient *http.Client
}

func NewClient() Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// This is gross, but there is some janky CA Dymo generates
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return Client{
		URL: "https://localhost:41951/DYMO/DLS/Printing",
		HTTPClient: &http.Client{
			Transport: transport,
		},
	}
}

func (c *Client) GetPrinters() ([]string, error) {
	resp, err := c.HTTPClient.Get(c.URL + "/GetPrinters")
	if err != nil {
		return nil, fmt.Errorf("failed to get printers: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received bad status code from api: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body from api: %w", err)
	}

	var printers Printers
	err = xml.Unmarshal(body, &printers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse api response: %w", err)
	}

	var printerList []string
	for _, p := range printers.Printers {
		printerList = append(printerList, p.Name)
	}
	return printerList, nil
}

func (c *Client) PrintLabels(buttons []sfdc.MemberButton) error {
	var ls LabelSet
	for _, b := range buttons {
		ls.LabelRecord = append(ls.LabelRecord,
			LabelRecord{
				ObjectData: []LabelField{
					{Name: "MEMBER", Text: b.Name},
					{Name: "BARCODE", Text: b.Barcode},
				},
			},
		)
	}
	labelSetXml, err := xml.Marshal(ls)
	if err != nil {
		return fmt.Errorf("failed to create labelset: %w", err)
	}

	data := url.Values{}
	data.Set("printerName", c.Printer)
	data.Set("printParamsXml", "")
	data.Set("labelXml", labelTemplate)
	data.Set("labelSetXml", string(labelSetXml))

	resp, err := c.HTTPClient.PostForm(c.URL+"/PrintLabel", data)
	if err != nil {
		return fmt.Errorf("failed to print label: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status code from api: %d", resp.StatusCode)
	}
	resp.Body.Close()

	return nil
}
