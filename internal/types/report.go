package types

type Report struct {
	Version         string          `json:"version" xml:"version"`
	ReportMetadata  ReportMetadata  `json:"source" xml:"report_metadata"`
	PolicyPublished PolicyPublished `json:"policy_published" xml:"policy_published"`
	Records         []Record        `json:"records" xml:"record"`
}

type ReportMetadata struct {
	OrgName          string    `json:"org_name" xml:"org_name"`
	Email            string    `json:"email" xml:"email"`
	ExtraContactInfo string    `json:"extra_contact_info" xml:"extra_contact_info"`
	ReportID         string    `json:"report_id" xml:"report_id"`
	DateRange        DateRange `json:"date_range" xml:"date_range"`
}

type DateRange struct {
	Begin int64 `json:"begin" xml:"begin"`
	End   int64 `json:"end" xml:"end"`
}

type PolicyPublished struct {
	Domain                  string `json:"domain" xml:"domain"`
	AlignmentModeDKIM       string `json:"adkim" xml:"adkim"`
	AlignmentModeSPF        string `json:"aspf" xml:"aspf"`
	Policy                  string `json:"p" xml:"p"`
	SubdomainPolicy         string `json:"sp" xml:"sp"`
	Percentage              int    `json:"pct" xml:"pct"`
	FailureReportingOptions rune   `json:"fo" xml:"fo"`
}

type Record struct {
	Row         Row         `json:"row" xml:"row"`
	Identifiers Identifiers `json:"identifiers" xml:"identifiers"`
	AuthResults AuthResult  `json:"auth_results" xml:"auth_results"`
}

type Row struct {
	SourceIP        string          `json:"source_ip" xml:"source_ip"`
	Count           int             `json:"count" xml:"count"`
	PolicyEvaluated PolicyEvaluated `json:"policy" xml:"policy_evaluated"`
}

type PolicyEvaluated struct {
	Disposition string `json:"disposition" xml:"disposition"`
	DKIM        string `json:"dkim" xml:"dkim"`
	SPF         string `json:"spf" xml:"spf"`
}

type Identifiers struct {
	EnvelopeTo   string `json:"envelope_to" xml:"envelope_to"`
	EnvelopeFrom string `json:"envelope_from" xml:"envelope_from"`
	HeaderFrom   string `json:"header_from" xml:"header_from"`
}

type AuthResult struct {
	DKIM DKIMAuthResult `json:"dkim" xml:"dkim"`
	SPF  SPFAuthResult  `json:"spf" xml:"spf"`
}

type DKIMAuthResult struct {
	Domain      string `json:"domain" xml:"domain"`
	Selector    string `json:"selector" xml:"selector"`
	Result      string `json:"result" xml:"result"`
	HumanResult string `json:"human_result" xml:"human_result"`
}

type SPFAuthResult struct {
	Domain      string `json:"domain" xml:"domain"`
	Scope       string `json:"scope" xml:"scope"`
	Result      string `json:"result" xml:"result"`
	HumanResult string `json:"human_result" xml:"human_result"`
}
