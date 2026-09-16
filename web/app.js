const API_KEY = 'circuitbreaker-secret-key';

async function api(path) {
  const res = await fetch(path, { headers: { 'X-Api-Key': API_KEY } });
  return res.json();
}

async function apiPost(path, body) {
  const res = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'X-Api-Key': API_KEY },
    body: JSON.stringify(body)
  });
  return res.json();
}

function renderStats(data) {
  const el = document.getElementById('stats');
  if (!data) return;
  const items = [
    { label: '下游服务', value: data.total_services },
    { label: '熔断器总数', value: data.total_breakers },
    { label: 'Open 状态', value: data.open_breakers },
    { label: 'HalfOpen 状态', value: data.half_open_breakers },
    { label: 'Closed 状态', value: data.closed_breakers },
    { label: '总调用量', value: data.total_calls },
    { label: '总失败数', value: data.total_failures },
    { label: '平均失败率', value: (data.avg_failure_ratio * 100).toFixed(2) + '%' },
    { label: '告警规则', value: data.total_alert_rules },
    { label: '已启用规则', value: data.enabled_alert_rules },
  ];
  el.innerHTML = items.map(i => `<div class="stat-card"><div class="label">${i.label}</div><div class="value">${i.value}</div></div>`).join('');
}

function renderServices(body) {
  const tbody = document.querySelector('#service-list tbody');
  tbody.innerHTML = '';
  const items = body.data?.items || [];
  items.forEach(s => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${s.id}</td><td>${s.name}</td><td>${s.address}</td><td>${s.protocol}</td><td>${s.status}</td><td>${s.weight}</td>`;
    tbody.appendChild(tr);
  });
}

async function transitionBreaker(id, state) {
  await apiPost(`/api/breakers/${id}/transition`, { target_state: state });
  loadBreakers();
  loadStats();
}

function renderBreakers(body) {
  const tbody = document.querySelector('#breaker-list tbody');
  tbody.innerHTML = '';
  const items = body.data?.items || [];
  items.forEach(b => {
    const tr = document.createElement('tr');
    let actions = '';
    if (b.state === 'closed') {
      actions += `<button class="btn-open" onclick="transitionBreaker('${b.id}', 'open')">熔断</button>`;
    } else if (b.state === 'open') {
      actions += `<button class="btn-half" onclick="transitionBreaker('${b.id}', 'half_open')">半开</button>`;
    } else if (b.state === 'half_open') {
      actions += `<button class="btn-closed" onclick="transitionBreaker('${b.id}', 'closed')">关闭</button>`;
      actions += `<button class="btn-open" onclick="transitionBreaker('${b.id}', 'open')">熔断</button>`;
    }
    tr.innerHTML = `<td>${b.id}</td><td>${b.service_id}</td><td>${b.rule_id}</td><td>${b.state}</td><td>${actions}</td>`;
    tbody.appendChild(tr);
  });
}

async function loadStats() {
  const body = await api('/api/stats/overview');
  renderStats(body.data);
}

async function loadServices() {
  const body = await api('/api/services');
  renderServices(body);
}

async function loadBreakers() {
  const body = await api('/api/breakers');
  renderBreakers(body);
}

async function load() {
  await loadStats();
  await loadServices();
  await loadBreakers();
}

load();
