package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Roster</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1100px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.green{color:var(--green)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:.6rem}
.card{background:var(--bg2);border:1px solid var(--bg3);padding:.9rem 1rem;cursor:pointer;transition:border-color .15s;display:flex;flex-direction:column;gap:.3rem}
.card:hover{border-color:var(--leather)}
.card.inactive{opacity:.55}
.card-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.card-name{font-size:.85rem;font-weight:700;color:var(--cream)}
.card-role{font-size:.62rem;color:var(--cd);margin-top:.1rem}
.card-meta{font-size:.55rem;color:var(--cm);display:flex;gap:.5rem;flex-wrap:wrap;margin-top:.3rem;align-items:center}
.card-meta a{color:var(--cd);text-decoration:none}
.card-meta a:hover{color:var(--rust)}
.card-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.card-extra-row{display:flex;gap:.4rem}
.card-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.card-extra-val{color:var(--cream)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.active{border-color:var(--green);color:var(--green)}
.badge.inactive{border-color:var(--cm);color:var(--cm)}
.badge.on_leave{border-color:var(--orange);color:var(--orange)}
.badge.dept{border-color:var(--leather);color:var(--leather)}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto;color:var(--red);border-color:#3a1a1a}
.acts .btn-del:hover{border-color:var(--red);color:var(--red)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}.trial-bar{flex-direction:column;align-items:stretch}.trial-bar input.key-input{width:100%}}
.trial-bar{display:none;background:linear-gradient(90deg,#3a2419,#2e1c14);border-bottom:2px solid var(--rust);padding:.7rem 1.5rem;font-family:var(--mono);font-size:.68rem;color:var(--cream);align-items:center;gap:1rem;flex-wrap:wrap}
.trial-bar.show{display:flex}
.trial-bar-msg{flex:1;min-width:240px;line-height:1.5}
.trial-bar-msg strong{color:var(--rust);text-transform:uppercase;letter-spacing:1px;font-size:.6rem;display:block;margin-bottom:.15rem}
.trial-bar-actions{display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
.trial-bar a.btn-trial{background:var(--rust);color:#fff;padding:.4rem .8rem;text-decoration:none;font-size:.65rem;text-transform:uppercase;letter-spacing:1px;font-weight:700;border:1px solid var(--rust);transition:all .2s}
.trial-bar a.btn-trial:hover{background:#f08545;border-color:#f08545}
.trial-bar-divider{color:var(--cm);font-size:.6rem}
.trial-bar input.key-input{padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.6rem;width:200px}
.trial-bar input.key-input:focus{outline:none;border-color:var(--rust)}
.trial-bar button.btn-activate{padding:.4rem .7rem;background:var(--bg2);color:var(--cream);border:1px solid var(--leather);font-family:var(--mono);font-size:.6rem;cursor:pointer;text-transform:uppercase;letter-spacing:1px}
.trial-bar button.btn-activate:hover{background:var(--bg3)}
.trial-bar button.btn-activate:disabled{opacity:.5;cursor:wait}
.trial-msg{font-size:.6rem;color:var(--cm);margin-left:.5rem}
.trial-msg.error{color:#e74c3c}
.trial-msg.success{color:#4ade80}
.btn-disabled-trial{opacity:.45;cursor:not-allowed!important}
</style>
</head>
<body>

<div class="trial-bar" id="trial-bar">
<div class="trial-bar-msg">
<strong>Trial Required</strong>
You can view your existing members, but adding, editing, or removing is locked until you start a 14-day free trial.
</div>
<div class="trial-bar-actions">
<a class="btn-trial" href="https://stockyard.dev/" target="_blank" rel="noopener">Start 14-Day Trial</a>
<span class="trial-bar-divider">or</span>
<input type="text" class="key-input" id="trial-key-input" placeholder="SY-..." autocomplete="off" spellcheck="false">
<button class="btn-activate" id="trial-activate-btn" onclick="activateLicense()">Activate</button>
<span class="trial-msg" id="trial-msg"></span>
</div>
</div>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> ROSTER</h1>
<button class="btn btn-p" onclick="openNew()">+ Add Member</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search name, email, role..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="active">Active</option>
<option value="on_leave">On Leave</option>
<option value="inactive">Inactive</option>
</select>
<select class="filter-sel" id="department-filter" onchange="render()">
<option value="">All Departments</option>
</select>
</div>
<div id="grid" class="grid"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='members';

var fields=[
{name:'name',label:'Name',type:'text',required:true},
{name:'email',label:'Email',type:'email'},
{name:'phone',label:'Phone',type:'text'},
{name:'role',label:'Role / Title',type:'text'},
{name:'department',label:'Department',type:'select_or_text',options:[]},
{name:'status',label:'Status',type:'select',options:['active','on_leave','inactive']},
{name:'join_date',label:'Join Date',type:'date'}
];

var members=[],memberExtras={},editId=null,searchTimer=null;

// ─── Helpers ──────────────────────────────────────────────────────

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s+'T12:00:00');
if(isNaN(d.getTime()))return s;
return d.toLocaleDateString('en-US',{year:'numeric',month:'short',day:'numeric'});
}catch(e){return s}
}

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/members').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
members=resps[0].members||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
memberExtras=ex||{};
members.forEach(function(m){
var x=memberExtras[m.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(m[k]===undefined)m[k]=x[k]});
});
}catch(e){memberExtras={}}

populateDepartmentFilter();
}catch(e){
console.error('load failed',e);
members=[];
}
render();
}

function populateDepartmentFilter(){
var sel=document.getElementById('department-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var depts=[];
members.forEach(function(m){
if(m.department&&!seen[m.department]){seen[m.department]=true;depts.push(m.department)}
});
depts.sort();
sel.innerHTML='<option value="">All Departments</option>'+depts.map(function(d){return'<option value="'+esc(d)+'"'+(d===current?' selected':'')+'>'+esc(d)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var active=s.active||0;
var byDept=s.by_department||{};
var deptCount=Object.keys(byDept).length;
var byStatus=s.by_status||{};
var onLeave=byStatus['on_leave']||0;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Members</div></div>'+
'<div class="st"><div class="st-v green">'+active+'</div><div class="st-l">Active</div></div>'+
'<div class="st"><div class="st-v">'+onLeave+'</div><div class="st-l">On Leave</div></div>'+
'<div class="st"><div class="st-v">'+deptCount+'</div><div class="st-l">Departments</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var df=document.getElementById('department-filter').value;

var f=members;
if(q)f=f.filter(function(m){
return(m.name||'').toLowerCase().includes(q)||
(m.email||'').toLowerCase().includes(q)||
(m.role||'').toLowerCase().includes(q)||
(m.department||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(m){return m.status===sf});
if(df)f=f.filter(function(m){return m.department===df});

if(!f.length){
var msg=window._emptyMsg||'No members yet. Add your first one.';
document.getElementById('grid').innerHTML='<div class="empty" style="grid-column:1/-1">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(m){h+=cardHTML(m)});
document.getElementById('grid').innerHTML=h;
}

function cardHTML(m){
var cls='card';
if(m.status&&m.status!=='active')cls+=' inactive';

var clickHandler=window._trialRequired?'showTrialNudge()':'openEdit(\\''+esc(m.id)+'\\')';
var h='<div class="'+cls+'" onclick="'+clickHandler+'">';
h+='<div class="card-top"><div style="flex:1;min-width:0">';
h+='<div class="card-name">'+esc(m.name)+'</div>';
if(m.role)h+='<div class="card-role">'+esc(m.role)+'</div>';
h+='</div></div>';

h+='<div class="card-meta">';
if(m.status)h+='<span class="badge '+esc(m.status)+'">'+esc(m.status.replace(/_/g,' '))+'</span>';
if(m.department)h+='<span class="badge dept">'+esc(m.department)+'</span>';
h+='</div>';

if(m.email||m.phone){
h+='<div class="card-meta">';
if(m.email)h+='<a href="mailto:'+esc(m.email)+'" onclick="event.stopPropagation()">'+esc(m.email)+'</a>';
if(m.phone)h+='<span>'+esc(m.phone)+'</span>';
h+='</div>';
}
if(m.join_date)h+='<div class="card-meta"><span>Joined '+esc(fmtDate(m.join_date))+'</span></div>';

// Custom field display
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=m[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="card-extra-row">';
customRows+='<span class="card-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="card-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="card-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=String(o).charAt(0).toUpperCase()+String(o).slice(1).replace(/_/g,' ');
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
var opts=(f.options||[]).slice();
members.forEach(function(md){
if(md.department&&opts.indexOf(md.department)===-1)opts.push(md.department);
});
opts.forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='number'||f.type==='integer'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(member){
var m=member||{};
var isEdit=!!member;
var h='<h2>'+(isEdit?'EDIT MEMBER':'NEW MEMBER')+'</h2>';

h+=fieldHTML(fieldByName('name'),m.name);
h+='<div class="row2">'+fieldHTML(fieldByName('email'),m.email)+fieldHTML(fieldByName('phone'),m.phone)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('role'),m.role)+fieldHTML(fieldByName('department'),m.department)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('status'),m.status||'active')+fieldHTML(fieldByName('join_date'),m.join_date)+'</div>';

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,m[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit){
h+='<button class="btn btn-del" onclick="delMember()">Delete</button>';
}
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button>';
h+='</div>';
return h;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var m=null;
for(var i=0;i<members.length;i++){if(members[i].id===id){m=members[i];break}}
if(!m)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(m);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='number'||f.type==='integer')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/members/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/members',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Add failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delMember(){
if(!editId)return;
if(!confirm('Delete this member?'))return;
await fetch(A+'/members/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.departments)){
var deptField=fieldByName('department');
if(deptField)deptField.options=cfg.departments;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
checkTrialState();
load();
});
})();

// ─── trial-required license gating ───
window._trialRequired=false;

async function checkTrialState(){
try{
var resp=await fetch('/api/tier');
if(!resp.ok)return;
var data=await resp.json();
window._trialRequired=!!data.trial_required;
if(window._trialRequired){
document.getElementById('trial-bar').classList.add('show');
disableWriteControls();
// Re-render so card onclick handlers pick up the new trial state
if(typeof render==='function')render();
}else{
document.getElementById('trial-bar').classList.remove('show');
}
}catch(e){}
}

function disableWriteControls(){
var buttons=document.querySelectorAll('.hdr .btn, .hdr .btn-p');
buttons.forEach(function(b){
var t=b.textContent||'';
if(t.indexOf('Add')!==-1||t.indexOf('New')!==-1){
b.classList.add('btn-disabled-trial');
b.title='Locked: trial required';
b.onclick=function(e){
e.preventDefault();
showTrialNudge();
return false;
};
}
});
}

function showTrialNudge(){
var input=document.getElementById('trial-key-input');
if(input){
input.focus();
input.style.borderColor='var(--rust)';
setTimeout(function(){if(input)input.style.borderColor=''},1500);
}
}

async function activateLicense(){
var input=document.getElementById('trial-key-input');
var btn=document.getElementById('trial-activate-btn');
var msg=document.getElementById('trial-msg');
if(!input||!btn||!msg)return;
var key=(input.value||'').trim();
if(!key){
msg.className='trial-msg error';
msg.textContent='Paste your license key first';
input.focus();
return;
}
btn.disabled=true;
msg.className='trial-msg';
msg.textContent='Activating...';
try{
var resp=await fetch('/api/license/activate',{
method:'POST',
headers:{'Content-Type':'application/json'},
body:JSON.stringify({license_key:key})
});
var data=await resp.json();
if(!resp.ok){
msg.className='trial-msg error';
msg.textContent=data.error||'Activation failed';
btn.disabled=false;
return;
}
msg.className='trial-msg success';
msg.textContent='Activated. Reloading...';
setTimeout(function(){location.reload()},800);
}catch(e){
msg.className='trial-msg error';
msg.textContent='Network error: '+e.message;
btn.disabled=false;
}
}

document.addEventListener('DOMContentLoaded',function(){
var input=document.getElementById('trial-key-input');
if(input){
input.addEventListener('keydown',function(e){
if(e.key==='Enter')activateLicense();
});
}
});
</script>
</body>
</html>`
