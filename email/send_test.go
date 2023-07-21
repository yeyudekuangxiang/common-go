package email

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	page := fmt.Sprintf("/pages/scene/charge/order/index?id=%s", "121212")
	println(page)

	orderTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2023-07-21 17:42:23", time.Local)
	if err != nil {

	}

	println(orderTime.Unix())

	a := time.Unix(orderTime.Unix(), 0)
	println(a.Unix())

	/*e := email.NewEmail()
	e.From = "dj <lvmiao@miotech.com>"
	e.To = []string{"18840853003@163.com"}
	e.Subject = "Awesome web"
	e.Text = []byte("Text Body is, of course, supported!")
	err := e.Send("smtp.126.com:25", smtp.PlainAuth("", "lvmiao@miotech.com", "lvmiao123", "smtp.126.com"))
	if err != nil {
		log.Fatal(err)
	}*/
}
