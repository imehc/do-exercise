import { createContext, useContext } from 'react'

export interface KeepAliveTabContextType {
  refreshTab: (path?: string) => void
  closeTab: (path?: string) => void
  closeOtherTab: (path?: string) => void
}

const defaultValue = {
  refreshTab: () => {},
  closeTab: () => {},
  closeOtherTab: () => {}
}

export const KeepAliveTabContext = createContext<KeepAliveTabContextType>(defaultValue)

export const useTabsHandler = () => {
  const context = useContext(KeepAliveTabContext)
  return context
}
