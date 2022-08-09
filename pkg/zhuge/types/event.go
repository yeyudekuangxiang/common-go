package types

type EventIOS struct {
	Ct   int64  `json:"$ct"`
	Eid  string `json:"$eid"`
	Cuid string `json:"$cuid"`
	Sid  int64  `json:"$sid"`
	Vn   string `json:"$vn"`
	Cn   string `json:"$cn"`
	Cr   int    `json:"$cr"`
	Os   string `json:"$os"`
	Ov   int    `json:"$ov"`
	Net  int    `json:"$net"`
}

func (EventIOS) F3127563263F0A80C8007E338109E07F() {}
func (e EventIOS) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"$ct":   e.Ct,
		"$eid":  e.Eid,
		"$cuid": e.Cuid,
		"$sid":  e.Sid,
		"$vn":   e.Vn,
		"$cn":   e.Cn,
		"$cr":   e.Cr,
		"$os":   e.Os,
		"$ov":   e.Ov,
		"$net":  e.Net,
	}
}

type EventAnd struct {
	Ct   int64  `json:"$ct"`
	Eid  string `json:"$eid"`
	Cuid string `json:"$cuid"`
	Sid  int64  `json:"$sid"`
	Vn   string `json:"$vn"`
	Cn   string `json:"$cn"`
	Cr   int    `json:"$cr"`
	Os   string `json:"$os"`
	Ov   int    `json:"$ov"`
	Net  int    `json:"$net"`
}

func (EventAnd) F3127563263F0A80C8007E338109E07F() {}
func (e EventAnd) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"$ct":   e.Ct,
		"$eid":  e.Eid,
		"$cuid": e.Cuid,
		"$sid":  e.Sid,
		"$vn":   e.Vn,
		"$cn":   e.Cn,
		"$cr":   e.Cr,
		"$os":   e.Os,
		"$ov":   e.Ov,
		"$net":  e.Net,
	}
}

type EventJs struct {
	Ct             int64  `json:"$ct"`
	Eid            string `json:"$eid"`
	Cuid           string `json:"$cuid"`
	Sid            int64  `json:"$sid"`
	ReferrerDomain string `json:"$referrer_domain"`
	Url            string `json:"$url"`
	Ref            string `json:"$ref"`
	UtmSource      string `json:"$utm_source"`
	UtmMedium      string `json:"$utm_medium"`
	UtmCampaign    string `json:"$utm_campaign"`
	UtmContent     string `json:"$utm_content"`
	UtmTerm        string `json:"$utm_term"`
}

func (EventJs) F3127563263F0A80C8007E338109E07F() {

}
func (e EventJs) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"$ct":              e.Ct,
		"$eid":             e.Eid,
		"$cuid":            e.Cuid,
		"$sid":             e.Sid,
		"$referrer_domain": e.ReferrerDomain,
		"$url":             e.Url,
		"$ref":             e.Ref,
		"$utm_source":      e.UtmSource,
		"$utm_medium":      e.UtmMedium,
		"$utm_campaign":    e.UtmCampaign,
		"$utm_content":     e.UtmContent,
		"$utm_term":        e.UtmTerm,
	}
}

type Event struct {
	Dt    string   `json:"dt"`
	Pl    string   `json:"pl"`
	Debug int      `json:"debug"`
	Ip    string   `json:"ip"`
	Pr    Attr     `json:"pr"`
	Usr   EventUsr `json:"usr"`
}

type EventUsr struct {
	Did string `json:"did"`
}
type EventWithAk struct {
	Ak string `json:"ak"`
	Event
}
