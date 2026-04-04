package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Roster</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(250px,1fr));gap:.6rem}
.card{background:var(--bg2);border:1px solid var(--bg3);padding:1rem}
.card-name{font-size:.9rem;color:var(--cream);margin-bottom:.2rem}
.card-role{font-size:.7rem;color:var(--gold)}
.card-dept{font-size:.65rem;color:var(--cm);margin-top:.1rem}
.card-contact{font-size:.65rem;color:var(--cd);margin-top:.4rem}
.card-contact a{color:var(--rust)}
.badge{font-size:.5rem;padding:.1rem .3rem;text-transform:uppercase;letter-spacing:1px}
.badge-active{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.badge-inactive{background:var(--bg3);color:var(--cm);border:1px solid var(--bg3)}
.btn{font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:380px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.search{width:100%;padding:.5rem .8rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem;margin-bottom:1rem}
.dept-bar{display:flex;gap:.3rem;margin-bottom:1rem;flex-wrap:wrap}
.dept-btn{font-size:.6rem;padding:.2rem .5rem;border:1px solid var(--bg3);background:var(--bg);color:var(--cm);cursor:pointer}.dept-btn:hover{border-color:var(--leather)}.dept-btn.active{border-color:var(--rust);color:var(--rust)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>ROSTER</h1><button class="btn btn-p" onclick="openForm()">+ Add Member</button></div>
<div class="main">
<input class="search" id="search" placeholder="Search by name, role, or department..." oninput="render()">
<div class="dept-bar" id="depts"></div>
<div class="grid" id="grid"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let members=[],filterDept='';
async function load(){const r=await fetch(A+'/members').then(r=>r.json());members=r.members||[];
const depts=[...new Set(members.map(m=>m.department).filter(d=>d))];
let dh='<button class="dept-btn'+(filterDept===''?' active':'')+'" onclick="setDept(\'\')">All ('+members.length+')</button>';
depts.forEach(d=>{const c=members.filter(m=>m.department===d).length;dh+='<button class="dept-btn'+(filterDept===d?' active':'')+'" onclick="setDept(\''+d+'\')">'+esc(d)+' ('+c+')</button>';});
document.getElementById('depts').innerHTML=dh;render();}
function setDept(d){filterDept=d;document.querySelectorAll('.dept-btn').forEach((b,i)=>{const depts=['',... new Set(members.map(m=>m.department).filter(d=>d))];b.classList.toggle('active',depts[i]===d)});render();}
function render(){const q=(document.getElementById('search').value||'').toLowerCase();
let filtered=members.filter(m=>{if(filterDept&&m.department!==filterDept)return false;if(q&&!(m.name+m.role+m.department+m.email).toLowerCase().includes(q))return false;return true;});
if(!filtered.length){document.getElementById('grid').innerHTML='<div class="empty">No team members'+(q?' matching "'+esc(q)+'"':'')+'.</div>';return;}
let h='';filtered.forEach(m=>{h+='<div class="card"><div style="display:flex;justify-content:space-between"><div class="card-name">'+esc(m.name)+'</div><span class="badge badge-'+m.status+'">'+m.status+'</span></div>';
if(m.role)h+='<div class="card-role">'+esc(m.role)+'</div>';
if(m.department)h+='<div class="card-dept">'+esc(m.department)+'</div>';
h+='<div class="card-contact">';if(m.email)h+='<a href="mailto:'+m.email+'">'+esc(m.email)+'</a>';if(m.phone)h+='<br>'+esc(m.phone);h+='</div>';
h+='<div style="margin-top:.5rem"><button class="btn" onclick="del(\''+m.id+'\')" style="font-size:.5rem;color:var(--cm)">Remove</button></div></div>';});
document.getElementById('grid').innerHTML=h;}
async function del(id){if(confirm('Remove?')){await fetch(A+'/members/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>Add Team Member</h2><div class="fr"><label>Name</label><input id="f-n"></div><div class="fr"><label>Email</label><input id="f-e" type="email"></div><div class="fr"><label>Role</label><input id="f-r" placeholder="e.g. Senior Engineer"></div><div class="fr"><label>Department</label><input id="f-d" placeholder="e.g. Engineering"></div><div class="fr"><label>Phone</label><input id="f-p"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Add</button></div>';document.getElementById('mbg').classList.add('open');}
async function sub(){await fetch(A+'/members',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-n').value,email:document.getElementById('f-e').value,role:document.getElementById('f-r').value,department:document.getElementById('f-d').value,phone:document.getElementById('f-p').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
