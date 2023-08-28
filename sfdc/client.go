package sfdc

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/simpleforce/simpleforce"
)

// Eventually make this configurable/dynamic
const MembershipYearStart = "2023-07-01"

const authSessionLength = 1 * time.Hour

type Client struct {
	SFClient          *simpleforce.Client
	clientSecret      string
	lastAuthenticated time.Time
}

func NewClient(url, clientID, clientSecret string) (Client, error) {
	sfc := simpleforce.NewClient(url, clientID, simpleforce.DefaultAPIVersion)
	err := sfc.LoginClientCredentials(clientSecret)
	if err != nil {
		return Client{}, fmt.Errorf("error making salesforce client: %w", err)
	}
	c := Client{
		SFClient:          sfc,
		clientSecret:      clientSecret,
		lastAuthenticated: time.Now(),
	}
	return c, nil
}

func (c *Client) Authenticate() error {
	err := c.SFClient.LoginClientCredentials(c.clientSecret)
	if err == nil {
		c.lastAuthenticated = time.Now()
	}
	return err
}

type MemberButton struct {
	ContactID   string `json:"id"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	LastPrinted string `json:"last_printed"`
}

func (c *Client) UnprintedButtons() ([]MemberButton, error) {
	if c.lastAuthenticated.Add(authSessionLength).Before(time.Now()) {
		c.Authenticate()
	}
	q := fmt.Sprintf(`
    SELECT
        Id,
        TFI_Display_Name_for_Button__c,
        TFI_Barcode_for_Button__c,
        Button_completed__c 
    FROM
        Contact 
    WHERE
		npo02__MembershipEndDate__c > %s
        AND Waivers_signed_date__c >= %s  
        AND ( Button_completed__c = null OR Button_completed__c < %s )  
        AND ( NOT Name LIKE '%%test%%' )
	`, MembershipYearStart, MembershipYearStart, MembershipYearStart)
	result, err := c.SFClient.Query(q)
	if err != nil {
		return nil, fmt.Errorf("error running SOQL query: %s", err)
	}

	var buttons []MemberButton
	for _, record := range result.Records {
		buttons = append(buttons, MemberButton{
			ContactID:   record.StringField("Id"),
			Name:        record.StringField("TFI_Display_Name_for_Button__c"),
			Barcode:     record.StringField("TFI_Barcode_for_Button__c"),
			LastPrinted: record.StringField("Button_completed__c"),
		})
	}
	return buttons, nil
}

func (c *Client) ButtonsByBarcode(barcode ...string) ([]MemberButton, error) {
	if c.lastAuthenticated.Add(authSessionLength).Before(time.Now()) {
		c.Authenticate()
	}
	var in string
	for i, b := range barcode {
		// test barcode is numeric
		_, err := strconv.Atoi(b)
		if err != nil {
			return nil, fmt.Errorf("barcode must be numeric: %w", err)
		}
		if i > 0 {
			in += ", "
		}
		in += "'" + b + "'"
	}
	q := fmt.Sprintf(`SELECT Id, TFI_Display_Name_for_Button__c, TFI_Barcode_for_Button__c, Button_completed__c
	FROM Contact
	WHERE TFI_Barcode_for_Button__c IN (%s)`, in)

	result, err := c.SFClient.Query(q)
	if err != nil {
		return nil, fmt.Errorf("error running SOQL query: %s", err)
	}

	var buttons []MemberButton
	for _, record := range result.Records {
		buttons = append(buttons, MemberButton{
			ContactID:   record.StringField("Id"),
			Name:        record.StringField("TFI_Display_Name_for_Button__c"),
			Barcode:     record.StringField("TFI_Barcode_for_Button__c"),
			LastPrinted: record.StringField("Button_completed__c"),
		})
	}
	return buttons, nil
}

func (c *Client) CompleteButton(id string) error {
	if c.lastAuthenticated.Add(authSessionLength).Before(time.Now()) {
		c.Authenticate()
	}
	updateObj := c.SFClient.SObject("Contact").
		Set("Id", id).
		Set("Button_completed__c", time.Now().Format("2006-01-02")).
		Update()

	if updateObj == nil {
		return errors.New("failed to update contact")
	}

	return nil
}
