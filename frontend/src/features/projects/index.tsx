import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, Plus, RefreshCw, Search } from 'lucide-react'
import { toast } from 'sonner'
import {
  createProject,
  deleteProject,
  downloadProjectPublicKey,
  fetchProject,
  fetchProjects,
  updateProject,
} from '@/api/lykn'
import { formatDateTime, getErrorMessage, saveBlobFile } from '@/lib/lykn'
import { useAuthStore } from '@/stores/auth-store'
import type { ProjectRecord } from '@/types/api'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'
import { ConfirmDialog } from '@/components/confirm-dialog'

const KEY_BITS = [2048, 3072, 4096] as const

type FormState = {
  name: string
  description: string
  key_bits: number
}

const emptyForm: FormState = {
  name: '',
  description: '',
  key_bits: 2048,
}

export function ProjectsPage() {
  const navigate = useNavigate()
  const loadProjects = useAuthStore((s) => s.auth.loadProjects)
  const selectProject = useAuthStore((s) => s.auth.selectProject)

  const [loading, setLoading] = useState(true)
  const [rows, setRows] = useState<ProjectRecord[]>([])
  const [keyword, setKeyword] = useState('')
  const [keyBitsFilter, setKeyBitsFilter] = useState<string>('all')

  const [modalOpen, setModalOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<ProjectRecord | null>(null)
  const [form, setForm] = useState<FormState>(emptyForm)

  const [deleteTarget, setDeleteTarget] = useState<ProjectRecord | null>(null)
  const [deleting, setDeleting] = useState(false)

  const [publicKeyOpen, setPublicKeyOpen] = useState(false)
  const [publicKeyLoading, setPublicKeyLoading] = useState(false)
  const [publicKeyProject, setPublicKeyProject] = useState<ProjectRecord | null>(null)
  const [publicKey, setPublicKey] = useState('')
  const [downloadingId, setDownloadingId] = useState<number>()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const data = await fetchProjects()
      setRows(data)
    } catch (error) {
      toast.error(getErrorMessage(error, '加载项目列表失败'))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const filtered = useMemo(() => {
    const q = keyword.trim().toLowerCase()
    return rows.filter((item) => {
      const matchKeyword =
        !q ||
        item.name.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q)
      const matchBits =
        keyBitsFilter === 'all' || item.key_bits === Number(keyBitsFilter)
      return matchKeyword && matchBits
    })
  }, [rows, keyword, keyBitsFilter])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setModalOpen(true)
  }

  function openEdit(record: ProjectRecord) {
    setEditing(record)
    setForm({
      name: record.name,
      description: record.description,
      key_bits: record.key_bits,
    })
    setModalOpen(true)
  }

  async function submitForm() {
    if (!form.name.trim()) {
      toast.error('请输入项目名称')
      return
    }
    setSubmitting(true)
    try {
      if (editing) {
        await updateProject(editing.id, {
          name: form.name.trim(),
          description: form.description.trim(),
        })
        toast.success('项目已更新')
      } else {
        await createProject({
          name: form.name.trim(),
          description: form.description.trim(),
          key_bits: form.key_bits,
        })
        toast.success('项目已创建')
      }
      setModalOpen(false)
      await load()
      await loadProjects()
    } catch (error) {
      toast.error(getErrorMessage(error))
    } finally {
      setSubmitting(false)
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return
    setDeleting(true)
    try {
      await deleteProject(deleteTarget.id)
      toast.success('项目已删除')
      setDeleteTarget(null)
      await load()
      await loadProjects()
    } catch (error) {
      toast.error(getErrorMessage(error, '删除项目失败'))
    } finally {
      setDeleting(false)
    }
  }

  async function showPublicKey(record: ProjectRecord) {
    setPublicKeyOpen(true)
    setPublicKeyProject(record)
    setPublicKey('')
    setPublicKeyLoading(true)
    try {
      const detail = await fetchProject(record.id)
      setPublicKeyProject(detail)
      setPublicKey(detail.public_key || '')
    } catch (error) {
      toast.error(getErrorMessage(error, '获取公钥失败'))
    } finally {
      setPublicKeyLoading(false)
    }
  }

  async function copyPublicKey() {
    if (!publicKey) return
    await navigator.clipboard.writeText(publicKey)
    toast.success('公钥已复制')
  }

  async function handleDownloadPublicKey(record: ProjectRecord) {
    setDownloadingId(record.id)
    try {
      saveBlobFile(await downloadProjectPublicKey(record.id))
    } catch (error) {
      toast.error(getErrorMessage(error, '下载公钥失败'))
    } finally {
      setDownloadingId(undefined)
    }
  }

  function goLicenses(record: ProjectRecord) {
    selectProject(record.id)
    navigate({ to: '/licenses' })
  }

  return (
    <>
      <Header>
        <div className='ms-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='mb-4 flex flex-wrap items-center justify-between gap-3'>
          <div>
            <h1 className='text-2xl font-bold tracking-tight'>项目管理</h1>
            <p className='text-muted-foreground'>创建 RSA 密钥对并管理授权项目</p>
          </div>
          <div className='flex items-center gap-2'>
            <Button variant='outline' size='sm' onClick={() => void load()}>
              <RefreshCw className={loading ? 'animate-spin' : ''} />
              刷新
            </Button>
            <Button size='sm' onClick={openCreate}>
              <Plus />
              新建项目
            </Button>
          </div>
        </div>

        <Card className='mb-4'>
          <CardContent className='flex flex-wrap gap-3 pt-6'>
            <div className='relative min-w-[220px] flex-1'>
              <Search className='absolute start-2.5 top-2.5 size-4 text-muted-foreground' />
              <Input
                className='ps-8'
                placeholder='搜索项目名称 / 描述'
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
              />
            </div>
            <Select value={keyBitsFilter} onValueChange={setKeyBitsFilter}>
              <SelectTrigger className='w-[160px]'>
                <SelectValue placeholder='密钥位数' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>全部位数</SelectItem>
                {KEY_BITS.map((bits) => (
                  <SelectItem key={bits} value={String(bits)}>
                    {bits} bit
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>项目列表</CardTitle>
            <CardDescription>共 {filtered.length} 个项目</CardDescription>
          </CardHeader>
          <CardContent className='overflow-x-auto'>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>项目</TableHead>
                  <TableHead>描述</TableHead>
                  <TableHead>密钥</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className='text-end'>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={5} className='h-24 text-center text-muted-foreground'>
                      <Loader2 className='mx-auto animate-spin' />
                    </TableCell>
                  </TableRow>
                ) : filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className='h-24 text-center text-muted-foreground'>
                      暂无项目
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((record) => (
                    <TableRow key={record.id}>
                      <TableCell className='font-medium'>{record.name}</TableCell>
                      <TableCell className='max-w-[280px] truncate text-muted-foreground'>
                        {record.description || '-'}
                      </TableCell>
                      <TableCell>
                        <Badge variant='secondary'>{record.key_bits} bit</Badge>
                      </TableCell>
                      <TableCell>{formatDateTime(record.updated_at)}</TableCell>
                      <TableCell className='text-end'>
                        <div className='flex flex-wrap justify-end gap-1'>
                          <Button size='sm' variant='outline' onClick={() => goLicenses(record)}>
                            License
                          </Button>
                          <Button size='sm' variant='outline' onClick={() => void showPublicKey(record)}>
                            公钥
                          </Button>
                          <Button
                            size='sm'
                            variant='outline'
                            disabled={downloadingId === record.id}
                            onClick={() => void handleDownloadPublicKey(record)}
                          >
                            下载
                          </Button>
                          <Button size='sm' variant='outline' onClick={() => openEdit(record)}>
                            编辑
                          </Button>
                          <Button
                            size='sm'
                            variant='destructive'
                            onClick={() => setDeleteTarget(record)}
                          >
                            删除
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </Main>

      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? '编辑项目' : '新建项目'}</DialogTitle>
          </DialogHeader>
          <div className='grid gap-4 py-2'>
            <div className='grid gap-2'>
              <Label htmlFor='project-name'>项目名称</Label>
              <Input
                id='project-name'
                value={form.name}
                onChange={(e) => setForm((s) => ({ ...s, name: e.target.value }))}
                placeholder='请输入项目名称'
              />
            </div>
            <div className='grid gap-2'>
              <Label htmlFor='project-desc'>项目描述</Label>
              <Textarea
                id='project-desc'
                value={form.description}
                onChange={(e) => setForm((s) => ({ ...s, description: e.target.value }))}
                placeholder='请输入项目描述'
              />
            </div>
            <div className='grid gap-2'>
              <Label>RSA 密钥位数</Label>
              <Select
                value={String(form.key_bits)}
                disabled={!!editing}
                onValueChange={(value) =>
                  setForm((s) => ({ ...s, key_bits: Number(value) }))
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {KEY_BITS.map((bits) => (
                    <SelectItem key={bits} value={String(bits)}>
                      {bits} bit
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant='outline' onClick={() => setModalOpen(false)}>
              取消
            </Button>
            <Button disabled={submitting} onClick={() => void submitForm()}>
              {submitting && <Loader2 className='animate-spin' />}
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={publicKeyOpen} onOpenChange={setPublicKeyOpen}>
        <DialogContent className='max-w-2xl'>
          <DialogHeader>
            <DialogTitle>{publicKeyProject?.name || '项目'} 公钥</DialogTitle>
          </DialogHeader>
          <p className='text-sm text-muted-foreground'>
            客户端 SDK 或 CLI 可使用该公钥验证离线 License
          </p>
          {publicKeyLoading ? (
            <div className='flex h-40 items-center justify-center'>
              <Loader2 className='animate-spin' />
            </div>
          ) : (
            <pre className='max-h-80 overflow-auto rounded-md bg-muted p-4 text-xs leading-6'>
              {publicKey || '暂无公钥'}
            </pre>
          )}
          <DialogFooter>
            <Button variant='outline' disabled={!publicKey} onClick={() => void copyPublicKey()}>
              复制公钥
            </Button>
            {publicKeyProject && (
              <Button
                disabled={downloadingId === publicKeyProject.id}
                onClick={() => void handleDownloadPublicKey(publicKeyProject)}
              >
                下载 public.pem
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title='删除项目'
        desc='确定删除该项目？相关 License 会一并删除。'
        confirmText={deleting ? '删除中...' : '删除'}
        destructive
        isLoading={deleting}
        handleConfirm={() => void confirmDelete()}
      />
    </>
  )
}
