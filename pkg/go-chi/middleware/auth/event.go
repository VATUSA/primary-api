package middleware

import (
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanEditEvent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetFacility := utils.GetFacilityCtx(r)
		credentials := GetCredentials(r)
		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsFacilitySeniorStaff(credentials.User, targetFacility.ID) {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsFacilityEventsStaff(credentials.User, targetFacility.ID) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to edit event: %s. No permissions.", credentials.User.CID, targetFacility.ID)
		}

		if credentials.Facility != nil {
			if credentials.Facility.ID == targetFacility.ID {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("Facility API Key %s, attempted to edit event: %s. No permissions.", credentials.Facility.ID, targetFacility.ID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

type EventSignupRequest struct {
	CID uint `json:"cid"`
}

func CanEventSignup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetFacility := utils.GetFacilityCtx(r)
		req := EventSignupRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		credentials := GetCredentials(r)
		if credentials.User != nil {
			if credentials.User.CID == req.CID {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to signup for event as: %d. No permissions.", credentials.User.CID, req.CID)
		}

		if credentials.Facility != nil {
			if credentials.Facility.ID == targetFacility.ID {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("Facility API Key %s, attempted to signup for event: %s. No permissions.", credentials.Facility.ID, targetFacility.ID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanDeleteEventSignup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetFacility := utils.GetFacilityCtx(r)
		signup := utils.GetEventSignupCtx(r)
		credentials := GetCredentials(r)
		if credentials.User != nil {
			if credentials.User.CID == signup.CID {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsFacilitySeniorStaff(credentials.User, targetFacility.ID) {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsFacilityEventsStaff(credentials.User, targetFacility.ID) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to delete event signup: %d. No permissions.", credentials.User.CID, signup.ID)
		}

		if credentials.Facility != nil {
			if credentials.Facility.ID == targetFacility.ID {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("Facility API Key %s, attempted to delete event signup: %d. No permissions.", credentials.Facility.ID, signup.ID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
