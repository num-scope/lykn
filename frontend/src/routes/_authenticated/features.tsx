import { createFileRoute } from '@tanstack/react-router'
import { FeaturesPage } from '@/features/features'

export const Route = createFileRoute('/_authenticated/features')({
  component: FeaturesPage,
})
