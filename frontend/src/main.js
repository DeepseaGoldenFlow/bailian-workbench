import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createVuetify } from 'vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import App from './App.vue'
import router from './router'

const saved = localStorage.getItem('nova-theme') || 'dark'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: saved,
    themes: {
      dark: {
        dark: true,
        colors: {
          background: '#0d1117',
          surface: '#161b22',
          'surface-variant': '#21262d',
          primary: '#58a6ff',
          secondary: '#8b949e',
          error: '#f85149',
          info: '#79c0ff',
          success: '#3fb950',
          warning: '#d29922',
        }
      },
      light: {
        dark: false,
        colors: {
          background: '#ffffff',
          surface: '#f6f8fa',
          'surface-variant': '#eaeef2',
          primary: '#0969da',
          secondary: '#57606a',
          error: '#cf222e',
          info: '#0969da',
          success: '#1a7f37',
          warning: '#9a6700',
        }
      }
    }
  }
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(vuetify)
app.mount('#app')
