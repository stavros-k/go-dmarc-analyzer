package parsers

import (
	"encoding/xml"
)

type Report struct {
	Version         string          `xml:"version"`
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Records         []Record        `xml:"record"`
}

type ReportMetadata struct {
	OrgName          string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ExtraContactInfo string    `xml:"extra_contact_info"`
	ReportID         string    `xml:"report_id"`
	DateRange        DateRange `xml:"date_range"`
}

type DateRange struct {
	Begin int64 `xml:"begin"`
	End   int64 `xml:"end"`
}

type PolicyPublished struct {
	Domain                  string `xml:"domain"`
	AlignmentModeDKIM       string `xml:"adkim"`
	AlignmentModeSPF        string `xml:"aspf"`
	Policy                  string `xml:"p"`
	SubdomainPolicy         string `xml:"sp"`
	Percentage              int    `xml:"pct"`
	FailureReportingOptions rune   `xml:"fo"`
}

type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResult  `xml:"auth_results"`
}

type Row struct {
	SourceIP        string          `xml:"source_ip"`
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

type Identifiers struct {
	EnvelopeTo   string `xml:"envelope_to"`
	EnvelopeFrom string `xml:"envelope_from"`
	HeaderFrom   string `xml:"header_from"`
}

type AuthResult struct {
	DKIM DKIMAuthResult `xml:"dkim"`
	SPF  SPFAuthResult  `xml:"spf"`
}

type DKIMAuthResult struct {
	Domain      string `xml:"domain"`
	Selector    string `xml:"selector"`
	Result      string `xml:"result"`
	HumanResult string `xml:"human_result"`
}

type SPFAuthResult struct {
	Domain      string `xml:"domain"`
	Scope       string `xml:"scope"`
	Result      string `xml:"result"`
	HumanResult string `xml:"human_result"`
}

// NewReport creates a new Report from a byte slice (opened file)
func NewReport(b []byte) (*Report, error) {
	report := &Report{}
	if err := xml.Unmarshal(b, &report); err != nil {
		return nil, err
	}

	return report, nil
}
