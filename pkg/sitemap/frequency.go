package sitemap

type Frequency int64

const (
	Always  Frequency = 0
	Hourly  Frequency = 1
	Daily   Frequency = 2
	Weekly  Frequency = 3
	Monthly Frequency = 4
	Yearly  Frequency = 5
	Never   Frequency = 6
)

func (freq Frequency) String() string {
	return []string{
		"Always",
		"Hourly",
		"Daily",
		"Weekly",
		"Monthly",
		"Yearly",
		"Never",
	}[freq]
}
