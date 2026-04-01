package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ conn *sql.DB }

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	conn, err := sql.Open("sqlite", filepath.Join(dataDir, "roster.db"))
	if err != nil {
		return nil, err
	}
	conn.Exec("PRAGMA journal_mode=WAL")
	conn.Exec("PRAGMA busy_timeout=5000")
	conn.SetMaxOpenConns(4)
	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() error { return db.conn.Close() }

func (db *DB) migrate() error {
	_, err := db.conn.Exec(`
CREATE TABLE IF NOT EXISTS contacts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT DEFAULT '',
    phone TEXT DEFAULT '',
    company TEXT DEFAULT '',
    title TEXT DEFAULT '',
    stage TEXT DEFAULT 'lead',
    tags TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS activities (
    id TEXT PRIMARY KEY,
    contact_id TEXT NOT NULL,
    type TEXT DEFAULT 'note',
    content TEXT DEFAULT '',
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_act_contact ON activities(contact_id);

CREATE TABLE IF NOT EXISTS deals (
    id TEXT PRIMARY KEY,
    contact_id TEXT NOT NULL,
    title TEXT NOT NULL,
    value_cents INTEGER DEFAULT 0,
    stage TEXT DEFAULT 'prospect',
    closed_at TEXT DEFAULT '',
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_deals_contact ON deals(contact_id);
CREATE INDEX IF NOT EXISTS idx_deals_stage ON deals(stage);

CREATE TABLE IF NOT EXISTS reminders (
    id TEXT PRIMARY KEY,
    contact_id TEXT NOT NULL,
    content TEXT NOT NULL,
    due_at TEXT NOT NULL,
    done INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_rem_due ON reminders(done, due_at);
`)
	return err
}

// --- Contacts ---

type Contact struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Title     string `json:"title"`
	Stage     string `json:"stage"`
	Tags      string `json:"tags"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (db *DB) CreateContact(name, email, phone, company, title, stage, tags, notes string) (*Contact, error) {
	id := "con_" + genID(8)
	now := time.Now().UTC().Format(time.RFC3339)
	if stage == "" {
		stage = "lead"
	}
	_, err := db.conn.Exec("INSERT INTO contacts (id,name,email,phone,company,title,stage,tags,notes,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		id, name, email, phone, company, title, stage, tags, notes, now, now)
	if err != nil {
		return nil, err
	}
	return &Contact{ID: id, Name: name, Email: email, Phone: phone, Company: company,
		Title: title, Stage: stage, Tags: tags, Notes: notes, CreatedAt: now, UpdatedAt: now}, nil
}

func (db *DB) ListContacts(stageFilter string) ([]Contact, error) {
	var rows *sql.Rows
	var err error
	if stageFilter != "" {
		rows, err = db.conn.Query("SELECT id,name,email,phone,company,title,stage,tags,notes,created_at,updated_at FROM contacts WHERE stage=? ORDER BY updated_at DESC", stageFilter)
	} else {
		rows, err = db.conn.Query("SELECT id,name,email,phone,company,title,stage,tags,notes,created_at,updated_at FROM contacts ORDER BY updated_at DESC")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Contact
	for rows.Next() {
		var c Contact
		rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Company, &c.Title, &c.Stage, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (db *DB) GetContact(id string) (*Contact, error) {
	var c Contact
	err := db.conn.QueryRow("SELECT id,name,email,phone,company,title,stage,tags,notes,created_at,updated_at FROM contacts WHERE id=?", id).
		Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Company, &c.Title, &c.Stage, &c.Tags, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (db *DB) UpdateContact(id string, name, email, phone, company, title, stage, tags, notes *string) (*Contact, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	if name != nil { db.conn.Exec("UPDATE contacts SET name=?, updated_at=? WHERE id=?", *name, now, id) }
	if email != nil { db.conn.Exec("UPDATE contacts SET email=?, updated_at=? WHERE id=?", *email, now, id) }
	if phone != nil { db.conn.Exec("UPDATE contacts SET phone=?, updated_at=? WHERE id=?", *phone, now, id) }
	if company != nil { db.conn.Exec("UPDATE contacts SET company=?, updated_at=? WHERE id=?", *company, now, id) }
	if title != nil { db.conn.Exec("UPDATE contacts SET title=?, updated_at=? WHERE id=?", *title, now, id) }
	if stage != nil { db.conn.Exec("UPDATE contacts SET stage=?, updated_at=? WHERE id=?", *stage, now, id) }
	if tags != nil { db.conn.Exec("UPDATE contacts SET tags=?, updated_at=? WHERE id=?", *tags, now, id) }
	if notes != nil { db.conn.Exec("UPDATE contacts SET notes=?, updated_at=? WHERE id=?", *notes, now, id) }
	return db.GetContact(id)
}

func (db *DB) DeleteContact(id string) error {
	db.conn.Exec("DELETE FROM activities WHERE contact_id=?", id)
	db.conn.Exec("DELETE FROM deals WHERE contact_id=?", id)
	db.conn.Exec("DELETE FROM reminders WHERE contact_id=?", id)
	_, err := db.conn.Exec("DELETE FROM contacts WHERE id=?", id)
	return err
}

func (db *DB) TotalContacts() int {
	var c int
	db.conn.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&c)
	return c
}

// --- Activities ---

type Activity struct {
	ID        string `json:"id"`
	ContactID string `json:"contact_id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (db *DB) AddActivity(contactID, atype, content string) (*Activity, error) {
	id := "act_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO activities (id,contact_id,type,content,created_at) VALUES (?,?,?,?,?)",
		id, contactID, atype, content, now)
	if err != nil {
		return nil, err
	}
	db.conn.Exec("UPDATE contacts SET updated_at=? WHERE id=?", now, contactID)
	return &Activity{ID: id, ContactID: contactID, Type: atype, Content: content, CreatedAt: now}, nil
}

func (db *DB) ListActivities(contactID string) ([]Activity, error) {
	rows, err := db.conn.Query("SELECT id,contact_id,type,content,created_at FROM activities WHERE contact_id=? ORDER BY created_at DESC LIMIT 50", contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Activity
	for rows.Next() {
		var a Activity
		rows.Scan(&a.ID, &a.ContactID, &a.Type, &a.Content, &a.CreatedAt)
		out = append(out, a)
	}
	return out, rows.Err()
}

// --- Deals ---

type Deal struct {
	ID         string `json:"id"`
	ContactID  string `json:"contact_id"`
	Title      string `json:"title"`
	ValueCents int    `json:"value_cents"`
	Stage      string `json:"stage"`
	ClosedAt   string `json:"closed_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

func (db *DB) CreateDeal(contactID, title string, valueCents int, stage string) (*Deal, error) {
	id := "deal_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	if stage == "" {
		stage = "prospect"
	}
	_, err := db.conn.Exec("INSERT INTO deals (id,contact_id,title,value_cents,stage,created_at) VALUES (?,?,?,?,?,?)",
		id, contactID, title, valueCents, stage, now)
	if err != nil {
		return nil, err
	}
	return &Deal{ID: id, ContactID: contactID, Title: title, ValueCents: valueCents, Stage: stage, CreatedAt: now}, nil
}

func (db *DB) ListDeals(stageFilter string) ([]Deal, error) {
	var rows *sql.Rows
	var err error
	if stageFilter != "" {
		rows, err = db.conn.Query("SELECT id,contact_id,title,value_cents,stage,closed_at,created_at FROM deals WHERE stage=? ORDER BY created_at DESC", stageFilter)
	} else {
		rows, err = db.conn.Query("SELECT id,contact_id,title,value_cents,stage,closed_at,created_at FROM deals ORDER BY created_at DESC")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Deal
	for rows.Next() {
		var d Deal
		rows.Scan(&d.ID, &d.ContactID, &d.Title, &d.ValueCents, &d.Stage, &d.ClosedAt, &d.CreatedAt)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (db *DB) UpdateDealStage(id, stage string) {
	if stage == "won" || stage == "lost" {
		now := time.Now().UTC().Format(time.RFC3339)
		db.conn.Exec("UPDATE deals SET stage=?, closed_at=? WHERE id=?", stage, now, id)
	} else {
		db.conn.Exec("UPDATE deals SET stage=? WHERE id=?", stage, id)
	}
}

func (db *DB) DeleteDeal(id string) error {
	_, err := db.conn.Exec("DELETE FROM deals WHERE id=?", id)
	return err
}

// --- Reminders ---

type Reminder struct {
	ID        string `json:"id"`
	ContactID string `json:"contact_id"`
	Content   string `json:"content"`
	DueAt     string `json:"due_at"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

func (db *DB) CreateReminder(contactID, content, dueAt string) (*Reminder, error) {
	id := "rem_" + genID(6)
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.conn.Exec("INSERT INTO reminders (id,contact_id,content,due_at,created_at) VALUES (?,?,?,?,?)",
		id, contactID, content, dueAt, now)
	if err != nil {
		return nil, err
	}
	return &Reminder{ID: id, ContactID: contactID, Content: content, DueAt: dueAt, CreatedAt: now}, nil
}

func (db *DB) ListReminders(dueOnly bool) ([]Reminder, error) {
	var rows *sql.Rows
	var err error
	if dueOnly {
		rows, err = db.conn.Query("SELECT id,contact_id,content,due_at,done,created_at FROM reminders WHERE done=0 ORDER BY due_at ASC")
	} else {
		rows, err = db.conn.Query("SELECT id,contact_id,content,due_at,done,created_at FROM reminders ORDER BY due_at DESC LIMIT 50")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Reminder
	for rows.Next() {
		var r Reminder
		var d int
		rows.Scan(&r.ID, &r.ContactID, &r.Content, &r.DueAt, &d, &r.CreatedAt)
		r.Done = d == 1
		out = append(out, r)
	}
	return out, rows.Err()
}

func (db *DB) CompleteReminder(id string) {
	db.conn.Exec("UPDATE reminders SET done=1 WHERE id=?", id)
}

// --- Stats ---

func (db *DB) Stats() map[string]any {
	var contacts, deals, pipelineValue, reminders int
	db.conn.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&contacts)
	db.conn.QueryRow("SELECT COUNT(*) FROM deals WHERE stage NOT IN ('won','lost')").Scan(&deals)
	db.conn.QueryRow("SELECT COALESCE(SUM(value_cents),0) FROM deals WHERE stage NOT IN ('won','lost')").Scan(&pipelineValue)
	db.conn.QueryRow("SELECT COUNT(*) FROM reminders WHERE done=0").Scan(&reminders)
	return map[string]any{"contacts": contacts, "active_deals": deals, "pipeline_cents": pipelineValue, "pending_reminders": reminders}
}

func genID(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
