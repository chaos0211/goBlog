function updateTime() {
    const timeElement = document.getElementById("time");
    const now = new Date();
    timeElement.textContent = now.toLocaleString("zh-CN");
}
setInterval(updateTime, 1000);

