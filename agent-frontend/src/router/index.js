import { createRouter, createWebHistory } from 'vue-router'

import Home from '../pages/Home.vue'
import LoveMaster from '../pages/LoveMaster.vue'
import SuperAgent from '../pages/SuperAgent.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/love', name: 'love', component: LoveMaster },
    { path: '/agent', name: 'agent', component: SuperAgent }
  ]
})

export default router


