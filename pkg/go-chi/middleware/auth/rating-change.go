package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanViewRatingChange(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetUser := utils.GetUserCtx(r)

		credentials := GetCredentials(r)
		if credentials.User != nil {
			if credentials.User.CID == targetUser.CID {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			for _, roster := range targetUser.Roster {
				if utils.IsFacilityStaff(credentials.User, roster.Facility) {
					next.ServeHTTP(w, r)
					return
				}

				if utils.IsInstructor(credentials.User, roster.Facility) {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Warnf("User %d, attempted to view rating change for user: %d. No permissions.", credentials.User.CID, targetUser.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanEditRatingChange(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Render(w, r, utils.ErrForbidden)
	})
}
