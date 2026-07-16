import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, type RenderResult } from 'vitest-browser-react'
import { type Locator, userEvent } from 'vitest/browser'
import { showSubmittedData } from '@/lib/show-submitted-data'
import { OtpForm } from './otp-form'

const navigate = vi.fn()

vi.mock('@tanstack/react-router', async (orig) => {
  const actual = await orig<typeof import('@tanstack/react-router')>()
  return { ...actual, useNavigate: () => navigate }
})

vi.mock('@/lib/show-submitted-data', () => ({ showSubmittedData: vi.fn() }))

describe('一次性密码表单', () => {
  let screen: RenderResult
  let otpInput: Locator
  let verifyButton: Locator

  beforeEach(async () => {
    vi.clearAllMocks()

    screen = await render(<OtpForm />)
    otpInput = screen.getByLabelText(/^一次性密码$/i)
    verifyButton = screen.getByRole('button', { name: /^验证$/i })
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('输入满 6 位数字前禁用验证按钮', async () => {
    await expect.element(verifyButton).toBeDisabled()

    await userEvent.fill(otpInput, '12345')
    await expect.element(verifyButton).toBeDisabled()

    await userEvent.fill(otpInput, '123456')
    await expect.element(verifyButton).toBeEnabled()
  })

  it('提交一次性密码并在超时后跳转', async () => {
    vi.useFakeTimers()

    await userEvent.fill(otpInput, '123456')
    await userEvent.click(verifyButton)

    expect(showSubmittedData).toHaveBeenCalledOnce()
    expect(showSubmittedData).toHaveBeenCalledWith({ otp: '123456' })

    await vi.advanceTimersByTimeAsync(1000)
    expect(navigate).toHaveBeenCalledWith({ to: '/' })
  })
})
