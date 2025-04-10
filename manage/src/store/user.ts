import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import localforage from 'localforage'
import { type Token } from '~/do-exercise-api'

type UserStore = {
  auth?: Token
  setAuth: (auth: Token) => void
  clearAuth: () => void
}

export const useUserStore = create(
  persist<UserStore>(
    set => ({
      auth: undefined,
      setAuth: (auth: Token) => set(() => ({ auth })),
      clearAuth: () => set(() => ({ auth: undefined }))
    }),
    {
      name: 'user-storage',
      storage: createJSONStorage(() => localforage)
    }
  )
)
