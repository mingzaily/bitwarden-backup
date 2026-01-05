// 加载服务器选项（用于目标服务器类型）
async function loadServerOptionsForDestination() {
    try {
        const response = await fetch('/api/servers');
        const servers = await response.json();

        const select = document.getElementById('targetServerId');
        const options = servers.map(s =>
            `<option value="${s.id}">${s.name}</option>`
        ).join('');

        select.innerHTML = '<option value="">请选择</option>' + options;
    } catch (error) {
        console.error('Failed to load servers:', error);
    }
}

// 保存备份目标
async function saveDestination(event) {
    event.preventDefault();

    const id = document.getElementById('destinationId').value;
    const type = document.getElementById('destinationType').value;

    const data = {
        name: document.getElementById('destinationName').value,
        type: type,
        enabled: document.getElementById('destinationEnabled').checked
    };

    // 根据类型添加配置
    if (type === 'local') {
        data.local_path = document.getElementById('localPath').value;
    } else if (type === 'webdav') {
        data.webdav_url = document.getElementById('webdavURL').value;
        data.webdav_username = document.getElementById('webdavUsername').value;
        data.webdav_password = document.getElementById('webdavPassword').value;
        data.webdav_path = document.getElementById('webdavPath').value;
    } else if (type === 'server') {
        const serverId = document.getElementById('targetServerId').value;
        data.target_server_id = serverId ? parseInt(serverId) : null;
    }

    try {
        const url = id ? `/api/destinations/${id}` : '/api/destinations';
        const method = id ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            closeDestinationModal();
            loadDestinations();
            alert('保存成功');
        } else {
            alert('保存失败');
        }
    } catch (error) {
        console.error('Failed to save destination:', error);
        alert('保存失败');
    }
}
