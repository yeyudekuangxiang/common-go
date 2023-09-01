package email

import (
	"bytes"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/idtool"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type mail struct {
	FromUser string
	Passwd   string
}

//初始化用户名和密码

func NewEmailClient(fromUser string, passwd string) mail {
	temp := mail{
		FromUser: fromUser,
		Passwd:   passwd,
	}
	return temp
}

type SendInvoiceParam struct {
	ToUser      string
	Subject     string
	ApplyDate   string
	InvoiceDate string
	Title       string
	Price       string
	Annex       []string
}

func (m mail) SendInvoice(param SendInvoiceParam) error {
	e := email.NewEmail()
	e.From = m.FromUser
	e.To = []string{param.ToUser}
	e.Subject = param.Subject
	/*// 解析html模板
	t1, err := template.ParseFiles("email-template.html")
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = t1.Execute(body, struct {
		ApplyDate   string
		InvoiceDate string
		Title       string
		Price       string
	}{
		ApplyDate:   param.ApplyDate,
		InvoiceDate: param.InvoiceDate,
		Title:       param.Title,
		Price:       param.Price,
	})
	if err != nil {
		return err
	}*/
	body3 := `<!DOCTYPE html><html lang='en'xmlns='http://www.w3.org/1999/html'><head><meta charset='UTF-8'><title>Title</title></head><body><div><table style='min-width: 348px; transform: scale(0.784483, 0.784483); transform-origin: left top 0px'><tbody><tr><td>尊敬的客户您好：</td></tr><tr><td>您于` + param.ApplyDate + `开具电子发票，我们将电子发票发送给您，以便作为您的维权保修凭证、报销凭证。</td></tr><tr><td>&nbsp;</td></tr><tr><td>发票信息如下：</td></tr><tr><td>开票日期：` + param.InvoiceDate + `</td></tr><tr><td>购方名称：` + param.Title + `</td></tr><tr><td>价税合计：` + param.Price + `</td></tr><tr><td>&nbsp;</td></tr><tr><td>详情请见附件</td></tr></tbody></table></div></body></html>`
	e.HTML = []byte(body3)
	for _, annex := range param.Annex {
		response, e1 := http.Get(annex)
		if e1 != nil {
			log.Fatal(e1)
		}
		defer response.Body.Close()
		fileName := idtool.UUID()
		file, err := os.Create(fileName + ".pdf")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}
		_, err = e.AttachFile(fileName + ".pdf")
		if err != nil {
			continue
		}
		err = os.Remove(fileName + ".pdf")
		if err != nil {
			log.Fatal(err)
		}
	}
	err := e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", m.FromUser, m.Passwd, "smtp.gmail.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})
	if err != nil {
		return err
	}
	return nil
}

func (m mail) SendInvoiceV2(param SendInvoiceParam) error {
	e := email.NewEmail()
	e.From = m.FromUser
	e.To = []string{param.ToUser}
	e.Subject = param.Subject
	// 解析html模板
	t1, err := template.ParseFiles("https://miotech-mio.oss-cn-shanghai.aliyuncs.com/static/mp2c/test/email-template.html")
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = t1.Execute(body, struct {
		ApplyDate   string
		InvoiceDate string
		Title       string
		Price       string
	}{
		ApplyDate:   param.ApplyDate,
		InvoiceDate: param.InvoiceDate,
		Title:       param.Title,
		Price:       param.Price,
	})
	if err != nil {
		return err
	}

	e.HTML = body.Bytes()
	for _, annex := range param.Annex {
		_, err = e.AttachFile(annex)
		if err != nil {
			continue
		}
	}
	err = e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", m.FromUser, m.Passwd, "smtp.gmail.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})
	if err != nil {
		return err
	}
	return nil
}
