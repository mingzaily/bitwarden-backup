// 显示服务器模态框
function showServerModal() {
    document.getElementById('serverModal').classList.add('show');
    document.getElementById('serverForm').reset();
    document.getElementById('serverId').value = '';
}

// 关闭服务器模态框
function closeServerModal() {
    document.getElementById('serverModal').classList.remove('show');
}

// 保存服务器配置
async function saveServer(event) {
    event.preventDefault();

    const id = document.getElementById('serverId').value;
    const data = {
        name: document.getElementById('serverName').value,
        server_url: document.getElementById('serverURL').value,
        client_id: document.getElementById('clientID').value,
        client_secret: document.getElementById('clientSecret').value,
        master_password: document.getElementById('masterPassword').value,
        is_official: document.getElementById('isOfficial').checked
    };

    try {
        const url = id ? `/api/servers/${id}` : '/api/servers';
        const method = id ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            closeServerModal();
            loadServers();
            alert('保存成功');
        } else {
            alert('保存失败');
        }
    } catch (error) {
        console.error('Failed to save server:', error);
        alert('保存失败');
    }
}
