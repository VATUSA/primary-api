package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CanEditFAQ(next http.Handler) http.Handler {
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

			log.Warnf("User %d, attempted to edit faq for facility: %s. No permissions.", credentials.User.CID, targetFacility.ID)
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
