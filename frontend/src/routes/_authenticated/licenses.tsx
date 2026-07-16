import { createFileRoute } from '@tanstack/react-router'
import { LicensesPage } from '@/features/licenses'

export const Route = createFileRoute('/_authenticated/licenses')({
  component: LicensesPage,
})
