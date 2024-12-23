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

type EventKey struct{}

func GetEventCtx(r *http.Request) *models.Event {
	ev, ok := r.Context().Value(EventKey{}).(*models.Event)
	if !ok {
		return nil
	}
	return ev
}

type EventPositionKey struct{}

func GetEventPositionCtx(r *http.Request) *models.EventPosition {
	ep, ok := r.Context().Value(EventPositionKey{}).(*models.EventPosition)
	if !ok {
		return nil
	}
	return ep
}

type EventSignupKey struct{}

func GetEventSignupCtx(r *http.Request) *models.EventSignup {
	es, ok := r.Context().Value(EventSignupKey{}).(*models.EventSignup)
	if !ok {
		return nil
	}
	return es
}

type EventRoutingKey struct{}

func GetEventRoutingCtx(r *http.Request) *models.EventRouting {
	er, ok := r.Context().Value(EventRoutingKey{}).(*models.EventRouting)
	if !ok {
		return nil
	}
	return er
}

type EventTemplateKey struct{}

func GetEventTemplateCtx(r *http.Request) *models.EventTemplate {
	et, ok := r.Context().Value(EventTemplateKey{}).(*models.EventTemplate)
	if !ok {
		return nil
	}
	return et
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

type UserNotificationKey struct{}

func GetUserNotificationCtx(r *http.Request) *models.UserNotification {
	un, ok := r.Context().Value(UserNotificationKey{}).(*models.UserNotification)
	if !ok {
		return nil
	}
	return un
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

type XFacility struct{}

func GetXFacility(r *http.Request) *models.Facility {
	fac, ok := r.Context().Value(XFacility{}).(*models.Facility)
	if !ok {
		return nil
	}
	return fac
}

type XGuest struct{}

func GetXGuest(r *http.Request) bool {
	guest, ok := r.Context().Value(XGuest{}).(bool)
	if !ok {
		return false
	}
	return guest
}
