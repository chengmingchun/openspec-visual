document.addEventListener('DOMContentLoaded', async () => {
    const viewSettings = document.getElementById('view-settings');
    const viewWorkflow = document.getElementById('view-workflow');
    const viewFiles = document.getElementById('view-files');
    
    const navSettings = document.getElementById('nav-settings');
    const navWorkflow = document.getElementById('nav-workflow');
    const navFiles = document.getElementById('nav-files');

    const cfgUrl = document.getElementById('cfg-url');
    const cfgKey = document.getElementById('cfg-key');
    const cfgModel = document.getElementById('cfg-model');
    const cfgFeature = document.getElementById('cfg-feature');
    const btnSaveCfg = document.getElementById('btn-save-cfg');

    const s1 = document.getElementById('s1');
    const s2 = document.getElementById('s2');
    const s3 = document.getElementById('s3');
    const s4 = document.getElementById('s4');
    
    const s2act = document.getElementById('s2-action');
    const s3act = document.getElementById('s3-action');
    const s4act = document.getElementById('s4-action');
    const s1act = document.getElementById('s1-action');

    const promptInput = document.getElementById('prompt-input');
    const btnStart = document.getElementById('btn-start');
    const btnViewFiles = document.getElementById('btn-view-files');
    const fileTreeContainer = document.getElementById('file-tree-container');
    const fileViewer = document.getElementById('file-viewer');

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
        [viewWorkflow, viewSettings, viewFiles].forEach(v => v.classList.add('hidden'));
        [navWorkflow, navSettings, navFiles].forEach(v => v.classList.remove('active'));

        if (view === 'settings') {
            viewSettings.classList.remove('hidden');
            navSettings.classList.add('active');
        } else if (view === 'files') {
            viewFiles.classList.remove('hidden');
            navFiles.classList.add('active');
            refreshFileTree();
        } else {
            viewWorkflow.classList.remove('hidden');
            navWorkflow.classList.add('active');
        }
    }

    navSettings.onclick = () => switchView('settings');
    navWorkflow.onclick = () => switchView('workflow');
    navFiles.onclick = () => switchView('files');
    btnViewFiles.onclick = () => switchView('files');

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

    btnStart.onclick = async () => {
        const prompt = promptInput.value.trim();
        if (!prompt) {
            alert("Please input a feature request.");
            return;
        }

        const featureName = cfgFeature.value.trim() || "new-feature";
        
        // Stage 1 -> 2
        s1.classList.remove('active');
        s1.classList.add('completed');
        s1act.classList.add('hidden');
        
        s2.classList.add('active');
        s2act.classList.remove('hidden');

        // AI Request Simulation / Real Wait
        let aiResult = "MOCK_MODE";
        if (apiKeyStr) {
            try {
                const res = await fetch("http://127.0.0.1:38192/api/prompt", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ Prompt: prompt, System: "You are OpenSpec analyst." })
                });
                if (!res.ok) {
                    throw new Error(await res.text());
                }
                const data = await res.json();
                aiResult = data.result;
            } catch(e) {
                console.error(e);
                alert("大模型请求失败: " + e);
            }
        } else {
            await new Promise(r => setTimeout(r, 4000));
        }

        // Stage 2 -> 3
        s2.classList.remove('active');
        s2.classList.add('completed');
        s2act.classList.add('hidden');

        s3.classList.add('active');
        s3act.classList.remove('hidden');

        // Task & File generation
        await new Promise(r => setTimeout(r, 4000)); // Simulating writing longer
        try {
            const res = await fetch("http://127.0.0.1:38192/api/generate", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ FeatureName: featureName, Content: prompt })
            });
            if (!res.ok) {
                alert("写入核心文件警告: " + await res.text());
            }
        } catch(e) {
            console.error("Write error:", e);
            alert("文件落盘通讯失败: " + e);
        }

        // Stage 3 -> 4
        s3.classList.remove('active');
        s3.classList.add('completed');
        s3act.classList.add('hidden');

        s4.classList.add('active');
        s4act.classList.remove('hidden');
    };

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
                    fileViewer.innerHTML = `<h3 style="margin-bottom:16px">${node.name}</h3><pre>${escapeHtml(content)}</pre>`;
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
