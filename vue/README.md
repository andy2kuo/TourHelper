# TourHelper Frontend

TourHelper 前端專案，使用 Vue.js 3 和 Vite 建構，提供旅遊推薦的網頁介面。

## 技術架構

- **框架**：Vue.js 3.5+
- **建置工具**：Vite 7.x
- **UI 框架**：Element Plus 2.x
- **狀態管理**：Pinia
- **圖標庫**：Element Plus Icons
- **開發語言**：JavaScript（可擴展為 TypeScript）

## 專案結構

```text
vue/
├── src/                     # 原始碼
│   ├── components/         # Vue 元件
│   ├── stores/             # Pinia 狀態管理
│   │   └── user.js        # 使用者狀態管理
│   ├── views/              # 頁面元件（待建立）
│   ├── router/             # Vue Router（待建立）
│   ├── api/                # API 請求封裝（待建立）
│   ├── utils/              # 工具函式（待建立）
│   ├── assets/             # 靜態資源
│   ├── App.vue             # 根元件
│   ├── main.js             # 應用程式入口
│   └── style.css           # 全域樣式
├── public/                  # 公開靜態資源
│   └── favicon.ico         # 網站圖標
├── index.html               # HTML 入口檔案
├── vite.config.js           # Vite 配置
├── package.json             # 專案依賴
├── package-lock.json        # 依賴鎖定檔案
└── README.md                # 本檔案
```

## 已安裝套件

### 核心依賴

- **vue** (^3.5.24)：Vue.js 框架
- **element-plus** (^2.11.7)：UI 元件庫
- **@element-plus/icons-vue** (^2.3.2)：Element Plus 圖標
- **pinia** (^3.0.4)：Vue 3 官方狀態管理

### 開發依賴

- **vite** (^7.2.2)：前端建置工具
- **@vitejs/plugin-vue** (^6.0.1)：Vite 的 Vue 插件

## 開發指令

### 安裝依賴

```bash
npm install
```

### 開發模式

啟動開發伺服器（支援熱重載）：

```bash
npm run dev
```

預設會在 `http://localhost:5173` 啟動。

### 建置生產版本

```bash
npm run build
```

建置產物會輸出到 `dist/` 目錄。

### 預覽生產版本

建置完成後可以本地預覽：

```bash
npm run preview
```

## 核心模組說明

### main.js

應用程式入口點，負責：

- 建立 Vue 應用程式實例
- 註冊 Element Plus UI 框架
- 註冊 Pinia 狀態管理
- 註冊所有 Element Plus 圖標為全域元件
- 掛載應用程式到 DOM

```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import './style.css'
import App from './App.vue'

const app = createApp(App)
const pinia = createPinia()

app.use(ElementPlus)
app.use(pinia)

// 註冊所有 Element Plus 圖標
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')
```

### stores/user.js

使用者狀態管理 Store，管理：

- 使用者 ID
- 使用者偏好設定（最大距離、偏好標籤、預算）
- 當前位置資訊

```javascript
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 設定位置
userStore.setLocation(25.0330, 121.5654)

// 更新偏好
userStore.updatePreferences({ maxDistance: 100 })

// 重置偏好
userStore.resetPreferences()
```

## Element Plus 使用

### 基本用法

Element Plus 已經全域註冊，可以直接在任何元件中使用：

```vue
<template>
  <div>
    <el-button type="primary">主要按鈕</el-button>
    <el-input v-model="input" placeholder="請輸入內容" />
  </div>
</template>

<script setup>
import { ref } from 'vue'

const input = ref('')
</script>
```

### 使用圖標

所有 Element Plus 圖標已全域註冊：

```vue
<template>
  <el-button type="primary" :icon="Search">搜尋</el-button>
  <el-icon><Location /></el-icon>
</template>

<script setup>
import { Search, Location } from '@element-plus/icons-vue'
</script>
```

### 常用元件

- **Layout**：`el-container`, `el-header`, `el-main`, `el-footer`
- **Button**：`el-button`, `el-button-group`
- **Form**：`el-form`, `el-form-item`, `el-input`, `el-select`
- **Data**：`el-table`, `el-pagination`, `el-tag`
- **Feedback**：`el-dialog`, `el-message`, `el-notification`
- **Navigation**：`el-menu`, `el-tabs`, `el-breadcrumb`

詳細文件請參考 [Element Plus 官方文件](https://element-plus.org/)

## Pinia 狀態管理

### Store 結構

使用 Composition API 風格定義 Store：

```javascript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useMyStore = defineStore('myStore', () => {
  // State
  const count = ref(0)

  // Getters
  const doubleCount = computed(() => count.value * 2)

  // Actions
  function increment() {
    count.value++
  }

  return { count, doubleCount, increment }
})
```

### 在元件中使用

```vue
<script setup>
import { useMyStore } from '@/stores/myStore'

const myStore = useMyStore()

// 直接訪問 state
console.log(myStore.count)

// 調用 actions
myStore.increment()
</script>
```

## 開發規範

### 元件命名

- **檔案名稱**：使用 PascalCase，例如 `MyComponent.vue`
- **元件名稱**：使用 PascalCase，例如 `<MyComponent />`

### 目錄結構建議

```text
src/
├── components/          # 可重用元件
│   ├── common/         # 通用元件（按鈕、輸入框等）
│   ├── layout/         # 佈局元件（Header、Footer、Sidebar）
│   └── business/       # 業務元件（RecommendationCard、MapView）
├── views/              # 頁面元件
│   ├── Home.vue
│   ├── Recommendations.vue
│   └── Settings.vue
├── stores/             # Pinia Stores
│   ├── user.js
│   ├── recommendation.js
│   └── location.js
├── api/                # API 請求
│   ├── request.js      # axios 封裝
│   ├── user.js
│   └── recommendation.js
├── utils/              # 工具函式
│   ├── format.js       # 格式化工具
│   └── validation.js   # 驗證工具
└── assets/             # 靜態資源
    ├── images/
    └── styles/
```

### 樣式規範

建議使用 scoped 樣式避免污染全域：

```vue
<style scoped>
.my-component {
  /* 元件樣式 */
}
</style>
```

全域樣式放在 `src/style.css` 或 `src/assets/styles/`。

## API 整合

### 安裝 Axios

```bash
npm install axios
```

### API 封裝範例

建立 `src/api/request.js`：

```javascript
import axios from 'axios'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000
})

// 請求攔截器
request.interceptors.request.use(
  config => {
    // 可在此添加 token
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 回應攔截器
request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    // 統一錯誤處理
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export default request
```

建立 `src/api/recommendation.js`：

```javascript
import request from './request'

export const getRecommendations = (params) => {
  return request({
    url: '/api/v1/recommendations',
    method: 'get',
    params
  })
}

export const updatePreferences = (data) => {
  return request({
    url: '/api/v1/user/preferences',
    method: 'post',
    data
  })
}
```

### 在元件中使用

```vue
<script setup>
import { ref } from 'vue'
import { getRecommendations } from '@/api/recommendation'

const recommendations = ref([])

const fetchRecommendations = async () => {
  try {
    const data = await getRecommendations({
      lat: 25.0330,
      lon: 121.5654
    })
    recommendations.value = data.recommendations
  } catch (error) {
    console.error('獲取推薦失敗:', error)
  }
}
</script>
```

## 環境變數

在專案根目錄建立環境變數檔案：

### .env.development

```bash
VITE_API_BASE_URL=http://localhost:8080
```

### .env.production

```bash
VITE_API_BASE_URL=https://api.yourserver.com
```

### 使用環境變數

```javascript
const apiUrl = import.meta.env.VITE_API_BASE_URL
```

## 路由設定（待實作）

### 安裝 Vue Router

```bash
npm install vue-router@4
```

### 基本設定

建立 `src/router/index.js`：

```javascript
import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/recommendations',
    name: 'Recommendations',
    component: () => import('@/views/Recommendations.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
```

在 `main.js` 中註冊：

```javascript
import router from './router'

app.use(router)
```

## 建置優化

### Vite 配置

編輯 `vite.config.js`：

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

## 測試（待實作）

### 安裝測試工具

```bash
npm install -D vitest @vue/test-utils
```

### 測試配置

在 `package.json` 中添加：

```json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui"
  }
}
```

## 部署

### 建置

```bash
npm run build
```

### 部署到靜態託管

建置後的 `dist/` 目錄可以部署到：

- Vercel
- Netlify
- GitHub Pages
- Cloudflare Pages

## 疑難排解

### 開發伺服器無法啟動

1. 確認 Node.js 版本是否符合要求（18+）
2. 刪除 `node_modules` 和 `package-lock.json`，重新執行 `npm install`
3. 檢查端口 5173 是否被佔用

### Element Plus 樣式未載入

確認 `main.js` 中有引入樣式：

```javascript
import 'element-plus/dist/index.css'
```

### API 請求 CORS 錯誤

1. 確認後端伺服器已設定 CORS
2. 或在 `vite.config.js` 中配置 proxy

## 相關連結

- [Vue.js 官方文件](https://vuejs.org/)
- [Vite 官方文件](https://vitejs.dev/)
- [Element Plus 官方文件](https://element-plus.org/)
- [Pinia 官方文件](https://pinia.vuejs.org/)
- [Vue Router 官方文件](https://router.vuejs.org/)
