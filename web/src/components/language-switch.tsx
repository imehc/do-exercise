import { IconCheck, IconLanguage } from '@tabler/icons-react'
import { i18n } from '@lingui/core'
import { Trans } from '@lingui/react/macro'
import { useAtom } from 'jotai'
import { languageAtom, languages } from '~/atoms'
import { cn } from '~/lib/utils'
import { Button } from '~/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'

export function LanguageSwitch() {
  const [lang, setLanguage] = useAtom(languageAtom)
  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild>
        <Button variant='ghost' size='icon' className='scale-95 rounded-full'>
          <IconLanguage className='size-[1.2rem]' />
          <span className='sr-only'>
            <Trans>切换语言</Trans>
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end'>
        {languages.map((language) => (
          <DropdownMenuItem
            key={language.value}
            onClick={() => {
              setLanguage(language.value)
              i18n.activate(language.value)
            }}
          >
            {language.label}
            <IconCheck
              size={14}
              className={cn('ml-auto', lang !== language.value && 'hidden')}
            />
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
