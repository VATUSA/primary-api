package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanEditOTSTemplate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := GetCredentials(r)
		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to edit an OTS.", credentials.User.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanCROTSRecord(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetUser := utils.GetUserCtx(r)

		credentials := GetCredentials(r)
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

				if utils.IsInstructor(credentials.User, roster.Facility) {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Warnf("User %d, attempted to create/read an OTS.", credentials.User.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
