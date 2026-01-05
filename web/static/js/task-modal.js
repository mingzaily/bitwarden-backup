// 显示任务模态框
async function showTaskModal() {
    await loadServerOptions();
    await loadDestinationCheckboxes();
    document.getElementById('taskModal').classList.add('show');
    document.getElementById('taskForm').reset();
    document.getElementById('taskId').value = '';
}

// 关闭任务模态框
function closeTaskModal() {
    document.getElementById('taskModal').classList.remove('show');
}

// 加载服务器选项
async function loadServerOptions() {
    try {
        const response = await fetch('/api/servers');
        const servers = await response.json();

        const sourceSelect = document.getElementById('sourceServerId');
        const options = servers.map(s =>
            `<option value="${s.id}">${s.name}</option>`
        ).join('');

        sourceSelect.innerHTML = '<option value="">请选择</option>' + options;
    } catch (error) {
        console.error('Failed to load servers:', error);
    }
}

// 加载备份目标复选框
async function loadDestinationCheckboxes() {
    try {
        const response = await fetch('/api/destinations');
        const destinations = await response.json();

        const container = document.getElementById('destinationCheckboxes');
        if (destinations.length === 0) {
            container.innerHTML = '<p style="color:#999;">暂无备份目标，请先添加备份目标</p>';
            return;
        }

        container.innerHTML = destinations.map(dest => `
            <label class="checkbox-label">
                <input type="checkbox" name="destinations" value="${dest.id}">
                ${dest.name} (${dest.type})
            </label>
        `).join('');
    } catch (error) {
        console.error('Failed to load destinations:', error);
    }
}
