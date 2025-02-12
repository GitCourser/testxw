// ==UserScript==
// @name         玄武日志
// @namespace    Courser Scripts
// @version      0.1
// @description  玄武日志
// @author       Courser
// @include      http://*:4165/admin/
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    const targetNode = document.body;
    const config = { childList: true, attributes: true, subtree: true };

    const callback = function(mutationsList, observer) {
        // console.log('Detected DOM changes:', mutationsList);
        for(let mutation of mutationsList) {
            if (mutation.type === 'childList' || mutation.type === 'attributes') {
                // 检查URL是否包含特定的哈希
                if (window.location.hash === '#/logs') {
                    // 查找表格
                    const table = document.querySelector('table');
                    if (table) {
                        console.log('找到表格');
                        // 获取所有名称单元格
                        const nameCells = table.querySelectorAll('td[data-label="名称"]');
                        nameCells.forEach(cell => {
                            // 移除可能已经存在的点击事件
                            cell.removeEventListener('click', handleCellClick);
                            // 创建一个点击事件监听器
                            cell.addEventListener('click', handleCellClick);
                        });
                    }
                }
            }
        }
    };

    const handleCellClick = async function() {
        const name = this.textContent.trim();
        const url = `${window.location.origin}/api/cron/getlog?name=${encodeURIComponent(name)}`;
        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('网络请求失败');
            }
            const data = await response.json();
            displayData(name, data.data);
        } catch (error) {
            alert('获取内容失败：' + error.message);
        }
    };

    let popup = null;

    const displayData = function(title, data) {
        // 创建一个自定义的弹出窗口
        popup = document.createElement('div');
        popup.classList.add('popup');
        popup.innerHTML = `
            <div class="popup-title">
                <span>${title}</span>
            </div>
            <div class="popup-content">
                <pre>${data}</pre>
            </div>
            <div class="popup-close">
                <span>&times;</span>
            </div>
        `;
        document.body.appendChild(popup);

        // 添加关闭按钮的事件
        const closeButton = popup.querySelector('.popup-close');
        closeButton.addEventListener('click', function(event) {
            event.stopPropagation();
            closePopup();
        });

        // 添加点击背景关闭弹出窗口的功能
        document.body.addEventListener('click', function(event) {
            if (popup && event.target !== popup && !popup.contains(event.target)) {
                closePopup();
            }
        });
    };

    const closePopup = function() {
        if (popup && document.body.contains(popup)) {
            document.body.removeChild(popup);
            popup = null;
        }
    };

    const observer = new MutationObserver(callback);
    observer.observe(targetNode, config);

    console.log('Observer started');

    // 向文档的<head>部分添加样式
    const style = document.createElement('style');
    style.innerHTML = `
        .popup {
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: rgba(0, 0, 0, 0.7);
            border: 1px solid #ccc;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            max-width: 80%;
            max-height: 80%;
            display: flex;
            flex-direction: column;
            z-index: 1000;
        }
        .popup-title {
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 5px;
            border-bottom: 1px solid #ccc;
            color: #eee;
            font-size: 1.4em;
        }
        .popup-content {
            flex: 1;
            overflow-y: auto;
            padding: 15px;
        }
        .popup-content pre {
            color: #0d0;
            font-family: '霞鹜文楷 屏幕阅读版', sans-serif;
            white-space: pre-line;
            word-wrap: break-word;
        }
        .popup-close {
            display: flex;
            justify-content: center;
            align-items: center;
            border-top: 1px solid #ccc;
            color: #a00;
            cursor: pointer;
            font-size: 2em;
        }
        .popup-close:hover {
            color: #f00;
        }
        td[data-label="名称"] {
            cursor: pointer;
        }
        /* 通用滚动条样式 */
        ::-webkit-scrollbar {
            width: 10px; /* 滚动条的宽度 */
            height: 10px; /* 滚动条的高度 */
        }
        /* 滚动条轨道样式 */
        ::-webkit-scrollbar-track {
            background-color: rgba(64, 64, 64, 0.3); /* 轨道颜色 */
            border-radius: 4px; /* 轨道圆角 */
        }
        /* 滚动条滑块样式 */
        ::-webkit-scrollbar-thumb {
            background-color: rgba(128, 128, 128, 0.5); /* 滑块颜色 */
            border-radius: 4px; /* 滑块圆角 */
        }
        /* 滚动条滑块在悬停时的样式 */
        ::-webkit-scrollbar-thumb:hover {
            background-color: rgba(160, 160, 160, 0.7); /* 滑块悬停颜色 */
        }
    `;
    document.head.appendChild(style);
})();