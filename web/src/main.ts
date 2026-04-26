import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import App from './App.vue'
import router from './router'
import permission from './directives/permission'
import { useConfigStore } from './store/config'
import 'ant-design-vue/dist/reset.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Antd)

// 注册权限指令
app.directive('permission', permission)

// 先加载配置，再挂载应用
const configStore = useConfigStore()
configStore.loadConfigs().finally(() => {
  app.mount('#app')
})
