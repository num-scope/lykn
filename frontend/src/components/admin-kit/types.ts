import type { Component, VNodeChild } from 'vue'

export type AdminMenuIcon = Component | (() => VNodeChild)
export type AdminMenuTheme = 'light' | 'dark'
export type AdminSearchFieldType = 'input' | 'select' | 'dateRange' | 'number' | 'custom'

export interface AdminMenuItem {
  key: string
  path?: string
  label: string
  hint?: string
  icon?: AdminMenuIcon
  children?: AdminMenuItem[]
  danger?: boolean
}

export interface AdminPaginationState {
  current: number
  pageSize: number
  total: number
}

export interface AdminSearchFieldOption {
  label: string
  value: string | number | boolean
  disabled?: boolean
}

export interface AdminSearchField {
  key: string
  label: string
  type: AdminSearchFieldType
  placeholder?: string
  options?: AdminSearchFieldOption[]
  clearable?: boolean
  disabled?: boolean
}

export interface AdminLoginPayload {
  username: string
  password: string
}

export interface AdminLoginHighlight {
  icon?: Component
  title: string
  note: string
}

export interface AdminRowAction<T = unknown> {
  key: string
  label: string
  danger?: boolean
  disabled?: boolean | ((record: T) => boolean)
  visible?: boolean | ((record: T) => boolean)
  confirm?: string | ((record: T) => string)
  confirmOkText?: string
  confirmCancelText?: string
  icon?: Component
  onClick: (record: T) => void | Promise<void>
}
