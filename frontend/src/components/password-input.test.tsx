import { useForm } from 'react-hook-form'
import { describe, expect, it } from 'vitest'
import { render } from 'vitest-browser-react'
import { userEvent } from 'vitest/browser'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form'
import { PasswordInput } from './password-input'

describe('密码输入框', () => {
  it('正确渲染密码输入框', async () => {
    const { getByPlaceholder, getByRole } = await render(
      <PasswordInput placeholder='密码' />
    )

    const passwordInput = getByPlaceholder('密码')
    const showPasswordButton = getByRole('button', { name: /显示密码/i })

    await expect.element(passwordInput).toBeInTheDocument()
    await expect.element(passwordInput).toHaveAttribute('type', 'password')
    await expect.element(showPasswordButton).toBeVisible()
  })

  it('点击显示密码按钮时切换密码可见性', async () => {
    const { getByPlaceholder, getByRole } = await render(
      <PasswordInput placeholder='密码' />
    )

    const passwordInput = getByPlaceholder('密码')
    const showPasswordButton = getByRole('button', { name: /显示密码/i })

    await expect.element(passwordInput).toHaveAttribute('type', 'password')
    await expect.element(showPasswordButton).toBeInTheDocument()

    await userEvent.click(showPasswordButton)

    await expect.element(passwordInput).toHaveAttribute('type', 'text')
    const hidePasswordButton = getByRole('button', { name: /隐藏密码/i })
    await expect.element(hidePasswordButton).toBeInTheDocument()

    await userEvent.click(hidePasswordButton)

    await expect.element(passwordInput).toHaveAttribute('type', 'password')
    await expect
      .element(getByRole('button', { name: /显示密码/i }))
      .toBeInTheDocument()
  })

  it('密码输入框禁用时同时禁用显示密码按钮', async () => {
    const { getByPlaceholder, getByRole } = await render(
      <PasswordInput placeholder='密码' disabled />
    )

    const passwordInput = getByPlaceholder('密码')
    const showPasswordButton = getByRole('button', { name: /显示密码/i })
    await expect.element(showPasswordButton).toBeDisabled()
    await expect.element(passwordInput).toBeDisabled()
  })

  it('配合 FormLabel 和 react-hook-form 字段展开使用', async () => {
    function PasswordInLabeledForm() {
      const form = useForm<{ password: string }>({
        defaultValues: { password: '' },
      })

      return (
        <Form {...form}>
          <form>
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>密码</FormLabel>
                  <FormControl>
                    <PasswordInput {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
          </form>
        </Form>
      )
    }

    const { getByLabelText } = await render(<PasswordInLabeledForm />)

    const password = getByLabelText(/^密码$/i)
    await expect.element(password).toHaveAttribute('type', 'password')

    await userEvent.type(password, 'secret-value')

    await expect.element(password).toHaveValue('secret-value')
  })
})
