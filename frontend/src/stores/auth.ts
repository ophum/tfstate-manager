import type { User } from '@/gen/api/user/v1/user_pb'
import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'

export interface AuthState {
  user: User | null
}

export const useAuthStore = defineStore('auth', {
  state: () =>
    ({
      user: null
    }) as AuthState,
  getters: {
    isSignedIn: (state) => state.user !== null
  },
  actions: {
    setUser(user: User | null) {
      this.user = user
    },
    sessionExpire() {
      const router = useRouter()

      this.setUser(null)
      router.replace('/sign-in')
    }
  }
})
