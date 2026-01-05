// 编辑服务器
async function editServer(id) {
    try {
        const response = await fetch(`/api/servers/${id}`);
        const server = await response.json();

        document.getElementById('serverId').value = server.id;
        document.getElementById('serverName').value = server.name;
        document.getElementById('serverURL').value = server.server_url;
        document.getElementById('clientID').value = server.client_id;
        document.getElementById('clientSecret').value = server.client_secret;
        document.getElementById('masterPassword').value = server.master_password;
        document.getElementById('isOfficial').checked = server.is_official;

        showServerModal();
    } catch (error) {
        console.error('Failed to load server:', error);
        alert('加载服务器信息失败');
    }
}

// 删除服务器
async function deleteServer(id) {
    if (!confirm('确定要删除这个服务器配置吗？')) return;

    try {
        const response = await fetch(`/api/servers/${id}`, { method: 'DELETE' });
        if (response.ok) {
            loadServers();
            alert('删除成功');
        } else {
            alert('删除失败');
        }
    } catch (error) {
        console.error('Failed to delete server:', error);
        alert('删除失败');
    }
}
