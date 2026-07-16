import type { SubmitEvent } from 'react'
import { describe, expect, it, vi } from 'vitest'
import { render } from 'vitest-browser-react'
import { userEvent } from 'vitest/browser'
import { ConfirmDialog } from './confirm-dialog'

describe('确认弹窗', () => {
  it('渲染标题、描述和默认按钮', async () => {
    const { getByRole, getByText } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='删除项目'
        desc='此操作无法撤销'
        handleConfirm={vi.fn()}
      />
    )

    await expect
      .element(getByRole('heading', { name: '删除项目' }))
      .toBeInTheDocument()
    await expect.element(getByText('此操作无法撤销')).toBeInTheDocument()
    await expect
      .element(getByRole('button', { name: '取消' }))
      .toBeInTheDocument()
    await expect
      .element(getByRole('button', { name: '继续' }))
      .toBeInTheDocument()
  })

  it('点击确认按钮时调用 确认处理函数', async () => {
    const handleConfirm = vi.fn()
    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='退出登录'
        desc='确定要继续吗？'
        confirmText='退出登录'
        handleConfirm={handleConfirm}
      />
    )

    await userEvent.click(getByRole('button', { name: '退出登录' }))
    expect(handleConfirm).toHaveBeenCalledOnce()
  })

  it('禁用状态为 true 时禁用确认按钮', async () => {
    const handleConfirm = vi.fn()
    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='危险操作'
        desc='...'
        disabled
        handleConfirm={handleConfirm}
      />
    )

    const confirm = getByRole('button', { name: '继续' })
    await expect.element(confirm).toBeDisabled()
    expect(handleConfirm).not.toHaveBeenCalled()
  })

  it('加载状态为 true 时禁用取消和确认按钮', async () => {
    const handleConfirm = vi.fn()
    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='加载中'
        desc='...'
        isLoading
        handleConfirm={handleConfirm}
      />
    )

    await expect.element(getByRole('button', { name: '取消' })).toBeDisabled()
    await expect.element(getByRole('button', { name: '继续' })).toBeDisabled()
  })

  it('支持自定义按钮文案', async () => {
    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='删除'
        desc='...'
        cancelBtnText='否'
        confirmText='是'
        handleConfirm={vi.fn()}
      />
    )

    await expect
      .element(getByRole('button', { name: '否' }))
      .toBeInTheDocument()
    await expect
      .element(getByRole('button', { name: '是' }))
      .toBeInTheDocument()
  })

  it('设置 `form` 时将确认按钮渲染为关联描述表单的提交按钮', async () => {
    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='删除任务'
        form='tasks-multi-delete-form'
        desc={
          <form id='tasks-multi-delete-form' className='space-y-4'>
            <p>输入 DELETE 以确认</p>
          </form>
        }
        confirmText='删除'
        destructive
      />
    )

    const deleteBtn = getByRole('button', { name: '删除' })
    await expect.element(deleteBtn).toHaveAttribute('type', 'submit')
    await expect
      .element(deleteBtn)
      .toHaveAttribute('form', 'tasks-multi-delete-form')
  })

  it('点击确认时提交描述表单（设置 form 且未传 确认处理函数）', async () => {
    const handleFormSubmit = vi.fn((e: SubmitEvent<HTMLFormElement>) => {
      e.preventDefault()
    })

    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='删除'
        form='users-delete-form'
        desc={
          <form
            id='users-delete-form'
            onSubmit={handleFormSubmit}
            className='space-y-4'
          >
            <p>确认删除</p>
          </form>
        }
        confirmText='删除'
        destructive
      />
    )

    await userEvent.click(getByRole('button', { name: '删除' }))

    expect(handleFormSubmit).toHaveBeenCalledOnce()
  })

  it('按下 回车键时提交表单', async () => {
    const handleFormSubmit = vi.fn((e: SubmitEvent<HTMLFormElement>) => {
      e.preventDefault()
    })

    const { getByPlaceholder } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='删除'
        form='users-delete-form'
        desc={
          <form
            id='users-delete-form'
            onSubmit={handleFormSubmit}
            className='space-y-4'
          >
            <input type='text' name='username' placeholder='用户名' />
          </form>
        }
        confirmText='删除'
        destructive
      />
    )

    await userEvent.fill(getByPlaceholder('用户名'), 'test')
    await userEvent.keyboard('{Enter}')
    expect(handleFormSubmit).toHaveBeenCalledOnce()
  })

  it('确认按钮禁用时不提交表单（输入确认内容不匹配）', async () => {
    const handleFormSubmit = vi.fn((e: SubmitEvent<HTMLFormElement>) => {
      e.preventDefault()
    })

    const { getByRole } = await render(
      <ConfirmDialog
        open
        onOpenChange={vi.fn()}
        title='删除'
        form='users-delete-form'
        disabled
        desc={
          <form id='users-delete-form' onSubmit={handleFormSubmit}>
            <p>输入用户名以启用删除</p>
          </form>
        }
        confirmText='删除'
        destructive
      />
    )

    const deleteBtn = getByRole('button', { name: '删除' })
    await expect.element(deleteBtn).toBeDisabled()
    expect(handleFormSubmit).not.toHaveBeenCalled()
  })
})
