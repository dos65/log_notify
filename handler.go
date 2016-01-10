package main

type Handler interface {
	Handle(string)
}

type NotifyHandler struct {
	notifications []Notification
}

func NewHandler() *NotifyHandler {
	arr := make([]Notification, 0)
	return &NotifyHandler{notifications: arr}
}

func (h *NotifyHandler) Add(notification Notification) {
	h.notifications = append(h.notifications, notification)
}

func (h *NotifyHandler) Handle(text string) {
	for _, notification := range h.notifications {
		go notification.Notify(text)
	}
}
