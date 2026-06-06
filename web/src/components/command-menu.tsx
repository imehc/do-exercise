import React from 'react'
import { useNavigate } from '@tanstack/react-router'
import { IconArrowRightDashed } from '@tabler/icons-react'
import { t } from '@lingui/core/macro'
import { Trans, useLingui } from '@lingui/react/macro'
import { themeList } from '~/commons/theme'
import { useSearch } from '~/provider/search'
import { useTheme } from '~/provider/theme'
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '~/components/ui/command'
import { NavGroup } from './layout/types'
import { ScrollArea } from './ui/scroll-area'

interface Props {
  navGroups: NavGroup[]
}
export function CommandMenu({ navGroups }: Props) {
  const navigate = useNavigate()
  const { setTheme } = useTheme()
  const { t: translate } = useLingui()
  const { open, setOpen } = useSearch()

  const runCommand = React.useCallback(
    (command: () => unknown) => {
      setOpen(false)
      command()
    },
    [setOpen]
  )

  return (
    <CommandDialog modal open={open} onOpenChange={setOpen}>
      <CommandInput placeholder={t`键入命令或搜索...`} />
      <CommandList>
        <ScrollArea type='hover' className='h-72 pr-1'>
          <CommandEmpty>
            <Trans>没有内容.</Trans>
          </CommandEmpty>
          {navGroups.map((group) => (
            <CommandGroup key={group.title} heading={group.title}>
              {group.items.map((navItem, i) => {
                if (navItem.url)
                  return (
                    <CommandItem
                      key={`${navItem.url}-${i}`}
                      value={navItem.title}
                      onSelect={() => {
                        runCommand(() => navigate({ to: navItem.url }))
                      }}
                    >
                      <div className='mr-2 flex h-4 w-4 items-center justify-center'>
                        <IconArrowRightDashed className='text-muted-foreground/80 size-2' />
                      </div>
                      {navItem.title}
                    </CommandItem>
                  )

                return navItem.items?.map((subItem, i) => (
                  <CommandItem
                    key={`${subItem.url}-${i}`}
                    value={subItem.title}
                    onSelect={() => {
                      runCommand(() => navigate({ to: subItem.url }))
                    }}
                  >
                    <div className='mr-2 flex h-4 w-4 items-center justify-center'>
                      <IconArrowRightDashed className='text-muted-foreground/80 size-2' />
                    </div>
                    {subItem.title}
                  </CommandItem>
                ))
              })}
            </CommandGroup>
          ))}
          <CommandSeparator />
          <CommandGroup heading={t`主题`}>
            {themeList.map((item) => {
              const label = translate(item.label)

              return (
                <CommandItem
                  key={item.value}
                  value={label}
                  onSelect={() => runCommand(() => setTheme(item.value))}
                >
                  {item.icon} <span>{label}</span>
                </CommandItem>
              )
            })}
          </CommandGroup>
        </ScrollArea>
      </CommandList>
    </CommandDialog>
  )
}
