// static/js/detail.js

// 展开回复
function expandReplies(commentID, totalReplies) {
    var container = document.getElementById('replies-' + commentID);
    var button = container.querySelector('.expand-replies');
    var replies = allReplies[commentID] || [];
    var displayedCount = container.querySelectorAll('.reply').length;

    // 每次展开10条
    var nextCount = Math.min(displayedCount + 10, totalReplies);
    for (var i = displayedCount; i < nextCount; i++) {
        var reply = replies[i];
        var replyDiv = document.createElement('div');
        replyDiv.className = 'reply';
        replyDiv.style.marginLeft = (40 * (Math.floor(i / replies.length) + 1)) + 'px'; // 动态缩进
        replyDiv.id = 'reply-' + reply.id;
        replyDiv.innerHTML = `
            <p>${reply.content}</p>
            <p>用户: ${reply.username}</p>
            <p>回复时间: ${reply.createdAt}</p>
            <div class="comment-actions">
                <button class="like-comment" data-id="${reply.id}">点赞 (<span class="like-count">${reply.likes}</span>)</button>
                <button class="dislike-comment" data-id="${reply.id}">踩 (<span class="dislike-count">${reply.dislikes}</span>)</button>
                <a href="#" class="reply-link" onclick="document.getElementById('reply-form-${reply.id}').style.display='block'; return false;">回复</a>
            </div>
            <div id="reply-form-${reply.id}" style="display: none; margin-left: 20px;">
                <form method="POST" action="/blog/comment">
                    <input type="hidden" name="article_id" value="${articleId}">
                    <input type="hidden" name="parent_id" value="${reply.id}">
                    <div>
                        <label for="username-${reply.id}">用户名:</label>
                        <input type="text" id="username-${reply.id}" name="username" placeholder="匿名用户">
                    </div>
                    <div>
                        <label for="content-${reply.id}">回复内容:</label>
                        <textarea id="content-${reply.id}" name="content" required></textarea>
                    </div>
                    <button type="submit">提交回复</button>
                </form>
            </div>
        `;
        container.insertBefore(replyDiv, button);

        // 渲染嵌套回复（默认显示两条）
        if (reply.replies && reply.replies.length > 0) {
            var nestedRepliesDiv = document.createElement('div');
            nestedRepliesDiv.className = 'replies';
            nestedRepliesDiv.id = 'replies-' + reply.id;
            var nestedDisplayCount = Math.min(2, reply.replies.length); // 默认显示两条
            for (var j = 0; j < nestedDisplayCount; j++) {
                var nestedReply = reply.replies[j];
                var nestedReplyDiv = document.createElement('div');
                nestedReplyDiv.className = 'reply';
                nestedReplyDiv.style.marginLeft = (40 * (Math.floor(i / replies.length) + 2)) + 'px';
                nestedReplyDiv.id = 'reply-' + nestedReply.id;
                nestedReplyDiv.innerHTML = `
                    <p>${nestedReply.content}</p>
                    <p>用户: ${nestedReply.username}</p>
                    <p>回复时间: ${nestedReply.createdAt}</p>
                    <div class="comment-actions">
                        <button class="like-comment" data-id="${nestedReply.id}">点赞 (<span class="like-count">${nestedReply.likes}</span>)</button>
                        <button class="dislike-comment" data-id="${nestedReply.id}">踩 (<span class="dislike-count">${nestedReply.dislikes}</span>)</button>
                        <a href="#" class="reply-link" onclick="document.getElementById('reply-form-${nestedReply.id}').style.display='block'; return false;">回复</a>
                    </div>
                    <div id="reply-form-${nestedReply.id}" style="display: none; margin-left: 20px;">
                        <form method="POST" action="/blog/comment">
                            <input type="hidden" name="article_id" value="${articleId}">
                            <input type="hidden" name="parent_id" value="${nestedReply.id}">
                            <div>
                                <label for="username-${nestedReply.id}">用户名:</label>
                                <input type="text" id="username-${nestedReply.id}" name="username" placeholder="匿名用户">
                            </div>
                            <div>
                                <label for="content-${nestedReply.id}">回复内容:</label>
                                <textarea id="content-${nestedReply.id}" name="content" required></textarea>
                            </div>
                            <button type="submit">提交回复</button>
                        </form>
                    </div>
                `;
                nestedRepliesDiv.appendChild(nestedReplyDiv);
            }
            // 如果嵌套回复超过2条，添加展开按钮
            if (reply.replies.length > 2) {
                var nestedExpandButton = document.createElement('button');
                nestedExpandButton.className = 'expand-replies';
                nestedExpandButton.innerHTML = `展开 ${reply.replies.length} 条回复`;
                nestedExpandButton.onclick = function() {
                    expandReplies(reply.id, reply.replies.length);
                };
                nestedRepliesDiv.appendChild(nestedExpandButton);
            }
            replyDiv.appendChild(nestedRepliesDiv);
        }

        // 重新绑定事件
        bindLikeDislikeEvents(replyDiv);
    }

    // 更新按钮状态
    if (nextCount >= totalReplies) {
        button.innerHTML = 'END';
        button.disabled = true;
    } else {
        button.innerHTML = `展开更多内容`;
    }
}

// 绑定点赞和踩的事件
function bindLikeDislikeEvents(container) {
    // 文章点赞
    (container || document).querySelectorAll('.like-article').forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const articleId = this.getAttribute('data-id');
            fetch(`/blog/like?id=${articleId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        const likeCountSpan = this.querySelector('.like-count');
                        likeCountSpan.textContent = data.likes;
                    } else {
                        alert('点赞失败');
                    }
                })
                .catch(error => {
                    console.error('点赞请求失败:', error);
                    alert('点赞失败');
                });
        });
    });

    // 评论和回复点赞
    (container || document).querySelectorAll('.like-comment').forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const commentId = this.getAttribute('data-id');
            fetch(`/blog/comment/like?id=${commentId}&article_id=${articleId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        const likeCountSpan = this.querySelector('.like-count');
                        likeCountSpan.textContent = data.likes;
                    } else {
                        alert('点赞失败');
                    }
                })
                .catch(error => {
                    console.error('点赞请求失败:', error);
                    alert('点赞失败');
                });
        });
    });

    // 评论和回复踩
    (container || document).querySelectorAll('.dislike-comment').forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const commentId = this.getAttribute('data-id');
            fetch(`/blog/comment/dislike?id=${commentId}&article_id=${articleId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        const dislikeCountSpan = this.querySelector('.dislike-count');
                        dislikeCountSpan.textContent = data.dislikes;
                    } else {
                        alert('点踩失败');
                    }
                })
                .catch(error => {
                    console.error('点踩请求失败:', error);
                    alert('点踩失败');
                });
        });
    });
}

// 页面加载完成后绑定事件
document.addEventListener('DOMContentLoaded', function() {
    // 调试：打印 allReplies
    console.log("allReplies:", allReplies);
    bindLikeDislikeEvents();
});