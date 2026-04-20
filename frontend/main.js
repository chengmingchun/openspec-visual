document.addEventListener('DOMContentLoaded', async () => {
    const viewSettings = document.getElementById('view-settings');
    const viewWorkflow = document.getElementById('view-workflow');
    const viewFiles = document.getElementById('view-files');
    const viewTdd = document.getElementById('view-tdd');
    const viewHistory = document.getElementById('view-history');
    
    const navSettings = document.getElementById('nav-settings');
    const navWorkflow = document.getElementById('nav-workflow');
    const navFiles = document.getElementById('nav-files');
    const navTdd = document.getElementById('nav-tdd');
    const navHistory = document.getElementById('nav-history');

    const cfgUrl = document.getElementById('cfg-url');
    const cfgKey = document.getElementById('cfg-key');
    const cfgModel = document.getElementById('cfg-model');
    const cfgFeature = document.getElementById('cfg-feature');
    const btnSaveCfg = document.getElementById('btn-save-cfg');

    const s1 = document.getElementById('s1');
    const s2 = document.getElementById('s2');
    const s3 = document.getElementById('s3');
    const s4 = document.getElementById('s4');
    
    // Removed legacy button constants
    const btnViewFiles = document.getElementById('btn-view-files');
    const fileTreeContainer = document.getElementById('file-tree-container');
    const fileViewer = document.getElementById('file-viewer');

    // Review Panel Elements
    const pendingReviewPanel = document.getElementById('pending-review-panel');
    const pendingSkillBadge = document.getElementById('pending-skill-badge');
    const botChecks = document.getElementById('bot-checks');
    const reviewFeedback = document.getElementById('review-feedback');
    const btnApprove = document.getElementById('btn-approve');
    const btnReject = document.getElementById('btn-reject');

    // History Panel Elements
    const historyTimelineNav = document.getElementById('history-timeline-nav');
    const diffViewer = document.getElementById('diff-viewer');
    const btnRollback = document.getElementById('btn-rollback');
    let selectedCommitHash = "";

    let apiKeyStr = "";
    let hasConfig = false;

    try {
        const res = await fetch("http://127.0.0.1:38192/api/config");
        if (res.ok) {
            const config = await res.json();
            cfgUrl.value = config.baseUrl || "";
            cfgKey.value = config.apiKey || "";
            cfgModel.value = config.model || "";
            apiKeyStr = config.apiKey || "";
            if (config.apiKey) hasConfig = true;
        }
    } catch (e) {
        console.warn("REST backend not connected.");
    }

    function switchView(view) {
        [viewWorkflow, viewSettings, viewFiles, viewTdd, viewHistory].forEach(v => Math.abs(v && v.classList.add('hidden')));
        [navWorkflow, navSettings, navFiles, navTdd, navHistory].forEach(v => Math.abs(v && v.classList.remove('active')));

        if (view === 'settings') {
            viewSettings.classList.remove('hidden');
            navSettings.classList.add('active');
        } else if (view === 'files') {
            viewFiles.classList.remove('hidden');
            navFiles.classList.add('active');
            refreshFileTree();
        } else if (view === 'tdd') {
            viewTdd.classList.remove('hidden');
            navTdd.classList.add('active');
            navTdd.style.color = '';
            navTdd.style.animation = ''; // Cancel pulses
        } else if (view === 'history') {
            viewHistory.classList.remove('hidden');
            navHistory.classList.add('active');
            refreshHistory();
        } else {
            viewWorkflow.classList.remove('hidden');
            navWorkflow.classList.add('active');
        }
    }

    navSettings.onclick = () => switchView('settings');
    navWorkflow.onclick = () => switchView('workflow');
    navFiles.onclick = () => switchView('files');
    navTdd.onclick = () => switchView('tdd');
    navHistory.onclick = () => switchView('history');

    btnSaveCfg.onclick = async () => {
        try {
            const res = await fetch("http://127.0.0.1:38192/api/config", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ APIKey: cfgKey.value, BaseURL: cfgUrl.value, Model: cfgModel.value })
            });
            if (!res.ok) throw new Error("Failed to save config");
            apiKeyStr = cfgKey.value;
            hasConfig = !!cfgKey.value;
            alert('Settings Saved.');
            switchView('workflow');
        } catch (e) {
            alert('Save failed: ' + e);
        }
    };

    let lastReportCount = 0;

    async function pollAgentReports() {
        if (!viewWorkflow.classList.contains('hidden')) {
            try {
                const res = await fetch("http://127.0.0.1:38192/api/reports");
                if (res.ok) {
                    const reports = await res.json();
                    if (reports && reports.length > lastReportCount) {
                        lastReportCount = reports.length;
                        updateDashboardStages(reports);
                    }
                }
            } catch (e) {
                // Silently fail during poll
            }
        }
    }

    function updateDashboardStages(reports) {
        const skillToStage = {
            'propose': 's1',
            'validate': 's2',
            'apply': 's3',
            'archive': 's4'
        };

        const stages = ['s1', 's2', 's3', 's4'];
        
        let highestStageIndex = -1;
        reports.forEach(r => {
            const sid = skillToStage[r.skill_name];
            if (sid) {
                const idx = stages.indexOf(sid);
                if (idx > highestStageIndex) highestStageIndex = idx;
            }
        });

        for (let i = 0; i < stages.length; i++) {
            const el = document.getElementById(stages[i]);
            const logPanel = el.querySelector('.stage-logs');
            
            if (i < highestStageIndex) {
                el.className = 'stage completed';
                logPanel.classList.remove('hidden');
            } else if (i === highestStageIndex) {
                el.className = 'stage active';
                logPanel.classList.remove('hidden');
            } else {
                el.className = 'stage';
            }
        }

        stages.forEach(sid => {
            const el = document.getElementById(sid);
            const logPanel = el.querySelector('.stage-logs');
            const targetSkill = el.getAttribute('data-skill');
            const relatedReports = reports.filter(r => skillToStage[r.skill_name] === sid || (targetSkill === 'propose' && !skillToStage[r.skill_name])); 
            
            if (relatedReports.length > 0) {
                logPanel.innerHTML = relatedReports.map(r => `
                    <div class="log-entry">
                        <span class="log-time">${new Date().toLocaleTimeString()}</span>
                        <span class="log-status">[${r.status}]</span>
                        <span class="log-file">${r.file_path || ''}</span>
                    </div>
                `).join('');
            }
        });
    }

    // Git Fetch & Render
    window.loadDiff = async (hash) => {
        selectedCommitHash = hash;
        btnRollback.style.display = 'inline-block';
        diffViewer.innerHTML = '<div style="color:#888; text-align:center;">Fetching patch...</div>';
        try {
            const res = await fetch("http://127.0.0.1:38192/api/diff?hash=" + hash);
            const text = await res.text();
            if (!text) {
                diffViewer.innerHTML = '<div style="color:#888; text-align:center;">(Empty diff)</div>';
                return;
            }
            const html = escapeHtml(text).split('\\n').map(line => {
                if (line.startsWith('+')) return `<div style="color: green; background:#e6ffec;">${line}</div>`;
                if (line.startsWith('-')) return `<div style="color: red; background:#ffebe9;">${line}</div>`;
                if (line.startsWith('@@')) return `<div style="color: #0969da; background:#ddf4ff; margin: 4px 0;">${line}</div>`;
                return `<div>${line}</div>`;
            }).join('');
            diffViewer.innerHTML = html;
        } catch(e) {
            diffViewer.innerHTML = 'Error loading git diff.';
        }
    };

    async function refreshHistory() {
        try {
            const res = await fetch("http://127.0.0.1:38192/api/history");
            if (res.ok) {
                const logs = await res.json();
                if (logs && logs.length > 0) {
                    historyTimelineNav.innerHTML = logs.map((l, idx) => `
                        <div class="commit-card" onclick="loadDiff('${l.hash}')" style="border-left: 2px solid var(--accent); padding-left: 16px; margin-bottom: 24px; position: relative; cursor:pointer;">
                            <div style="position: absolute; left: -5px; top: 0; width: 8px; height: 8px; background: var(--accent); border-radius: 50%;"></div>
                            <div style="font-weight: 500; font-family: monospace; font-size: 13px; color: var(--text-main);">📌 版本快照：${l.hash.substring(0,7)}</div>
                            <div style="font-size: 12px; color: var(--accent); margin-top: 4px;">📝 ${escapeHtml(l.message)}</div>
                            <div style="font-size: 11px; color: #888; margin-top: 2px;">🕰️ ${new Date(l.date).toLocaleString()} by ${escapeHtml(l.author)}</div>
                        </div>
                    `).join('');
                } else {
                    historyTimelineNav.innerHTML = '<div style="color:#888;text-align:center;">暂无任何变更记录长廊数据。</div>';
                }
            }
        } catch(e) {}
    }

    btnRollback.onclick = async () => {
        if (!selectedCommitHash) return;
        if (confirm("🚨 危险隔离操作！\\n\\n确定要将项目目录规范强制抹除并撤回至节点 " + selectedCommitHash.substring(0,7) + " 吗？\\n此动作相当于时间逆流，会彻底删掉 Agent 在此之后的胡言乱语产物！")) {
            try {
                const res = await fetch("http://127.0.0.1:38192/api/rollback?hash=" + selectedCommitHash, {method: 'POST'});
                if (res.ok) {
                    alert('回滚成功！本地硬盘状态已强行重置！');
                    refreshHistory();
                    refreshFileTree();
                    selectedCommitHash = '';
                    btnRollback.style.display = 'none';
                    diffViewer.innerHTML = '<div style="color:#888; text-align:center;">回滚重置成功！工作区已降级！</div>';
                } else {
                    alert('回滚执行失败！');
                }
            } catch(e) {
                alert('回滚请求异常');
            }
        }
    };

    let isReviewPanelOpen = false;

    async function pollPendingReview() {
        if (!viewWorkflow.classList.contains('hidden') && !isReviewPanelOpen) {
            try {
                const res = await fetch("http://127.0.0.1:38192/api/pending");
                if (res.ok) {
                    const pending = await res.json();
                    if (pending && pending.request && pending.request.skill_name) {
                        isReviewPanelOpen = true;
                        showPendingReview(pending);
                    }
                }
            } catch (e) {
                // Silently fail
            }
        }
    }

    function showPendingReview(pending) {
        pendingSkillBadge.textContent = pending.request.skill_name;
        
        botChecks.innerHTML = pending.checker_results.map(ch => {
            const icon = ch.passed ? '✅' : '❌';
            const color = ch.passed ? 'var(--success)' : '#e00';
            return `<div style="color:${color}; margin-bottom:4px;"><span style="display:inline-block;width:20px;">${icon}</span> [<b>${ch.rule_name}</b>]: ${ch.message}</div>`;
        }).join('');

        reviewFeedback.value = '';
        document.getElementById('no-review-placeholder').classList.add('hidden');
        pendingReviewPanel.classList.remove('hidden');

        // Alert user visually on Sidebar
        if (viewTdd.classList.contains('hidden')) {
            navTdd.style.color = '#e00';
            navTdd.style.animation = 'pulse 1s infinite alternate';
        }
    }

    btnApprove.onclick = async () => {
        await submitReviewDecision(true, "通过");
    };

    btnReject.onclick = async () => {
        const feedback = reviewFeedback.value.trim() || "未提供反馈意见。";
        await submitReviewDecision(false, feedback);
    };

    async function submitReviewDecision(approved, feedback) {
        try {
            await fetch("http://127.0.0.1:38192/api/review", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ approved, feedback })
            });
            pendingReviewPanel.classList.add('hidden');
            document.getElementById('no-review-placeholder').classList.remove('hidden');
            navTdd.style.color = '';
            navTdd.style.animation = '';
            isReviewPanelOpen = false;
        } catch (e) {
            alert("提交审批决策失败：" + e);
        }
    }

    setInterval(() => {
        pollAgentReports();
        pollPendingReview();
    }, 1500);

    // File Tree Fetch & Render
    async function refreshFileTree() {
        try {
            const res = await fetch("http://127.0.0.1:38192/api/list");
            if (!res.ok) throw new Error("Backend connection failed.");
            const root = await res.json();
            if (!root) {
                fileTreeContainer.innerHTML = '<div style="padding:16px;font-size:13px;color:#888;">No openspec/ folder generated yet. Go to Workflow and run a generation phase.</div>';
                return;
            }
            fileTreeContainer.innerHTML = '';
            const frag = renderNode(root);
            fileTreeContainer.appendChild(frag);
        } catch(e) {
            fileTreeContainer.innerHTML = 'Error listing files: ' + e;
        }
    }

    function renderNode(node) {
        const div = document.createElement('div');
        div.className = 'tree-node';
        
        const item = document.createElement('div');
        item.className = 'tree-item';
        item.textContent = (node.isDir ? '📁 ' : '📄 ') + node.name;
        
        div.appendChild(item);

        if (node.isDir && node.children) {
            const childrenContainer = document.createElement('div');
            childrenContainer.className = 'tree-children';
            node.children.forEach(c => childrenContainer.appendChild(renderNode(c)));
            div.appendChild(childrenContainer);
        } else if (!node.isDir) {
            item.onclick = async (e) => {
                e.stopPropagation();
                try {
                    const res = await fetch("http://127.0.0.1:38192/api/read?path=" + encodeURIComponent(node.path));
                    const content = await res.text();
                    if (node.path.endsWith('.md') || node.path.endsWith('.mdx')) {
                        fileViewer.innerHTML = `<div class="markdown-body" style="text-align:left;">${window.marked.parse(content)}</div>`;
                    } else {
                        fileViewer.innerHTML = `<h3 style="margin-bottom:16px">${node.name}</h3><pre>${escapeHtml(content)}</pre>`;
                    }
                } catch(error) {
                    fileViewer.innerHTML = 'Failed to load file.';
                }
            };
        }

        return div;
    }

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
});
