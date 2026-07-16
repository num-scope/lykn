import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, type RenderResult } from 'vitest-browser-react'
import { type Locator, userEvent } from 'vitest/browser'
import { UserAuthForm } from './user-auth-form'

const FORM_MESSAGES = {
  usernameEmpty: '请输入用户名',
  passwordEmpty: '请输入密码',
} as const

const navigate = vi.fn()
const loginMock = vi.fn()

vi.mock('@/stores/auth-store', () => ({
  useAuthStore: () => ({
    auth: {
      login: loginMock,
    },
  }),
}))

vi.mock('@tanstack/react-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@tanstack/react-router')>()
  return {
    ...actual,
    useNavigate: () => navigate,
  }
})

describe('登录表单', () => {
  describe('未提供 redirectTo 时的渲染', () => {
    let screen: RenderResult
    let usernameInput: Locator
    let passwordInput: Locator
    let signInButton: Locator

    beforeEach(async () => {
      vi.clearAllMocks()
      loginMock.mockResolvedValue(undefined)
      screen = await render(<UserAuthForm />)
      usernameInput = screen.getByRole('textbox', { name: /^用户名$/i })
      passwordInput = screen.getByLabelText(/^密码$/i)
      signInButton = screen.getByRole('button', { name: /^登录$/i })
    })

    it('渲染字段和提交按钮', async () => {
      await expect.element(usernameInput).toBeInTheDocument()
      await expect.element(passwordInput).toBeInTheDocument()
      await expect.element(signInButton).toBeInTheDocument()
    })

    it('提交空表单时显示校验消息', async () => {
      await userEvent.clear(usernameInput)
      await userEvent.click(signInButton)

      await expect
        .element(screen.getByText(FORM_MESSAGES.usernameEmpty))
        .toBeInTheDocument()
      await expect
        .element(screen.getByText(FORM_MESSAGES.passwordEmpty))
        .toBeInTheDocument()
    })

    it('认证成功后跳转到默认路由', async () => {
      await userEvent.clear(usernameInput)
      await userEvent.fill(usernameInput, 'admin')
      await userEvent.fill(passwordInput, 'admin123')

      await userEvent.click(signInButton)

      await vi.waitFor(() => expect(loginMock).toHaveBeenCalledOnce())
      expect(loginMock).toHaveBeenCalledWith('admin', 'admin123')

      await vi.waitFor(() =>
        expect(navigate).toHaveBeenCalledWith({ to: '/', replace: true })
      )
    })
  })

  it('提供 redirectTo 时跳转到指定地址', async () => {
    vi.clearAllMocks()
    loginMock.mockResolvedValue(undefined)

    const { getByRole, getByLabelText } = await render(
      <UserAuthForm redirectTo='/projects' />
    )

    const usernameInput = getByRole('textbox', { name: /用户名/i })
    await userEvent.clear(usernameInput)
    await userEvent.fill(usernameInput, 'admin')
    await userEvent.fill(getByLabelText('密码'), 'admin123')
    await userEvent.click(getByRole('button', { name: /登录/i }))

    await vi.waitFor(() => expect(loginMock).toHaveBeenCalledOnce())
    await vi.waitFor(() =>
      expect(navigate).toHaveBeenCalledWith({
        to: '/projects',
        replace: true,
      })
    )
  })
})
