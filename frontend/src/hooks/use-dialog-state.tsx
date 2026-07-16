import { useState } from 'react'

/**
 * 确认弹窗使用的状态钩子
 * @param initialState 字符串或空值
 * @returns 当前状态值和用于更新它的函数
 * @example const [open, setOpen] = useDialogState<"approve" | "reject">()
 */
export default function useDialogState<T extends string | boolean>(
  initialState: T | null = null
) {
  const [open, _setOpen] = useState<T | null>(initialState)

  const setOpen = (str: T | null) =>
    _setOpen((prev) => (prev === str ? null : str))

  return [open, setOpen] as const
}
