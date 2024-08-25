package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func CanViewFeedback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := GetCredentials(r)

		targetFacility := utils.GetFacilityCtx(r)

		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsFacilitySeniorStaff(credentials.User, targetFacility.ID) {
				next.ServeHTTP(w, r)
				return
			}

			cid := r.URL.Query().Get("cid")
			if cid != "" {
				cidInt, err := strconv.Atoi(cid)
				if err != nil {
					utils.Render(w, r, utils.ErrInvalidRequest(err))
					return
				}

				if credentials.User.CID == uint(cidInt) {
					q := r.URL.Query()
					q.Set("cid", cid)           // Preserve the original 'cid'
					q.Set("status", "accepted") // Set the 'status' to 'accepted'
					r.URL.RawQuery = q.Encode()
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Warnf("User %d, attempted to view feedback for facility: %s. No permissions.", credentials.User.CID, targetFacility.ID)
		}

		if credentials.Facility != nil {
			if credentials.Facility.ID == targetFacility.ID {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("Facility %s, attempted to view feedback for facility: %s. No permissions.", credentials.Facility.ID, targetFacility.ID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanEditFeedback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := GetCredentials(r)

		targetFacility := utils.GetFacilityCtx(r)

		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsFacilitySeniorStaff(credentials.User, targetFacility.ID) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to edit feedback for facility: %s. No permissions.", credentials.User.CID, targetFacility.ID)
		}

		if credentials.Facility != nil {
			if credentials.Facility.ID == targetFacility.ID {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("Facility %s, attempted to edit feedback for facility: %s. No permissions.", credentials.Facility.ID, targetFacility.ID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
