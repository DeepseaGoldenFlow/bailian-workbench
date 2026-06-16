import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createVuetify } from 'vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import App from './App.vue'
import router from './router'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'dark',
    themes: {
      dark: {
        dark: true,
        colors: {
          background: '#121212',
          surface: '#1e1e1e',
          'surface-variant': '#2a2a2a',
          primary: '#8ab4f8',
          secondary: '#a8c7fa',
          error: '#f28b82',
          info: '#aecbfa',
          success: '#81c995',
          warning: '#fdd663',
        }
      },
      light: {
        dark: false,
        colors: {
          background: '#fafafa',
          surface: '#ffffff',
          'surface-variant': '#f1f3f4',
          primary: '#1a73e8',
          secondary: '#5f6368',
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
