import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'
import { client } from '@/main'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    },
    {
      path: '/states',
      name: 'stateList',
      component: () => import('../views/StateListView.vue')
    },
    {
      path: '/sign-in',
      name: 'sign-in',
      component: () => import('../views/SignInView.vue')
    }
  ]
})

const getUser = async () => {
  const user = await client.user.getProfile({})
  return user.data
}
router.beforeEach(async (to) => {
  const publicPages = ['/sign-in']
  const auth = useAuthStore()

  // is public
  if (publicPages.includes(to.path)) {
    return
  }

  if (auth.user === null) {
    try {
      const user = await getUser()
      if (!user) throw new Error('session expired')
      auth.setUser(user)
    } catch (err) {
      return 'sign-in'
    }
  }
})

export default router
