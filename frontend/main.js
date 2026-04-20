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

        // Update History Timeline View
        const historyTimeline = document.getElementById('history-timeline');
        if (reports.length > 0) {
            historyTimeline.innerHTML = reports.map((r, idx) => `
                <div style="border-left: 2px solid var(--accent); padding-left: 16px; margin-bottom: 24px; position: relative;">
                    <div style="position: absolute; left: -5px; top: 0; width: 8px; height: 8px; background: var(--accent); border-radius: 50%;"></div>
                    <div style="font-weight: 500; font-family: monospace; font-size: 13px; color: var(--text-main);">迭代拦截位 #${idx+1}： 触发技能 <b>[${r.skill_name}]</b></div>
                    <div style="font-size: 12px; color: var(--text-muted); margin-top: 4px;">🎯 拦截节点：执行状态 [${r.status}]；相关文件：<code>${r.file_path || 'Unknown'}</code></div>
                </div>
            `).reverse().join('');
        }
    }

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
