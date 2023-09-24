package parsers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net"
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
	FailureReportingOptions string `xml:"fo"`
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

	if err := report.Validate(); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *Report) Validate() error {
	if err := r.PolicyPublished.Validate(); err != nil {
		return err
	}

	if err := r.ReportMetadata.Validate(); err != nil {
		return err
	}

	for _, record := range r.Records {
		record.Validate()
	}

	return nil
}

func (r *PolicyPublished) Validate() error {
	// PolicyPublished is required
	if r == (&PolicyPublished{}) {
		return errors.New("policy published is required")
	}

	// Domain is required
	if r.Domain == "" {
		return errors.New("policy published - [domain] is required")
	}

	// AlignmentModeDKIM is optional
	if r.AlignmentModeDKIM != "" {
		// AlignmentModeDKIM must be one of these values if present
		if r.AlignmentModeDKIM != "r" && r.AlignmentModeDKIM != "s" {
			return errors.New("policy published - [alignment mode dkim] must be one of these values: [r, s], got: " + string(r.AlignmentModeDKIM))
		}
	}

	// AlignmentModeSPF is optional
	if r.AlignmentModeSPF != "" {
		// AlignmentModeSPF must be one of these values if present
		if r.AlignmentModeSPF != "r" && r.AlignmentModeSPF != "s" {
			return errors.New("policy published - [alignment mode spf] must be one of these values: [r, s], got: " + string(r.AlignmentModeSPF))
		}
	}

	// Policy is required
	if r.Policy == "" {
		return errors.New("policy published - [policy] is required")
	}

	// Policy must be one of these values
	if r.Policy != "none" &&
		r.Policy != "quarantine" &&
		r.Policy != "reject" {
		return errors.New("policy published - [policy] must be one of these values: [none, quarantine, reject], got: " + r.Policy)
	}

	// SubdomainPolicy is optional
	if r.SubdomainPolicy != "" {
		// SubdomainPolicy must be one of these values
		if r.SubdomainPolicy != "none" &&
			r.SubdomainPolicy != "quarantine" &&
			r.SubdomainPolicy != "reject" {
			return errors.New("policy published - [subdomain policy] must be one of these values: [none, quarantine, reject], got: " + r.SubdomainPolicy)
		}
	}

	// Percentage must be between 0 and 100
	if r.Percentage < 0 || r.Percentage > 100 {
		return errors.New("policy published - [percentage] must be between 0 and 100, got: " + fmt.Sprintf("%d", r.Percentage))
	}

	// FailureReportingOptions is optional
	if r.FailureReportingOptions != "" {
		// FailureReportingOptions must be one of these values
		if r.FailureReportingOptions != "0" &&
			r.FailureReportingOptions != "1" &&
			r.FailureReportingOptions != "d" &&
			r.FailureReportingOptions != "s" {
			return errors.New("policy published - [failure reporting options] must be one of these values: [0, 1, d, s], got: " + string(r.FailureReportingOptions))
		}
	}

	return nil
}

func (r *ReportMetadata) Validate() error {
	// ReportMetadata is required
	if r == (&ReportMetadata{}) {
		return errors.New("report metadata is required")
	}

	// OrgName is required
	if r.OrgName == "" {
		return errors.New("report metadata - [org name] is required")
	}

	// Email is required
	if r.Email == "" {
		return errors.New("report metadata - [email] is required")
	}

	// ReportID is required
	if r.ReportID == "" {
		return errors.New("report metadata - [report id] is required")
	}

	// DateRange is required
	if err := r.DateRange.Validate(); err != nil {
		return err
	}

	return nil
}

func (d *DateRange) Validate() error {
	// DateRange is required
	if d == (&DateRange{}) {
		return errors.New("date range is required")
	}

	// Begin is required
	if d.Begin == 0 {
		return errors.New("date range - [begin] is required")
	}

	// End is required
	if d.End == 0 {
		return errors.New("date range - [end] is required")
	}

	return nil
}

func (r *Record) Validate() error {
	if err := r.AuthResults.Validate(); err != nil {
		return err
	}

	if err := r.Identifiers.Validate(); err != nil {
		return err
	}

	if err := r.Row.Validate(); err != nil {
		return err
	}

	return nil
}

func (a *AuthResult) Validate() error {
	if err := a.DKIM.Validate(); err != nil {
		return err
	}

	if err := a.SPF.Validate(); err != nil {
		return err
	}

	return nil
}

func (spf *SPFAuthResult) Validate() error {
	// SPF is required
	if spf == (&SPFAuthResult{}) {
		return errors.New("spf is required")
	}

	// Domain is required
	if spf.Domain == "" {
		return errors.New("spf - [domain] is required")
	}

	// Result is required
	if spf.Result == "" {
		return errors.New("spf - [result] is required")
	}

	// Scope is optional
	if spf.Scope != "" {
		// Scope must be one of these two values if present
		if spf.Scope != "mfrom" && spf.Scope != "helo" {
			return errors.New("spf - [scope] is not empty, so it must be one of these values: [mfrom, helo], got: " + spf.Scope)
		}
	}

	// Result must be one of these values
	if spf.Result != "pass" &&
		spf.Result != "fail" &&
		spf.Result != "none" &&
		spf.Result != "neutral" &&
		spf.Result != "softfail" &&
		spf.Result != "permerror" &&
		spf.Result != "temperror" {
		return errors.New("spf - [result] must be one of these values: [pass, fail, none, neutral, softfail, permerror, temperror], got: " + spf.Result)
	}

	return nil
}

func (dkim *DKIMAuthResult) Validate() error {
	// DKIM is optional
	if dkim == (&DKIMAuthResult{}) {
		return nil
	}

	// If DKIM is not nil, then it must be valid
	// Domain is required
	if dkim.Domain == "" {
		return errors.New("dkim - [domain] is required")
	}

	// Selector is required
	if dkim.Selector == "" {
		return errors.New("dkim - [selector] is required")
	}

	// Result is required
	if dkim.Result == "" {
		return errors.New("dkim - [result] is required")
	}

	// Result must be one of these values
	if dkim.Result != "pass" &&
		dkim.Result != "fail" &&
		dkim.Result != "none" &&
		dkim.Result != "neutral" &&
		dkim.Result != "policy" &&
		dkim.Result != "permerror" &&
		dkim.Result != "temperror" {
		return errors.New("dkim - [result] must be one of these values: [pass, fail, none, neutral, policy, permerror, temperror], got: " + dkim.Result)
	}

	return nil
}

func (i *Identifiers) Validate() error {
	// Identifiers is required
	if i == (&Identifiers{}) {
		return errors.New("identifiers is required")
	}

	// EnvelopeFrom is required
	if i.EnvelopeFrom == "" {
		return errors.New("identifiers - [envelope from] is required")
	}

	// HeaderFrom is required
	if i.HeaderFrom == "" {
		return errors.New("identifiers - [header from] is required")
	}

	return nil
}

func (r *Row) Validate() error {
	// Row is required
	if r == (&Row{}) {
		return errors.New("row is required")
	}

	// SourceIP is required
	if r.SourceIP == "" {
		return errors.New("row - [source ip] is required")
	}

	if net.ParseIP(r.SourceIP) == nil {
		return errors.New("row - [source ip] is not a valid IP address")
	}

	// Count is required
	if r.Count == 0 {
		return errors.New("row - [count] cannot be 0")
	}

	if err := r.PolicyEvaluated.Validate(); err != nil {
		return err
	}

	return nil
}

func (p *PolicyEvaluated) Validate() error {
	// PolicyEvaluated is required
	if p == (&PolicyEvaluated{}) {
		return errors.New("policy evaluated is required")
	}

	// Disposition is required
	if p.Disposition == "" {
		return errors.New("policy evaluated - [disposition] is required")
	}

	// Disposition must be one of these values
	if p.Disposition != "none" &&
		p.Disposition != "reject" &&
		p.Disposition != "quarantine" {
		return errors.New("policy evaluated - [disposition] must be one of these values: [none, quarantine, reject], got: " + p.Disposition)
	}

	// DKIM is optional
	if p.DKIM != "" {
		// DKIM must be one of these values
		if p.DKIM != "pass" &&
			p.DKIM != "fail" {
			return errors.New("policy evaluated - [dkim] must be one of these values: [pass, fail], got: " + p.DKIM)
		}
	}

	// SPF is optional
	if p.SPF != "" {
		if p.SPF != "pass" &&
			p.SPF != "fail" {
			return errors.New("policy evaluated - [spf] must be one of these values: [pass, fail], got: " + p.SPF)
		}
	}

	return nil
}
