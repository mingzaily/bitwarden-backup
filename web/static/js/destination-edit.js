// 删除备份目标
async function deleteDestination(id) {
    if (!confirm('确定要删除这个备份目标吗？')) return;

    try {
        const response = await fetch(`/api/destinations/${id}`, { method: 'DELETE' });
        if (response.ok) {
            loadDestinations();
            alert('删除成功');
        } else {
            alert('删除失败');
        }
    } catch (error) {
        console.error('Failed to delete destination:', error);
        alert('删除失败');
    }
}

// 编辑备份目标
async function editDestination(id) {
    try {
        const response = await fetch(`/api/destinations/${id}`);
        const dest = await response.json();

        document.getElementById('destinationId').value = dest.id;
        document.getElementById('destinationName').value = dest.name;
        document.getElementById('destinationType').value = dest.type;
        document.getElementById('destinationEnabled').checked = dest.enabled;

        toggleDestinationType();

        if (dest.type === 'local') {
            document.getElementById('localPath').value = dest.local_path || '';
        } else if (dest.type === 'webdav') {
            document.getElementById('webdavURL').value = dest.webdav_url || '';
            document.getElementById('webdavUsername').value = dest.webdav_username || '';
            document.getElementById('webdavPassword').value = dest.webdav_password || '';
            document.getElementById('webdavPath').value = dest.webdav_path || '';
        } else if (dest.type === 'server') {
            document.getElementById('targetServerId').value = dest.target_server_id || '';
        }

        showDestinationModal();
    } catch (error) {
        console.error('Failed to load destination:', error);
        alert('加载备份目标信息失败');
    }
}
