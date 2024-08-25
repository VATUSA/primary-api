package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanReadDisciplinaryLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vatusaOnly := r.URL.Query().Get("vatusa_only")

		targetUser := utils.GetUserCtx(r)

		credentials := GetCredentials(r)

		if vatusaOnly == "true" {
			if credentials.User != nil {
				if utils.IsVATUSAStaff(credentials.User) {
					next.ServeHTTP(w, r)
					return
				}
			}
		} else {
			if credentials.User != nil {
				if utils.IsVATUSAStaff(credentials.User) {
					next.ServeHTTP(w, r)
					return
				}

				for _, roster := range targetUser.Roster {
					if utils.IsFacilitySeniorStaff(credentials.User, roster.Facility) {
						next.ServeHTTP(w, r)
						return
					}
				}

				log.Warnf("User %d, attempted to view disciplinary log for user: %d. No permissions.", credentials.User.CID, targetUser.CID)
			}

			if credentials.Facility != nil {
				for _, roster := range targetUser.Roster {
					if roster.Facility == credentials.Facility.ID {
						next.ServeHTTP(w, r)
						return
					}
				}

				log.Warnf("Facility API Key %s, attempted to view disciplinary log for user's disciplinary log: %d. No permissions.", credentials.Facility.ID, targetUser.CID)
			}
		}
		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanEditDisciplinaryLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := GetCredentials(r)
		targetUser := utils.GetUserCtx(r)

		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}
		}

		log.Warnf("User %d, attempted to edit disciplinary log for user: %d. No permissions.", credentials.User.CID, targetUser.CID)

		utils.Render(w, r, utils.ErrForbidden)
	})
}
