import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createVuetify } from 'vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import './styles.css'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import App from './App.vue'
import router from './router'

const saved = localStorage.getItem('nova-theme') || 'light'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: saved,
    themes: {
      dark: {
        dark: true,
        colors: {
          background: '#0b1020', surface: '#12192b', 'surface-variant': '#1b2540', 'on-background': '#e5e7eb', 'on-surface': '#e5e7eb', 'on-surface-variant': '#94a3b8', primary: '#a5b4fc', secondary: '#7dd3fc', error: '#fb7185', info: '#38bdf8', success: '#34d399', warning: '#fbbf24',
        }
      },
      light: {
        dark: false,
        colors: {
          background: '#f5f7fb', surface: '#ffffff', 'surface-variant': '#eef1f8', 'on-background': '#172033', 'on-surface': '#172033', 'on-surface-variant': '#64748b', primary: '#4f46e5', secondary: '#0284c7', error: '#e11d48', info: '#0284c7', success: '#059669', warning: '#d97706',
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
