// 加载日志列表
async function loadLogs() {
    try {
        const response = await fetch('/api/logs');
        const logs = await response.json();

        const list = document.getElementById('logs-list');
        if (logs.length === 0) {
            list.innerHTML = '<p class="empty-state">暂无执行日志</p>';
            return;
        }

        list.innerHTML = `
            <table class="log-table">
                <thead>
                    <tr>
                        <th>时间</th>
                        <th>任务ID</th>
                        <th>状态</th>
                        <th>消息</th>
                        <th>备份文件</th>
                    </tr>
                </thead>
                <tbody>
                    ${logs.map(log => `
                        <tr class="log-${log.status}">
                            <td>${new Date(log.start_time).toLocaleString()}</td>
                            <td>${log.task_id}</td>
                            <td><span class="status-badge status-${log.status}">${log.status}</span></td>
                            <td>${log.message || '-'}</td>
                            <td>${log.backup_file || '-'}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (error) {
        console.error('Failed to load logs:', error);
        alert('加载日志失败');
    }
}
