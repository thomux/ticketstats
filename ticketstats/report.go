package ticketstats

// Report groups all evaluation results.
// This is the input for rendering the HTML report.
type Report struct {
	Issues  []*Issue
	Active  []*Issue
	OldBugs []*Issue
	// ResolutionTimes by ticket type
	ResolutionTimes map[string]TimeRanges
}

// NewReport creates a new report.
func NewReport() *Report {
	report := new(Report)

	// initialize lists to avoid nil issues
	report.Issues = make([]*Issue, 0)
	report.Active = make([]*Issue, 0)
	report.OldBugs = make([]*Issue, 0)
	report.ResolutionTimes = make(map[string]TimeRanges)

	return report
}
