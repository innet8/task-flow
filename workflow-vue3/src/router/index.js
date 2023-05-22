import { createRouter, createWebHistory } from 'vue-router'

const  routes = [
  {
    name:'manage',
    path: '/',
    component:()=>import('@/views/manage.vue'),
    children: [
      {
        name:'manage-setting',
        path: 'setting',
        component:()=>import('@/views/manage/setting.vue'),
      },
      {
        name:'manage-review',
        path: 'review',
        component:()=>import('@/views/manage/review/index.vue'),
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
