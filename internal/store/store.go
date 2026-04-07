package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Member is a person in the team directory. Status is one of:
// active, inactive, on_leave (custom statuses can be added via config).
type Member struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
	JoinDate   string `json:"join_date"`
	CreatedAt  string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "roster.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS members(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT DEFAULT '',
		role TEXT DEFAULT '',
		department TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		status TEXT DEFAULT 'active',
		join_date TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_members_status ON members(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_members_department ON members(department)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Member) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Status == "" {
		e.Status = "active"
	}
	_, err := d.db.Exec(
		`INSERT INTO members(id, name, email, role, department, phone, status, join_date, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Email, e.Role, e.Department, e.Phone, e.Status, e.JoinDate, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Member {
	var e Member
	err := d.db.QueryRow(
		`SELECT id, name, email, role, department, phone, status, join_date, created_at
		 FROM members WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Department, &e.Phone, &e.Status, &e.JoinDate, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Member {
	rows, _ := d.db.Query(
		`SELECT id, name, email, role, department, phone, status, join_date, created_at
		 FROM members ORDER BY name ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Member
	for rows.Next() {
		var e Member
		rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Department, &e.Phone, &e.Status, &e.JoinDate, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Member) error {
	_, err := d.db.Exec(
		`UPDATE members SET name=?, email=?, role=?, department=?, phone=?, status=?, join_date=?
		 WHERE id=?`,
		e.Name, e.Email, e.Role, e.Department, e.Phone, e.Status, e.JoinDate, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM members WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM members`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Member {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR email LIKE ? OR role LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	if v, ok := filters["department"]; ok && v != "" {
		where += " AND department=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, email, role, department, phone, status, join_date, created_at
		 FROM members WHERE `+where+`
		 ORDER BY name ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Member
	for rows.Next() {
		var e Member
		rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Department, &e.Phone, &e.Status, &e.JoinDate, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns aggregate metrics: total members, counts by status,
// counts by department, and the active member count.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":         d.Count(),
		"active":        0,
		"by_status":     map[string]int{},
		"by_department": map[string]int{},
	}

	var active int
	d.db.QueryRow(`SELECT COUNT(*) FROM members WHERE status='active'`).Scan(&active)
	m["active"] = active

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM members GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	if rows, _ := d.db.Query(`SELECT department, COUNT(*) FROM members WHERE department != '' GROUP BY department`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_department"] = by
	}

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
