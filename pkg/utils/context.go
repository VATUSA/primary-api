package utils

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AleKey struct{}

func GetActionLogCtx(r *http.Request) *models.ActionLogEntry {
	return r.Context().Value(AleKey{}).(*models.ActionLogEntry)
}

type DleKey struct{}

func GetDisciplinaryLogCtx(r *http.Request) *models.DisciplinaryLogEntry {
	return r.Context().Value(DleKey{}).(*models.DisciplinaryLogEntry)
}

type DocumentKey struct{}

func GetDocumentCtx(r *http.Request) *models.Document {
	return r.Context().Value(DocumentKey{}).(*models.Document)
}

func GetFacilityCtx(r *http.Request) (*models.Facility, error) {
	id := chi.URLParam(r, "facilityID")

	fac := &models.Facility{
		ID: constants.FacilityID(id),
	}

	if err := fac.Get(); err != nil {
		return nil, err
	}

	return fac, nil
}

type FacilityLogKey struct{}

func GetFacilityLogCtx(r *http.Request) *models.FacilityLogEntry {
	return r.Context().Value(FacilityLogKey{}).(*models.FacilityLogEntry)
}

type FAQKey struct{}

func GetFAQCtx(r *http.Request) *models.FAQ {
	return r.Context().Value(FAQKey{}).(*models.FAQ)
}

type FeedbackKey struct{}

func GetFeedbackCtx(r *http.Request) *models.Feedback {
	return r.Context().Value(FeedbackKey{}).(*models.Feedback)
}

type NewsKey struct{}

func GetNewsCtx(r *http.Request) *models.News {
	return r.Context().Value(NewsKey{}).(*models.News)
}

type NotificationKey struct{}

func GetNotificationCtx(r *http.Request) *models.Notification {
	return r.Context().Value(NotificationKey{}).(*models.Notification)
}

type UserKey struct{}

func GetUserCtx(r *http.Request) *models.User {
	return r.Context().Value(UserKey{}).(*models.User)
}

type RosterKey struct{}

func GetRosterCtx(r *http.Request) *models.Roster {
	return r.Context().Value(RosterKey{}).(*models.Roster)
}

type RatingChangeKey struct{}

func GetRatingChangeCtx(r *http.Request) *models.RatingChange {
	return r.Context().Value(RatingChangeKey{}).(*models.RatingChange)
}

type RosterRequestKey struct{}

func GetRosterRequestCtx(r *http.Request) *models.RosterRequest {
	return r.Context().Value(RosterRequestKey{}).(*models.RosterRequest)
}

type UserFlagKey struct{}

func GetUserFlagCtx(r *http.Request) *models.UserFlag {
	return r.Context().Value(UserFlagKey{}).(*models.UserFlag)
}

type UserRoleKey struct{}

func GetUserRoleCtx(r *http.Request) *models.UserRole {
	return r.Context().Value(UserRoleKey{}).(*models.UserRole)
}

type XUser struct{}

func GetXUser(r *http.Request) *models.User {
	return r.Context().Value(XUser{}).(*models.User)
}

type XGuest struct{}

func GetXGuest(r *http.Request) bool {
	return r.Context().Value(XGuest{}).(bool)
}
