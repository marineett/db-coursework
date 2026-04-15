const FETCH_INTERVAL = 2000;
const HISTORY_LENGTH = 60;

const ids = {
  active: document.getElementById('active'),
  accepted: document.getElementById('accepted'),
  handled: document.getElementById('handled'),
  requests: document.getElementById('requests'),
  reading: document.getElementById('reading'),
  writing: document.getElementById('writing'),
  waiting: document.getElementById('waiting'),
};

let history = Array(HISTORY_LENGTH).fill(0);

async function fetchRaw() {
  try {
    const res = await fetch('/status/raw', {cache: 'no-store'});
    if (!res.ok) throw new Error('HTTP ' + res.status);
    const text = await res.text();
    return text;
  } catch (e) {
    console.error('Fetch status failed', e);
    return null;
  }
}

function parseStubStatus(text) {
  // Example nginx stub_status output:
  // Active connections: 1
  // server accepts handled requests
  //  123 123 456
  // Reading: 0 Writing: 1 Waiting: 0
  const out = {};
  if (!text) return out;
  const lines = text.trim().split('\n').map(l => l.trim());
  for (const line of lines) {
    if (line.startsWith('Active connections:')) {
      out.active = parseInt(line.split(':')[1]) || 0;
    } else if (/^\d+\s+\d+\s+\d+$/.test(line)) {
      const parts = line.split(/\s+/).map(s => parseInt(s));
      out.accepted = parts[0] || 0;
      out.handled = parts[1] || 0;
      out.requests = parts[2] || 0;
    } else if (line.startsWith('Reading:')) {
      // Reading: 0 Writing: 1 Waiting: 0
      const m = line.match(/Reading:\s*(\d+)\s+Writing:\s*(\d+)\s+Waiting:\s*(\d+)/i);
      if (m) {
        out.reading = parseInt(m[1]) || 0;
        out.writing = parseInt(m[2]) || 0;
        out.waiting = parseInt(m[3]) || 0;
      }
    }
  }
  return out;
}

function updateDOM(values) {
  if (!values) return;
  ids.active && (ids.active.textContent = values.active ?? '—');
  ids.accepted && (ids.accepted.textContent = values.accepted ?? '—');
  ids.handled && (ids.handled.textContent = values.handled ?? '—');
  ids.requests && (ids.requests.textContent = values.requests ?? '—');
  ids.reading && (ids.reading.textContent = values.reading ?? '—');
  ids.writing && (ids.writing.textContent = values.writing ?? '—');
  ids.waiting && (ids.waiting.textContent = values.waiting ?? '—');
}

function pushHistory(value) {
  history.push(value);
  if (history.length > HISTORY_LENGTH) history.shift();
}

function drawChart() {
  const canvas = document.getElementById('requestsChart');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  const w = canvas.width = canvas.clientWidth * devicePixelRatio;
  const h = canvas.height = 160 * devicePixelRatio;
  ctx.clearRect(0,0,w,h);

  const max = Math.max(...history, 1);

  // background grid
  ctx.save();
  ctx.scale(devicePixelRatio, devicePixelRatio);
  ctx.fillStyle = 'rgba(255,255,255,0.02)';
  ctx.fillRect(0,0,canvas.clientWidth,160);
  ctx.restore();

  // draw line
  ctx.beginPath();
  for (let i = 0; i < history.length; i++) {
    const x = (i / (history.length - 1)) * (canvas.clientWidth - 12) + 6;
    const y = 20 + (120 * (1 - history[i] / max));
    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  const grad = ctx.createLinearGradient(0,0,0,160);
  grad.addColorStop(0, 'rgba(96,165,250,0.95)');
  grad.addColorStop(1, 'rgba(124,58,237,0.6)');
  ctx.lineWidth = 2 * devicePixelRatio;
  ctx.strokeStyle = grad;
  ctx.stroke();

  // fill under line
  ctx.lineTo(canvas.clientWidth - 6, 160);
  ctx.lineTo(6,160);
  ctx.closePath();
  ctx.fillStyle = 'rgba(96,165,250,0.08)';
  ctx.fill();

  // current value
  const latest = history[history.length - 1] || 0;
  ctx.font = `${14 * devicePixelRatio}px sans-serif`;
  ctx.fillStyle = '#cfe7ff';
  ctx.fillText(`now: ${latest}`, 8 * devicePixelRatio, 16 * devicePixelRatio);
}

let lastRequests = null;
async function tick() {
  const text = await fetchRaw();
  if (!text) return;
  const vals = parseStubStatus(text);
  // compute requests per interval if possible
  if (typeof vals.requests === 'number') {
    if (lastRequests === null) lastRequests = vals.requests;
    const delta = Math.max(0, vals.requests - lastRequests);
    lastRequests = vals.requests;
    pushHistory(delta);
  } else {
    pushHistory(0);
  }
  updateDOM(vals);
  drawChart();
}

// init
(function(){
  // try to fill initial history
  for (let i=0;i<HISTORY_LENGTH;i++) history[i]=0;
  tick();
  setInterval(tick, FETCH_INTERVAL);
})();
