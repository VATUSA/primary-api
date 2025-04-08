package training

type Rating string

const (
	Empty          Rating = ""
	Satisfactory   Rating = "Satisfactory"
	Unsatisfactory Rating = "Unsatisfactory"
)

// KPI - Key Performance Indicator
type KPI struct {
	Order       uint   `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Rating      `json:"rating"`
}

// KPA - Key Performance Area
type KPA struct {
	Order       uint   `json:"order"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Indicators  []KPI  `json:"indicators"`
}

// KPC - Key performance Category
type KPC struct {
	Order uint   `json:"order"`
	Title string `json:"title"`
	Areas []KPA  `json:"areas"`
}
type OTS struct {
	Title      string `json:"title"`
	Categories []KPC  `json:"categories"`
}

var S2OTS = OTS{
	Title: "S2 (Tower) Rating Form",
	Categories: []KPC{
		{
			Order: 1,
			Title: "Theory",
			Areas: []KPA{
				{
					Order:       1,
					Title:       "Clearance Delivery + Ground Control",
					Description: "Demonstrates knowledge of Delivery and Ground Controller duties and responsibilities",
					Indicators: []KPI{
						{
							Order:       1,
							Title:       "Defines all parts of a clearance",
							Description: "...",
							Required:    true,
							Rating:      Empty,
						},
						{
							Order:       2,
							Title:       "Explains all types of SIDs",
							Description: "...",
							Required:    true,
							Rating:      Empty,
						},
					},
				},
				{
					Order:       2,
					Title:       "Local Control",
					Description: "Demonstrates knowledge of Local Controller duties and responsibilities",
					Indicators: []KPI{
						{
							Order:       1,
							Title:       "Identifies difference between movement and non-movement areas",
							Description: "...",
							Required:    true,
							Rating:      Empty,
						},
						{
							Order:       2,
							Title:       "Defines all parts of VFR traffic pattern",
							Description: "...",
							Required:    true,
							Rating:      Empty,
						},
					},
				},
			},
		},
	},
}
