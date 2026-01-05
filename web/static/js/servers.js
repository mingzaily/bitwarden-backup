// 加载服务器列表
async function loadServers() {
    try {
        const response = await fetch('/api/servers');
        const servers = await response.json();

        const list = document.getElementById('servers-list');
        if (servers.length === 0) {
            list.innerHTML = '<p class="empty-state">暂无服务器配置</p>';
            return;
        }

        list.innerHTML = servers.map(server => `
            <div class="card">
                <h3>${server.name}</h3>
                <p><strong>URL:</strong> ${server.server_url}</p>
                <p><strong>类型:</strong> ${server.is_official ? '官方服务器' : '自建服务器'}</p>
                <div class="card-actions">
                    <button class="btn btn-secondary" onclick="editServer(${server.id})">编辑</button>
                    <button class="btn btn-danger" onclick="deleteServer(${server.id})">删除</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load servers:', error);
        alert('加载服务器列表失败');
    }
}
