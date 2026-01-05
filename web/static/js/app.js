// 标签页切换
document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const tabName = btn.dataset.tab;

        // 切换按钮状态
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');

        // 切换内容面板
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        document.getElementById(tabName).classList.add('active');

        // 加载对应数据
        if (tabName === 'servers') loadServers();
        if (tabName === 'destinations') loadDestinations();
        if (tabName === 'tasks') loadTasks();
        if (tabName === 'logs') loadLogs();
    });
});

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    loadServers();
});
