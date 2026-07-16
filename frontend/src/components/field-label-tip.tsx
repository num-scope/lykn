import { CircleHelp } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Label } from '@/components/ui/label'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

type FieldLabelTipProps = React.ComponentProps<typeof Label> & {
  tip: React.ReactNode
}

/** 表单标签 + 问号提示（shadcn Tooltip） */
export function FieldLabelTip({
  tip,
  children,
  className,
  ...props
}: FieldLabelTipProps) {
  return (
    <div className={cn('flex items-center gap-1.5', className)}>
      <Label {...props}>{children}</Label>
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            type='button'
            className='inline-flex size-4 shrink-0 items-center justify-center rounded-full text-muted-foreground transition-colors hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none'
            aria-label='字段说明'
          >
            <CircleHelp className='size-3.5' />
          </button>
        </TooltipTrigger>
        <TooltipContent side='top' className='max-w-xs'>
          {tip}
        </TooltipContent>
      </Tooltip>
    </div>
  )
}
