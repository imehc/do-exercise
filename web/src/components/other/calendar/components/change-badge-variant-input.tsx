import { Trans } from '@lingui/react/macro'
import { DotIcon, PaletteIcon } from 'lucide-react'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '~/components/ui/tooltip'
import { buttonHover } from '../animations'
import { MotionButton } from '../components/header/calendar-header'
import { useCalendar } from '../contexts/calendar-context'

export function ChangeBadgeVariantInput() {
  const { badgeVariant, setBadgeVariant } = useCalendar()

  return (
    <div className='space-y-1'>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            <MotionButton
              variant='outline'
              size='icon'
              variants={buttonHover}
              whileHover='hover'
              whileTap='tap'
              onClick={() =>
                setBadgeVariant(badgeVariant === 'dot' ? 'colored' : 'dot')
              }
            >
              {badgeVariant === 'dot' ? (
                <DotIcon className='h-5 w-5' />
              ) : (
                <PaletteIcon className='h-5 w-5' />
              )}
            </MotionButton>
          </TooltipTrigger>
          <TooltipContent>
            <Trans>徽章样式</Trans>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  )
}
