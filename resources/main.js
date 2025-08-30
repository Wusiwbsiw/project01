// 当整个 HTML 文档加载完成后，再执行我们的代码
document.addEventListener('DOMContentLoaded', () => {

    // --- 1. 处理注册表单 ---
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (event) => {
            // 阻止表单的默认提交行为（即刷新页面）
            event.preventDefault();

            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const resultDiv = document.getElementById('result');

            // 使用 fetch API 向后端发送 POST 请求
            try {
                const response = await fetch('/api/user/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    // 将 JavaScript 对象转换为 JSON 字符串
                    body: JSON.stringify({ username: username, password: password }),
                });

                // 解析后端返回的 JSON 响应
                const data = await response.json();

                if (response.ok) { // HTTP 状态码在 200-299 之间
                    resultDiv.textContent = `注册成功！您的用户 ID 是: ${data.userId}`;
                    resultDiv.className = 'result success';
                } else {
                    resultDiv.textContent = `注册失败: ${data.message}`;
                    resultDiv.className = 'result error';
                }
            } catch (error) {
                resultDiv.textContent = `请求失败: ${error}`;
                resultDiv.className = 'result error';
            }
        });
    }

    // --- 2. 处理登录表单 ---
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (event) => {
            event.preventDefault();

            const userId = parseInt(document.getElementById('userId').value, 10);
            const password = document.getElementById('password').value;
            const resultDiv = document.getElementById('result');
            // 获取提交按钮，以便在请求期间禁用它，防止重复提交
            const submitButton = loginForm.querySelector('button');

            // 禁用按钮，并显示“登录中...”
            submitButton.disabled = true;
            submitButton.textContent = '登录中...';
            resultDiv.textContent = ''; // 清空之前的结果

            try {
                const response = await fetch('/api/user/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ userId: userId, password: password }),
                });

                const data = await response.json();

                if (response.ok && data.success) {
                    // --- 登录成功的核心逻辑 ---

                    // 1. 在 localStorage 中存储 Token
                    // localStorage 是浏览器提供的一个持久化存储区域
                    localStorage.setItem('authToken', data.token);
                    localStorage.setItem('currentUserid', userId);
                    localStorage.setItem('currentUsername', data.username);

                    // 2. 显示成功信息，并准备跳转
                    resultDiv.textContent = `登录成功！欢迎您, ${data.username}。正在跳转到个人主页...`;
                    resultDiv.className = 'result success';

                    // 3. 延迟 2 秒后，跳转到个人主页
                    setTimeout(() => {
                        window.location.href = 'profile.html'; // <-- 这就是跳转命令！
                    }, 500); // 2000 毫秒 = 2 秒

                } else {
                    // 登录失败
                    resultDiv.textContent = `登录失败: ${data.message || '未知错误'}`;
                    resultDiv.className = 'result error';
                    // 重新启用按钮
                    submitButton.disabled = false;
                    submitButton.textContent = '登录';
                }
            } catch (error) {
                // 网络或请求错误
                resultDiv.textContent = `请求失败: ${error}`;
                resultDiv.className = 'result error';
                // 重新启用按钮
                submitButton.disabled = false;
                submitButton.textContent = '登录';
            }
        });
    }

    // --- 3. 处理修改密码表单 ---
    const resetPasswordForm = document.getElementById('resetPasswordForm');
    if (resetPasswordForm) {
        resetPasswordForm.addEventListener('submit', async (event) => {
            event.preventDefault();

            const userId = parseInt(document.getElementById('userId').value, 10);
            const newPassword = document.getElementById('newPassword').value;
            const resultDiv = document.getElementById('result');

            try {
                const response = await fetch('/api/user/reset-password', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ userId: userId, newPassword: newPassword }),
                });

                const data = await response.json();

                if (response.ok) {
                    resultDiv.textContent = `密码修改成功！`;
                    resultDiv.className = 'result success';
                } else {
                    resultDiv.textContent = `密码修改失败: ${data.message}`;
                    resultDiv.className = 'result error';
                }
            } catch (error) {
                resultDiv.textContent = `请求失败: ${error}`;
                resultDiv.className = 'result error';
            }
        });
    }


    if (window.location.pathname.endsWith('/profile.html')) {
        // 1. 从 localStorage 中读取用户信息
        const userId = localStorage.getItem('currentUserid');
        const username = localStorage.getItem('currentUsername');
        const authToken = localStorage.getItem('authToken');

        if (!userId || !username || !authToken) {
            // 如果缺少任何信息，都视为未登录
            alert("请先登录！");
            window.location.href = 'login.html';
            return; // 结束后续代码的执行
        }

        // 2. 立即用本地数据填充页面，实现秒开效果
        document.getElementById('userId').textContent = userId;
        document.getElementById('username').textContent = username;

        // 3. (可选，但推荐) 在后台验证 Token 有效性
        // fetch('/api/profile', {
        //     method: 'GET',
        //     headers: { 'Authorization': `Bearer ${authToken}` },
        // })
        //     .then(response => {
        //         if (!response.ok) {
        //             // 如果 Token 验证失败（比如已过期）
        //             localStorage.clear(); // 清空所有本地存储
        //             alert("您的登录会话已失效，请重新登录。");
        //             window.location.href = 'login.html';
        //         }
        //         // 如果验证成功，什么都不用做，因为页面已经显示了
        //         console.log("Token a validation successful in the background.");
        //     })
        //     .catch(error => {
        //         console.error("Background token validation failed:", error);
        //     });
    }
});