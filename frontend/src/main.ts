import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createConnectTransport } from '@connectrpc/connect-web'
import { createPromiseClient } from '@connectrpc/connect'
import { StateService } from './gen/api/state/v1/state_connect'
import { UserService } from './gen/api/user/v1/user_connect'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin)

app.mount('#app')

const transport = createConnectTransport({
  baseUrl: 'http://localhost:8000/cgi-bin/tfstate-manager/',
  credentials: 'include'
})
export const client = {
  state: createPromiseClient(StateService, transport),
  user: createPromiseClient(UserService, transport)
}
