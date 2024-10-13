package event

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

type EventTemplateRequest struct {
	Title      string                 `json:"title" example:"KDEN FNO" validate:"required"`
	Positions  []string               `json:"positions" example:"[\"ZDV_APP\", \"ZDV_TWR\"]" validate:"required"`
	Facilities []constants.FacilityID `json:"facilities" example:"[\"ZDV\", \"ZAB\", \"ZLC\"]" validate:"required"`
	Fields     []string               `json:"fields" example:"[\"KDEN\", \"KBJC\", \"KAPA\"]" validate:"required"`
	Shifts     bool                   `json:"shifts" example:"true"`
}

func (req *EventTemplateRequest) Validate() error {
	return nil
}

func (req *EventTemplateRequest) Bind(r *http.Request) error {
	return nil
}

type EventTemplateResponse struct {
	*models.EventTemplate
}

func NewEventTemplateResponse(et *models.EventTemplate) *EventTemplateResponse {
	return &EventTemplateResponse{EventTemplate: et}
}

func (res *EventTemplateResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.EventTemplate == nil {
		return nil
	}
	return nil
}

func NewEventTemplateListResponse(ets []models.EventTemplate) []render.Renderer {
	list := []render.Renderer{}
	for idx := range ets {
		list = append(list, NewEventTemplateResponse(&ets[idx]))
	}
	return list
}

type EventRequest struct {
	Title       string                 `json:"title" example:"ZDV FNO" validate:"required"`
	Description string                 `json:"description" example:"Join us for a fun night of flying in and out of Denver!" validate:"required"`
	BannerURL   string                 `json:"banner_url" example:"https://zdvartcc.org/banner.jpg"`
	StartDate   time.Time              `json:"start_date" example:"2021-01-01T00:00:00Z" validate:"required"`
	EndDate     time.Time              `json:"end_date" example:"2021-01-01T00:00:00Z" validate:"required"`
	Fields      []string               `json:"fields" example:"[\"KDEN\", \"KBJC\", \"KAPA\"]" validate:"required"`
	Facilities  []constants.FacilityID `json:"facilities" example:"[\"ZDV\", \"ZAB\", \"ZLC\"]" validate:"required"`
}

func (req *EventRequest) Validate() error {
	return nil
}

func (req *EventRequest) Bind(r *http.Request) error {
	return nil
}

type EventResponse struct {
	*models.Event
}

func NewEventResponse(e *models.Event) *EventResponse {
	return &EventResponse{Event: e}
}

func (res *EventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Event == nil {
		return nil
	}
	return nil
}

func NewEventListResponse(es []models.Event) []render.Renderer {
	list := []render.Renderer{}
	for idx := range es {
		list = append(list, NewEventResponse(&es[idx]))
	}
	return list
}

type EventPositionRequest struct {
	Position          string `json:"position" example:"ZDV_APP" validate:"required"`
	Assignee          *uint  `json:"assignee" example:"1293257"`
	Shifts            *bool  `json:"shifts" example:"true"`
	SecondaryAssignee *uint  `json:"secondary_assignee" example:"1293257"`
}

func (req *EventPositionRequest) Validate() error {
	return nil
}

func (req *EventPositionRequest) Bind(r *http.Request) error {
	return nil
}

type EventPositionResponse struct {
	*models.EventPosition
}

func NewEventPositionResponse(ep *models.EventPosition) *EventPositionResponse {
	return &EventPositionResponse{EventPosition: ep}
}

func (res *EventPositionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.EventPosition == nil {
		return nil
	}
	return nil
}

func NewEventPositionListResponse(eps []models.EventPosition) []render.Renderer {
	list := []render.Renderer{}
	for idx := range eps {
		list = append(list, NewEventPositionResponse(&eps[idx]))
	}
	return list
}

type EventSignupRequest struct {
	PositionID uint `json:"position_id" example:"1" validate:"required"`
	CID        uint `json:"cid" example:"1293257" validate:"required"`
	Shift      uint `json:"shift" example:"1" validate:"required"` // 1 = Primary, 2 = Secondary
}

func (req *EventSignupRequest) Validate() error {
	return nil
}

func (req *EventSignupRequest) Bind(r *http.Request) error {
	return nil
}

type EventSignupResponse struct {
	*models.EventSignup
}

func NewEventSignupResponse(es *models.EventSignup) *EventSignupResponse {
	return &EventSignupResponse{EventSignup: es}
}

func (res *EventSignupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.EventSignup == nil {
		return nil
	}
	return nil
}

func NewEventSignupListResponse(ess []models.EventSignup) []render.Renderer {
	list := []render.Renderer{}
	for idx := range ess {
		list = append(list, NewEventSignupResponse(&ess[idx]))
	}
	return list
}

type EventRoutingRequest struct {
	Origin      string `json:"origin" example:"ZDV" validate:"required"`
	Destination string `json:"destination" example:"ZAB" validate:"required"`
	Routing     string `json:"routing" example:"ZDV J80 DBL J80 FQF J80 HCT J80 HBU J80 HCT J80 FQF J80 DBL J80 ZAB" validate:"required"`
	Notes       string `json:"notes" example:"Expect vectors to final at DBL" validate:"required"`
}

func (req *EventRoutingRequest) Validate() error {
	return nil
}

func (req *EventRoutingRequest) Bind(r *http.Request) error {
	return nil
}

type EventRoutingResponse struct {
	*models.EventRouting
}

func NewEventRoutingResponse(er *models.EventRouting) *EventRoutingResponse {
	return &EventRoutingResponse{EventRouting: er}
}

func (res *EventRoutingResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.EventRouting == nil {
		return nil
	}
	return nil
}

func NewEventRoutingListResponse(ers []models.EventRouting) []render.Renderer {
	list := []render.Renderer{}
	for idx := range ers {
		list = append(list, NewEventRoutingResponse(&ers[idx]))
	}
	return list
}
