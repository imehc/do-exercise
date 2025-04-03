import { GlobeIcon } from 'lucide-react'
import i18n from 'i18next'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '~/components/ui/dropdown-menu'
import { useReducer, useEffect } from 'react'

export function LanguageToggle() {
  const [, forceUpdate] = useReducer(x => x + 1, 0)

  useEffect(() => {
    const handleLanguageChange = () => forceUpdate()
    i18n.on('languageChanged', handleLanguageChange)
    return () => i18n.off('languageChanged', handleLanguageChange)
  }, [])

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="icon">
          <GlobeIcon className="h-[1.2rem] w-[1.2rem]" />
          <span className="sr-only">Toggle language</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem
          onClick={() => i18n.changeLanguage('en')}
          className={i18n.language === 'en' ? 'bg-accent' : ''}
        >
          English
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={() => i18n.changeLanguage('zh')}
          className={i18n.language === 'zh' ? 'bg-accent' : ''}
        >
          中文
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
