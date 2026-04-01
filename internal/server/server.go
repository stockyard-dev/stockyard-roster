package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stockyard-dev/stockyard-roster/internal/store"
)

type Server struct {
	db     *store.DB
	mux    *http.ServeMux
	port   int
	limits Limits
}

func New(db *store.DB, port int, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), port: port, limits: limits}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /api/contacts", s.handleCreateContact)
	s.mux.HandleFunc("GET /api/contacts", s.handleListContacts)
	s.mux.HandleFunc("GET /api/contacts/{id}", s.handleGetContact)
	s.mux.HandleFunc("PUT /api/contacts/{id}", s.handleUpdateContact)
	s.mux.HandleFunc("DELETE /api/contacts/{id}", s.handleDeleteContact)

	s.mux.HandleFunc("POST /api/contacts/{id}/activities", s.handleAddActivity)
	s.mux.HandleFunc("GET /api/contacts/{id}/activities", s.handleListActivities)

	s.mux.HandleFunc("POST /api/deals", s.handleCreateDeal)
	s.mux.HandleFunc("GET /api/deals", s.handleListDeals)
	s.mux.HandleFunc("PUT /api/deals/{id}/stage", s.handleUpdateDealStage)
	s.mux.HandleFunc("DELETE /api/deals/{id}", s.handleDeleteDeal)

	s.mux.HandleFunc("POST /api/reminders", s.handleCreateReminder)
	s.mux.HandleFunc("GET /api/reminders", s.handleListReminders)
	s.mux.HandleFunc("POST /api/reminders/{id}/done", s.handleCompleteReminder)

	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /ui", s.handleUI)
	s.mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"product": "stockyard-roster", "version": "0.1.0"})
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("[roster] listening on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleCreateContact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Company string `json:"company"`
		Title   string `json:"title"`
		Stage   string `json:"stage"`
		Tags    string `json:"tags"`
		Notes   string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeJSON(w, 400, map[string]string{"error": "name is required"})
		return
	}
	if s.limits.MaxContacts > 0 && LimitReached(s.limits.MaxContacts, s.db.TotalContacts()) {
		writeJSON(w, 402, map[string]string{"error": fmt.Sprintf("free tier limit: %d contacts — upgrade to Pro", s.limits.MaxContacts), "upgrade": "https://stockyard.dev/roster/"})
		return
	}
	c, err := s.db.CreateContact(req.Name, req.Email, req.Phone, req.Company, req.Title, req.Stage, req.Tags, req.Notes)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, map[string]any{"contact": c})
}

func (s *Server) handleListContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := s.db.ListContacts(r.URL.Query().Get("stage"))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	if contacts == nil { contacts = []store.Contact{} }
	writeJSON(w, 200, map[string]any{"contacts": contacts, "count": len(contacts)})
}

func (s *Server) handleGetContact(w http.ResponseWriter, r *http.Request) {
	c, err := s.db.GetContact(r.PathValue("id"))
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "contact not found"})
		return
	}
	activities, _ := s.db.ListActivities(c.ID)
	if activities == nil { activities = []store.Activity{} }
	writeJSON(w, 200, map[string]any{"contact": c, "activities": activities})
}

func (s *Server) handleUpdateContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := s.db.GetContact(id); err != nil {
		writeJSON(w, 404, map[string]string{"error": "contact not found"})
		return
	}
	var req struct {
		Name *string `json:"name"`; Email *string `json:"email"`; Phone *string `json:"phone"`
		Company *string `json:"company"`; Title *string `json:"title"`; Stage *string `json:"stage"`
		Tags *string `json:"tags"`; Notes *string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	c, _ := s.db.UpdateContact(id, req.Name, req.Email, req.Phone, req.Company, req.Title, req.Stage, req.Tags, req.Notes)
	writeJSON(w, 200, map[string]any{"contact": c})
}

func (s *Server) handleDeleteContact(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteContact(r.PathValue("id"))
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) handleAddActivity(w http.ResponseWriter, r *http.Request) {
	cid := r.PathValue("id")
	if _, err := s.db.GetContact(cid); err != nil {
		writeJSON(w, 404, map[string]string{"error": "contact not found"})
		return
	}
	var req struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeJSON(w, 400, map[string]string{"error": "content is required"})
		return
	}
	if req.Type == "" { req.Type = "note" }
	a, _ := s.db.AddActivity(cid, req.Type, req.Content)
	writeJSON(w, 201, map[string]any{"activity": a})
}

func (s *Server) handleListActivities(w http.ResponseWriter, r *http.Request) {
	acts, _ := s.db.ListActivities(r.PathValue("id"))
	if acts == nil { acts = []store.Activity{} }
	writeJSON(w, 200, map[string]any{"activities": acts, "count": len(acts)})
}

func (s *Server) handleCreateDeal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContactID  string `json:"contact_id"`
		Title      string `json:"title"`
		ValueCents int    `json:"value_cents"`
		Stage      string `json:"stage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.ContactID == "" {
		writeJSON(w, 400, map[string]string{"error": "contact_id and title are required"})
		return
	}
	d, err := s.db.CreateDeal(req.ContactID, req.Title, req.ValueCents, req.Stage)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, map[string]any{"deal": d})
}

func (s *Server) handleListDeals(w http.ResponseWriter, r *http.Request) {
	deals, _ := s.db.ListDeals(r.URL.Query().Get("stage"))
	if deals == nil { deals = []store.Deal{} }
	writeJSON(w, 200, map[string]any{"deals": deals, "count": len(deals)})
}

func (s *Server) handleUpdateDealStage(w http.ResponseWriter, r *http.Request) {
	var req struct { Stage string `json:"stage"` }
	json.NewDecoder(r.Body).Decode(&req)
	s.db.UpdateDealStage(r.PathValue("id"), req.Stage)
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (s *Server) handleDeleteDeal(w http.ResponseWriter, r *http.Request) {
	s.db.DeleteDeal(r.PathValue("id"))
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) handleCreateReminder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContactID string `json:"contact_id"`
		Content   string `json:"content"`
		DueAt     string `json:"due_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeJSON(w, 400, map[string]string{"error": "content is required"})
		return
	}
	rem, _ := s.db.CreateReminder(req.ContactID, req.Content, req.DueAt)
	writeJSON(w, 201, map[string]any{"reminder": rem})
}

func (s *Server) handleListReminders(w http.ResponseWriter, r *http.Request) {
	rems, _ := s.db.ListReminders(r.URL.Query().Get("due") == "true")
	if rems == nil { rems = []store.Reminder{} }
	writeJSON(w, 200, map[string]any{"reminders": rems, "count": len(rems)})
}

func (s *Server) handleCompleteReminder(w http.ResponseWriter, r *http.Request) {
	s.db.CompleteReminder(r.PathValue("id"))
	writeJSON(w, 200, map[string]string{"status": "done"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Stats())
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
