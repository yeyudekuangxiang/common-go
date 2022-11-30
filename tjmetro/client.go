package tjmetro

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	mrand "math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	ticketAllotPath = "/tj-metro-api/open-forward/api/eTicket/allot"
	publicKey       = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCMfK901Hst5gspb0doMp2+33P+BuZ+MHJLwRUje1+/Vx9K8BvPqbpIX6V120CJYchDnmHRDiVT8gq2jRMre0T20K9cOagJdKTKEAxNrTOKdx/CrtYxNSREGZ5yHLvN8+LAC7zorMdeb+FoPvDshPpSgnBYzQm7qr+kmM6GHA7eVQIDAQAB
-----END PUBLIC KEY-----`
)

type Client struct {
	//https://app.trtpazyz.com
	Domain    string
	AppId     string
	Version   string
	htpClient http.Client
}

// TicketAllot 发放电子票
func (c *Client) TicketAllot(param TicketAllotParam) (resp *BaseResponse, bizId string, err error) {
	return c.request(c.Domain+ticketAllotPath, param)
}
func (c *Client) request(url string, v interface{}) (resp *BaseResponse, bizId string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}
	s, err := c.sign(data)
	if err != nil {
		return nil, "", err
	}
	fmt.Println(string(data))
	fmt.Println(s)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, "", err
	}
	bizId = time.Now().Format("20060102150405") + c.rand()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appid", c.AppId)
	req.Header.Add("sequence", bizId)
	req.Header.Add("version", c.Version)
	req.Header.Add("signature", s)

	htpRes, err := c.htpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer htpRes.Body.Close()

	resp = &BaseResponse{}
	err = json.NewDecoder(htpRes.Body).Decode(resp)
	if err != nil {
		return nil, "", err
	}
	return resp, bizId, nil
}
func (c *Client) sign(v []byte) (string, error) {
	signData, err := c.rsa([]byte(c.md5(v)))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signData), nil
}
func (c *Client) md5(d []byte) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(d)))
}
func (c *Client) rsa(d []byte) ([]byte, error) {
	pk, err := c.publicKey()
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pk, d)
}
func (c *Client) publicKey() (*rsa.PublicKey, error) {
	b, _ := pem.Decode([]byte(publicKey))
	cert, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.(*rsa.PublicKey), nil
}
func (c *Client) rand() string {
	s := ""
	for i := 0; i < 10; i++ {
		s += strconv.Itoa(mrand.Intn(10))
	}
	return s
}
