// 保存任务配置
async function saveTask(event) {
    event.preventDefault();

    const id = document.getElementById('taskId').value;

    // 获取选中的备份目标
    const checkboxes = document.querySelectorAll('input[name="destinations"]:checked');
    const destinationIds = Array.from(checkboxes).map(cb => parseInt(cb.value));

    const data = {
        name: document.getElementById('taskName').value,
        source_server_id: parseInt(document.getElementById('sourceServerId').value),
        cron_expression: document.getElementById('cronExpression').value,
        enabled: document.getElementById('enabled').checked,
        destination_ids: destinationIds
    };

    try {
        const url = id ? `/api/tasks/${id}` : '/api/tasks';
        const method = id ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            closeTaskModal();
            loadTasks();
            alert('保存成功');
        } else {
            alert('保存失败');
        }
    } catch (error) {
        console.error('Failed to save task:', error);
        alert('保存失败');
    }
}

// 删除任务
async function deleteTask(id) {
    if (!confirm('确定要删除这个任务吗？')) return;

    try {
        const response = await fetch(`/api/tasks/${id}`, { method: 'DELETE' });
        if (response.ok) {
            loadTasks();
            alert('删除成功');
        } else {
            alert('删除失败');
        }
    } catch (error) {
        console.error('Failed to delete task:', error);
        alert('删除失败');
    }
}
