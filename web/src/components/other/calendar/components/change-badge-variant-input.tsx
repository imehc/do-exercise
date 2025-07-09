import { Trans } from '@lingui/react/macro'
import { IconPoint, IconPalette } from '@tabler/icons-react'
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
                <IconPoint className='h-5 w-5' />
              ) : (
                <IconPalette className='h-5 w-5' />
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
