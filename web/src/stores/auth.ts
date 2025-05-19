import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import { type Token } from '~/do-exercise-api'
import { addSeconds } from 'date-fns'

type AuthStore = {
  auth?: Token
  setAuth: (auth: Token) => void
  clearAuth: () => void
}

export const useAuthStore = create(
  persist<AuthStore>(
    set => ({
      auth: undefined,
      setAuth: (auth: Token) =>
        set(() => ({
          auth: {
            ...auth,
            expireTime: addSeconds(new Date(), auth.expireTime).getTime(),
            refreshExpireTime: addSeconds(new Date(), auth.refreshExpireTime).getTime()
          }
        })),
      clearAuth: () => set(() => ({ auth: undefined }))
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage)
    }
  )
)
