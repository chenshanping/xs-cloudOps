```bash
// 在 App.vue 或 store 中
function setupTokenRefreshTimer() {
  setInterval(async () => {
    const userStore = useUserStore()
    if (!userStore.token) return
    
    // 解析 Token 过期时间
    try {
      const payload = JSON.parse(atob(userStore.token.split('.')[1]))
      const exp = payload.exp * 1000 // 转毫秒
      const now = Date.now()
      
      // 还剩 5 分钟时刷新
      if (exp - now < 5 * 60 * 1000 && exp > now) {
        const res = await refreshToken()
        userStore.token = res.data.token
        localStorage.setItem('token', res.data.token)
      }
    } catch {}
  }, 60 * 1000) // 每分钟检查一次
}
```

