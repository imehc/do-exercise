import { createStore } from 'jotai'
import { atomWithStorage } from 'jotai/utils'
import { type Token } from '~/do-exercise-api'

export const originTokenAtom = atomWithStorage<Token>('tokenAtom', {
  accessToken: '',
  expireTime: 0,
  refreshToken: '',
  refreshExpireTime: 0,
})

const store = createStore()
store.sub(originTokenAtom, () => {})

export { store }
