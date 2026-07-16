import { create } from 'zustand'
import { getCookie, removeCookie, setCookie } from '@/lib/cookies'
import { getCurrentUser, login as loginApi } from '@/api/auth'
import { fetchProjects } from '@/api/lykn'
import type { ProjectRecord, UserAccount } from '@/types/api'

const ACCESS_TOKEN = 'lykn_access_token'
const SELECTED_PROJECT = 'lykn_selected_project'

function readToken() {
  const cookieState = getCookie(ACCESS_TOKEN)
  if (!cookieState) return ''
  try {
    return JSON.parse(cookieState) as string
  } catch {
    return cookieState
  }
}

function readSelectedProjectId() {
  const raw = getCookie(SELECTED_PROJECT)
  if (!raw) return undefined
  try {
    const value = JSON.parse(raw) as number
    return Number.isFinite(value) ? value : undefined
  } catch {
    const value = Number(raw)
    return Number.isFinite(value) ? value : undefined
  }
}

export type AuthUser = {
  id: number
  username: string
}

interface AuthState {
  auth: {
    user: AuthUser | null
    accessToken: string
    projects: ProjectRecord[]
    selectedProjectId?: number
    loadingProjects: boolean
    setUser: (user: AuthUser | null) => void
    setAccessToken: (accessToken: string) => void
    resetAccessToken: () => void
    reset: () => void
    login: (username: string, password: string) => Promise<void>
    hydrateUser: () => Promise<boolean>
    loadProjects: () => Promise<void>
    selectProject: (projectId?: number) => void
  }
}

function toAuthUser(account: UserAccount): AuthUser {
  return {
    id: account.id,
    username: account.username,
  }
}

export const useAuthStore = create<AuthState>()((set, get) => {
  const initToken = readToken()
  const initSelectedProjectId = readSelectedProjectId()

  return {
    auth: {
      user: null,
      accessToken: initToken,
      projects: [],
      selectedProjectId: initSelectedProjectId,
      loadingProjects: false,

      setUser: (user) =>
        set((state) => ({ ...state, auth: { ...state.auth, user } })),

      setAccessToken: (accessToken) =>
        set((state) => {
          setCookie(ACCESS_TOKEN, JSON.stringify(accessToken))
          return { ...state, auth: { ...state.auth, accessToken } }
        }),

      resetAccessToken: () =>
        set((state) => {
          removeCookie(ACCESS_TOKEN)
          return { ...state, auth: { ...state.auth, accessToken: '' } }
        }),

      reset: () =>
        set((state) => {
          removeCookie(ACCESS_TOKEN)
          removeCookie(SELECTED_PROJECT)
          return {
            ...state,
            auth: {
              ...state.auth,
              user: null,
              accessToken: '',
              projects: [],
              selectedProjectId: undefined,
              loadingProjects: false,
            },
          }
        }),

      login: async (username, password) => {
        const data = await loginApi(username, password)
        setCookie(ACCESS_TOKEN, JSON.stringify(data.access_token))
        set((state) => ({
          ...state,
          auth: {
            ...state.auth,
            accessToken: data.access_token,
            user: toAuthUser(data.user),
          },
        }))
        await get().auth.loadProjects()
      },

      hydrateUser: async () => {
        const token = get().auth.accessToken
        if (!token) return false

        try {
          const user = await getCurrentUser()
          set((state) => ({
            ...state,
            auth: { ...state.auth, user: toAuthUser(user) },
          }))
          await get().auth.loadProjects()
          return true
        } catch {
          get().auth.reset()
          return false
        }
      },

      loadProjects: async () => {
        if (!get().auth.accessToken) return

        set((state) => ({
          ...state,
          auth: { ...state.auth, loadingProjects: true },
        }))

        try {
          const projects = await fetchProjects()
          const currentId = get().auth.selectedProjectId
          const nextId =
            currentId && projects.some((item) => item.id === currentId)
              ? currentId
              : projects[0]?.id

          if (nextId) {
            setCookie(SELECTED_PROJECT, JSON.stringify(nextId))
          } else {
            removeCookie(SELECTED_PROJECT)
          }

          set((state) => ({
            ...state,
            auth: {
              ...state.auth,
              projects,
              selectedProjectId: nextId,
              loadingProjects: false,
            },
          }))
        } catch {
          set((state) => ({
            ...state,
            auth: { ...state.auth, loadingProjects: false },
          }))
        }
      },

      selectProject: (projectId) => {
        if (projectId) {
          setCookie(SELECTED_PROJECT, JSON.stringify(projectId))
        } else {
          removeCookie(SELECTED_PROJECT)
        }
        set((state) => ({
          ...state,
          auth: { ...state.auth, selectedProjectId: projectId },
        }))
      },
    },
  }
})
