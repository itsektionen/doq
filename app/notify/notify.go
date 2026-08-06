package notify

import (
	"html/template"
	"net/http"

	"github.com/itsektionen/doq/app/hx"
)

var defaultTmpl *template.Template

func SetTemplate(t *template.Template) {
	defaultTmpl = t
}

type NotificationData struct {
	Title     string
	Message   string
	Type      string
	Timeout   int
	TimeoutMS int
}

type Notifier struct {
	tmpl *template.Template
}

func New(t *template.Template) *Notifier {
	return &Notifier{tmpl: t}
}

func (n *Notifier) Send(w http.ResponseWriter, title, msg, notifType string, timeoutSec int) error {
	hx.Retarget(w, "#notifications")
	hx.Reswap(w, "beforeend")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := n.tmpl.Clone()
	if err != nil {
		return err
	}

	data := NotificationData{
		Title:     title,
		Message:   msg,
		Type:      notifType,
		Timeout:   timeoutSec,
		TimeoutMS: timeoutSec * 1000,
	}

	return tmpl.ExecuteTemplate(w, "notification", data)
}

func Notify(w http.ResponseWriter, r *http.Request, title, msg, notifType string, timeoutSec int) error {
	if defaultTmpl == nil {
		http.Error(w, "notification template not initialized", http.StatusInternalServerError)
		return nil
	}
	return New(defaultTmpl).Send(w, title, msg, notifType, timeoutSec)
}

func NotifyError(w http.ResponseWriter, r *http.Request, msg string) error {
	return Notify(w, r, "Error", msg, "error", 5)
}

func NotifySuccess(w http.ResponseWriter, r *http.Request, msg string) error {
	return Notify(w, r, "Success", msg, "success", 3)
}
