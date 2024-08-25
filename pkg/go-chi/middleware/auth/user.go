package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanViewUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetUser := utils.GetUserCtx(r)

		credentials := GetCredentials(r)
		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			for _, roster := range targetUser.Roster {
				if utils.IsFacilityStaff(credentials.User, roster.Facility) {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Warnf("User %d, attempted to view user: %d. No permissions.", credentials.User.CID, targetUser.CID)
		}

		if credentials.Facility != nil {
			for _, roster := range targetUser.Roster {
				if roster.Facility == credentials.Facility.ID {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Warnf("Facility API Key %s, attempted to view user: %d. No permissions.", credentials.Facility.ID, targetUser.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanEditUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetUser := utils.GetUserCtx(r)

		credentials := GetCredentials(r)
		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to edit user: %d. No permissions.", credentials.User.CID, targetUser.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
