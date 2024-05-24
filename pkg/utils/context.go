package utils

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"net/http"
)

type AleKey struct{}

func GetActionLogCtx(r *http.Request) *models.ActionLogEntry {
	ale, ok := r.Context().Value(AleKey{}).(*models.ActionLogEntry)
	if !ok {
		return nil
	}
	return ale
}

type DleKey struct{}

func GetDisciplinaryLogCtx(r *http.Request) *models.DisciplinaryLogEntry {
	dle, ok := r.Context().Value(DleKey{}).(*models.DisciplinaryLogEntry)
	if !ok {
		return nil
	}
	return dle
}

type DocumentKey struct{}

func GetDocumentCtx(r *http.Request) *models.Document {
	doc, ok := r.Context().Value(DocumentKey{}).(*models.Document)
	if !ok {
		return nil
	}
	return doc
}

type FacilityKey struct{}

func GetFacilityCtx(r *http.Request) *models.Facility {
	fac, ok := r.Context().Value(FacilityKey{}).(*models.Facility)
	if !ok {
		return nil
	}
	return fac
}

type FacilityLogKey struct{}

func GetFacilityLogCtx(r *http.Request) *models.FacilityLogEntry {
	fle, ok := r.Context().Value(FacilityLogKey{}).(*models.FacilityLogEntry)
	if !ok {
		return nil
	}
	return fle
}

type FAQKey struct{}

func GetFAQCtx(r *http.Request) *models.FAQ {
	faq, ok := r.Context().Value(FAQKey{}).(*models.FAQ)
	if !ok {
		return nil
	}
	return faq
}

type FeedbackKey struct{}

func GetFeedbackCtx(r *http.Request) *models.Feedback {
	fb, ok := r.Context().Value(FeedbackKey{}).(*models.Feedback)
	if !ok {
		return nil
	}
	return fb
}

type NewsKey struct{}

func GetNewsCtx(r *http.Request) *models.News {
	news, ok := r.Context().Value(NewsKey{}).(*models.News)
	if !ok {
		return nil
	}
	return news
}

type NotificationKey struct{}

func GetNotificationCtx(r *http.Request) *models.Notification {
	notif, ok := r.Context().Value(NotificationKey{}).(*models.Notification)
	if !ok {
		return nil
	}
	return notif
}

type UserKey struct{}

func GetUserCtx(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserKey{}).(*models.User)
	if !ok {
		return nil
	}
	return user
}

type RosterKey struct{}

func GetRosterCtx(r *http.Request) *models.Roster {
	roster, ok := r.Context().Value(RosterKey{}).(*models.Roster)
	if !ok {
		return nil
	}
	return roster
}

type RatingChangeKey struct{}

func GetRatingChangeCtx(r *http.Request) *models.RatingChange {
	rc, ok := r.Context().Value(RatingChangeKey{}).(*models.RatingChange)
	if !ok {
		return nil
	}
	return rc
}

type RosterRequestKey struct{}

func GetRosterRequestCtx(r *http.Request) *models.RosterRequest {
	rr, ok := r.Context().Value(RosterRequestKey{}).(*models.RosterRequest)
	if !ok {
		return nil
	}
	return rr
}

type UserFlagKey struct{}

func GetUserFlagCtx(r *http.Request) *models.UserFlag {
	uf, ok := r.Context().Value(UserFlagKey{}).(*models.UserFlag)
	if !ok {
		return nil
	}
	return uf
}

type UserRoleKey struct{}

func GetUserRoleCtx(r *http.Request) *models.UserRole {
	ur, ok := r.Context().Value(UserRoleKey{}).(*models.UserRole)
	if !ok {
		return nil
	}
	return ur
}

type XUser struct{}

func GetXUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(XUser{}).(*models.User)
	if !ok {
		return nil
	}
	return user
}

type XGuest struct{}

func GetXGuest(r *http.Request) bool {
	guest, ok := r.Context().Value(XGuest{}).(bool)
	if !ok {
		return false
	}
	return guest
}
