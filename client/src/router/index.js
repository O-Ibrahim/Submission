import Layout from '@/components/layout/Layout.vue'
import DashboardView from '@/views/DashboardView.vue'
import Logs from '@/views/LogsView.vue'
import NewCommandView from '@/views/NewCommandView.vue'
import { createRouter, createWebHistory } from 'vue-router'
// import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'layout',
      component: Layout,
      children: [
        {
          path: "",
          name: "home",
          component: DashboardView
        },
        {
          path: "/logs/:id",
          name: "logs",
          component: Logs
        },
        {
          path: "/new",
          name: "new",
          component: NewCommandView
        }
      ]
    },
    { path: '/:pathMatch(.*)*', name: "catchall", redirect: '/' }
  ],
})

export default router
