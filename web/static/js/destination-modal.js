// 显示备份目标模态框
async function showDestinationModal() {
    await loadServerOptionsForDestination();
    document.getElementById('destinationModal').classList.add('show');
    document.getElementById('destinationForm').reset();
    document.getElementById('destinationId').value = '';
    hideAllDestinationConfigs();
}

// 关闭备份目标模态框
function closeDestinationModal() {
    document.getElementById('destinationModal').classList.remove('show');
}

// 切换目标类型配置
function toggleDestinationType() {
    const type = document.getElementById('destinationType').value;
    hideAllDestinationConfigs();

    if (type === 'local') {
        document.getElementById('localConfig').style.display = 'block';
    } else if (type === 'webdav') {
        document.getElementById('webdavConfig').style.display = 'block';
    } else if (type === 'server') {
        document.getElementById('serverConfig').style.display = 'block';
    }
}

// 隐藏所有目标配置
function hideAllDestinationConfigs() {
    document.getElementById('localConfig').style.display = 'none';
    document.getElementById('webdavConfig').style.display = 'none';
    document.getElementById('serverConfig').style.display = 'none';
}
