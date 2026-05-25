const API = '';
const addForm = document.getElementById('add-form');
const addMsg = document.getElementById('add-message');
const reviewArea = document.getElementById('review-word');
const rWord = document.getElementById('r-word');
const rTrans = document.getElementById('r-translation');
const rEx = document.getElementById('r-example');
const wordsTbody = document.getElementById('words-tbody');
let currentId = null;

addForm.addEventListener('submit', async e => {
    e.preventDefault();
    const data = {
        word: document.getElementById('word').value,
        translation: document.getElementById('translation').value,
        example: document.getElementById('example').value,
        difficulty: +document.getElementById('difficulty').value
    };
    try {
        const res = await fetch(`${API}/api/words`, { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(data) });
        if(res.ok) { addMsg.textContent='✅ Добавлено!'; addMsg.className='success'; addForm.reset(); loadWords(); }
        else throw new Error('Ошибка');
    } catch(err) { addMsg.textContent='❌ '+err.message; addMsg.className='error'; }
    setTimeout(()=> addMsg.textContent='', 3000);
});

document.getElementById('get-review-btn').addEventListener('click', async () => {
    try {
        const res = await fetch(`${API}/api/words/review`);
        if(res.status===204) { alert('Нет слов для повторения'); reviewArea.classList.add('hidden'); return; }
        if(!res.ok) throw new Error('Ошибка');
        const w = await res.json();
        currentId = w.id;
        rWord.textContent=w.word; rTrans.textContent=w.translation; rEx.textContent=w.example||'—';
        reviewArea.classList.remove('hidden');
    } catch(err) { alert(err.message); }
});

async function updateProgress(recalled) {
    if(!currentId) return;
    try {
        const res = await fetch(`${API}/api/words/progress?id=${currentId}`, { method:'PUT', headers:{'Content-Type':'application/json'}, body:JSON.stringify({recalled}) });
        if(res.ok) { alert(recalled?'Отлично!':'Повторим позже.'); reviewArea.classList.add('hidden'); currentId=null; loadWords(); }
    } catch(err) { alert(err.message); }
}
document.getElementById('recalled-btn').addEventListener('click', ()=> updateProgress(true));
document.getElementById('not-recalled-btn').addEventListener('click', ()=> updateProgress(false));

async function loadWords() {
    try {
        const res = await fetch(`${API}/api/words`);
        const words = await res.json();
        wordsTbody.innerHTML = words.map(w => `<tr><td>${w.id}</td><td>${w.word}</td><td>${w.translation}</td><td>${w.difficulty}</td><td>${w.status}</td><td>${new Date(w.next_review).toLocaleString()}</td></tr>`).join('');
    } catch(e) { console.error(e); }
}
document.getElementById('refresh-words-btn').addEventListener('click', loadWords);
loadWords();