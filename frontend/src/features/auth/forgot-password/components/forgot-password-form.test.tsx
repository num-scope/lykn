import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, type RenderResult } from 'vitest-browser-react'
import { userEvent, type Locator } from 'vitest/browser'
import { ForgotPasswordForm } from './forgot-password-form'

const navigateMock = vi.fn()

vi.mock('@tanstack/react-router', async (orig) => {
  const actual = await orig<typeof import('@tanstack/react-router')>()
  return { ...actual, useNavigate: () => navigateMock }
})

vi.mock('@/lib/utils', async (orig) => ({
  ...(await orig()),
  sleep: vi.fn(() => Promise.resolve()),
}))

describe('忘记密码表单', () => {
  let screen: RenderResult
  let emailInput: Locator
  let continueButton: Locator

  beforeEach(async () => {
    vi.clearAllMocks()

    screen = await render(<ForgotPasswordForm />)
    emailInput = screen.getByRole('textbox', { name: /^邮箱$/i })
    continueButton = screen.getByRole('button', { name: /^继续$/i })
  })

  it('渲染邮箱字段和继续按钮', async () => {
    await expect.element(emailInput).toBeInTheDocument()
    await expect.element(continueButton).toBeInTheDocument()
  })

  it('提交空表单时显示校验消息', async () => {
    await userEvent.click(continueButton)
    await expect.element(screen.getByText(/^请输入邮箱$/i)).toBeInTheDocument()
  })

  it('成功后重置表单并跳转到 /otp', async () => {
    await userEvent.fill(emailInput, 'a@b.com')
    await userEvent.click(continueButton)

    await vi.waitFor(() =>
      expect(navigateMock).toHaveBeenCalledWith({ to: '/otp' })
    )

    // 成功后应重置表单
    await expect.element(emailInput).toHaveValue('')
  })
})
