package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanViewNotifications(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := GetCredentials(r)
		targetUser := utils.GetXUser(r)

		if credentials.User != nil {
			if credentials.User.CID == targetUser.CID {
				next.ServeHTTP(w, r)
				return
			}

			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to view action log for user: %d. No permissions.", credentials.User.CID, targetUser.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}

func CanEditNotifications(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := GetCredentials(r)
		targetUser := utils.GetXUser(r)

		if credentials.User != nil {
			if utils.IsVATUSAStaff(credentials.User) {
				next.ServeHTTP(w, r)
				return
			}

			log.Warnf("User %d, attempted to view action log for user: %d. No permissions.", credentials.User.CID, targetUser.CID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
