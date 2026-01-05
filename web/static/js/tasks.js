// 加载任务列表
async function loadTasks() {
    try {
        const response = await fetch('/api/tasks');
        const tasks = await response.json();

        const list = document.getElementById('tasks-list');
        if (tasks.length === 0) {
            list.innerHTML = '<p class="empty-state">暂无备份任务</p>';
            return;
        }

        list.innerHTML = tasks.map(task => {
            const destinations = task.destinations || [];
            const destInfo = destinations.length > 0
                ? destinations.map(d => `${d.name} (${d.type})`).join(', ')
                : '未配置备份目标';

            return `
                <div class="card">
                    <h3>${task.name}</h3>
                    <p><strong>Cron:</strong> ${task.cron_expression}</p>
                    <p><strong>备份目标:</strong> ${destInfo}</p>
                    <p><strong>状态:</strong> ${task.enabled ? '✅ 启用' : '❌ 禁用'}</p>
                    <div class="card-actions">
                        <button class="btn btn-success" onclick="executeTask(${task.id})">立即执行</button>
                        <button class="btn btn-secondary" onclick="editTask(${task.id})">编辑</button>
                        <button class="btn btn-danger" onclick="deleteTask(${task.id})">删除</button>
                    </div>
                </div>
            `;
        }).join('');
    } catch (error) {
        console.error('Failed to load tasks:', error);
        alert('加载任务列表失败');
    }
}

// 立即执行任务
async function executeTask(id) {
    if (!confirm('确定要立即执行此备份任务吗？')) {
        return;
    }

    try {
        const response = await fetch(`/api/tasks/${id}/execute`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            const data = await response.json();
            alert('任务已开始执行，请查看执行日志查看进度');
            // 切换到日志标签页
            showTab('logs');
            loadLogs();
        } else {
            alert('执行失败，请重试');
        }
    } catch (error) {
        console.error('Error executing task:', error);
        alert('执行出错: ' + error.message);
    }
}
