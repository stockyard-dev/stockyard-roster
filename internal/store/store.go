package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Member struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	Department string `json:"department"`
	Phone string `json:"phone"`
	Status string `json:"status"`
	JoinDate string `json:"join_date"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"roster.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS members(id TEXT PRIMARY KEY,name TEXT NOT NULL,email TEXT DEFAULT '',role TEXT DEFAULT '',department TEXT DEFAULT '',phone TEXT DEFAULT '',status TEXT DEFAULT 'active',join_date TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Member)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO members(id,name,email,role,department,phone,status,join_date,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Email,e.Role,e.Department,e.Phone,e.Status,e.JoinDate,e.CreatedAt);return err}
func(d *DB)Get(id string)*Member{var e Member;if d.db.QueryRow(`SELECT id,name,email,role,department,phone,status,join_date,created_at FROM members WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Email,&e.Role,&e.Department,&e.Phone,&e.Status,&e.JoinDate,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Member{rows,_:=d.db.Query(`SELECT id,name,email,role,department,phone,status,join_date,created_at FROM members ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Member;for rows.Next(){var e Member;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Role,&e.Department,&e.Phone,&e.Status,&e.JoinDate,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Member)error{_,err:=d.db.Exec(`UPDATE members SET name=?,email=?,role=?,department=?,phone=?,status=?,join_date=? WHERE id=?`,e.Name,e.Email,e.Role,e.Department,e.Phone,e.Status,e.JoinDate,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM members WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM members`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Member{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR email LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,email,role,department,phone,status,join_date,created_at FROM members WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Member;for rows.Next(){var e Member;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Role,&e.Department,&e.Phone,&e.Status,&e.JoinDate,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM members GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
