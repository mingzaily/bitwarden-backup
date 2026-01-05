// 加载备份目标列表
async function loadDestinations() {
    try {
        const response = await fetch('/api/destinations');
        const destinations = await response.json();

        const list = document.getElementById('destinations-list');
        if (destinations.length === 0) {
            list.innerHTML = '<p class="empty-state">暂无备份目标配置</p>';
            return;
        }

        list.innerHTML = destinations.map(dest => {
            let typeLabel = '';
            let configInfo = '';

            if (dest.type === 'local') {
                typeLabel = '本地存储';
                configInfo = `<p><strong>路径:</strong> ${dest.local_path}</p>`;
            } else if (dest.type === 'webdav') {
                typeLabel = 'WebDAV';
                configInfo = `<p><strong>URL:</strong> ${dest.webdav_url}</p>
                             <p><strong>路径:</strong> ${dest.webdav_path}</p>`;
            } else if (dest.type === 'server') {
                typeLabel = '目标服务器';
                configInfo = `<p><strong>服务器ID:</strong> ${dest.target_server_id}</p>`;
            }

            return `
                <div class="card">
                    <h3>${dest.name}</h3>
                    <p><strong>类型:</strong> ${typeLabel}</p>
                    ${configInfo}
                    <p><strong>状态:</strong> ${dest.enabled ? '✅ 启用' : '❌ 禁用'}</p>
                    <div class="card-actions">
                        <button class="btn btn-secondary" onclick="editDestination(${dest.id})">编辑</button>
                        <button class="btn btn-danger" onclick="deleteDestination(${dest.id})">删除</button>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Failed to load destinations:', error);
        alert('加载备份目标列表失败');
    }
}
