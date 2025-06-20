import { createStore } from 'jotai'
import { atomWithStorage } from 'jotai/utils'
import { type Token } from '~/do-exercise-api'
import { languageAtom } from './language'

export const originTokenAtom = atomWithStorage<Token>('tokenAtom', {
  accessToken: '',
  expireTime: 0,
  refreshToken: '',
  refreshExpireTime: 0,
})

const store = createStore()
store.sub(originTokenAtom, () => {
  // console.log('originTokenAtom change')
})
store.sub(languageAtom, () => {
  // console.log('languageAtom change')
})

export { store }
