package handlers

import (
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/server/utils"
)

type ContactsHandler struct {
	Template template.Template
	DB       db.DB
}

func NewContactsHandler(db db.DB, t template.Template) ContactsHandler {
	return ContactsHandler{
		DB:       db,
		Template: t,
	}
}

func (h ContactsHandler) GetContactsPage(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.DB.GetContacts()
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	_, _, pageData := utils.WithPageData(w, r, map[string]interface{}{
		"Contacts": contacts,
	})
	h.Template.ExecuteTemplate(w, "/contacts", pageData)
}

func (h ContactsHandler) PutContacts(w http.ResponseWriter, r *http.Request) {
	var contacts models.Contacts
	contacts.Links = r.FormValue("links")
	contacts.Update()
	err := h.DB.SetContacts(contacts)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
