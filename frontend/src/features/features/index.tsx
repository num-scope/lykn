import { useCallback, useEffect, useMemo, useState } from 'react'
import { Loader2, Plus, RefreshCw, Search } from 'lucide-react'
import { toast } from 'sonner'
import {
  createFeature,
  deleteFeature,
  fetchFeatures,
  updateFeature,
} from '@/api/lykn'
import { formatDateTime, getErrorMessage } from '@/lib/lykn'
import type { FeaturePayload, FeatureRecord } from '@/types/api'
import { ConfirmDialog } from '@/components/confirm-dialog'
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
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'

const emptyForm: FeaturePayload = {
  code: '',
  name: '',
  description: '',
  enabled: true,
}

export function FeaturesPage() {
  const [loading, setLoading] = useState(true)
  const [rows, setRows] = useState<FeatureRecord[]>([])
  const [keyword, setKeyword] = useState('')
  const [enabledFilter, setEnabledFilter] = useState<string>('all')

  const [modalOpen, setModalOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<FeatureRecord | null>(null)
  const [form, setForm] = useState<FeaturePayload>(emptyForm)

  const [deleteTarget, setDeleteTarget] = useState<FeatureRecord | null>(null)
  const [deleting, setDeleting] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      setRows(await fetchFeatures())
    } catch (error) {
      toast.error(getErrorMessage(error, '加载功能列表失败'))
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
        item.code.toLowerCase().includes(q) ||
        item.name.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q)
      const matchEnabled =
        enabledFilter === 'all' || String(item.enabled) === enabledFilter
      return matchKeyword && matchEnabled
    })
  }, [rows, keyword, enabledFilter])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setModalOpen(true)
  }

  function openEdit(record: FeatureRecord) {
    setEditing(record)
    setForm({
      code: record.code,
      name: record.name,
      description: record.description,
      enabled: record.enabled,
    })
    setModalOpen(true)
  }

  async function submitForm() {
    if (!form.code.trim() || !form.name.trim()) {
      toast.error('请填写功能编码和名称')
      return
    }
    setSubmitting(true)
    try {
      const payload: FeaturePayload = {
        code: form.code.trim(),
        name: form.name.trim(),
        description: form.description.trim(),
        enabled: form.enabled,
      }
      if (editing) {
        await updateFeature(editing.id, payload)
        toast.success('功能已更新')
      } else {
        await createFeature(payload)
        toast.success('功能已创建')
      }
      setModalOpen(false)
      await load()
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
      await deleteFeature(deleteTarget.id)
      toast.success('功能已删除')
      setDeleteTarget(null)
      await load()
    } catch (error) {
      toast.error(getErrorMessage(error, '删除功能失败'))
    } finally {
      setDeleting(false)
    }
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
            <h1 className='text-2xl font-bold tracking-tight'>功能管理</h1>
            <p className='text-muted-foreground'>定义可授权的功能点</p>
          </div>
          <div className='flex items-center gap-2'>
            <Button variant='outline' size='sm' onClick={() => void load()}>
              <RefreshCw className={loading ? 'animate-spin' : ''} />
              刷新
            </Button>
            <Button size='sm' onClick={openCreate}>
              <Plus />
              新建功能
            </Button>
          </div>
        </div>

        <Card className='mb-4'>
          <CardContent className='flex flex-wrap gap-3 pt-6'>
            <div className='relative min-w-[220px] flex-1'>
              <Search className='absolute start-2.5 top-2.5 size-4 text-muted-foreground' />
              <Input
                className='ps-8'
                placeholder='搜索编码 / 名称 / 描述'
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
              />
            </div>
            <Select value={enabledFilter} onValueChange={setEnabledFilter}>
              <SelectTrigger className='w-[140px]'>
                <SelectValue placeholder='状态' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>全部状态</SelectItem>
                <SelectItem value='true'>启用</SelectItem>
                <SelectItem value='false'>停用</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>功能列表</CardTitle>
            <CardDescription>共 {filtered.length} 个功能</CardDescription>
          </CardHeader>
          <CardContent className='overflow-x-auto'>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>功能</TableHead>
                  <TableHead>编码</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>描述</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className='text-end'>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={6} className='h-24 text-center'>
                      <Loader2 className='mx-auto animate-spin' />
                    </TableCell>
                  </TableRow>
                ) : filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className='h-24 text-center text-muted-foreground'>
                      暂无功能
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((record) => (
                    <TableRow key={record.id}>
                      <TableCell className='font-medium'>{record.name}</TableCell>
                      <TableCell className='font-mono text-sm'>{record.code}</TableCell>
                      <TableCell>
                        <Badge variant={record.enabled ? 'default' : 'secondary'}>
                          {record.enabled ? '启用' : '停用'}
                        </Badge>
                      </TableCell>
                      <TableCell className='max-w-[260px] truncate text-muted-foreground'>
                        {record.description || '-'}
                      </TableCell>
                      <TableCell>{formatDateTime(record.updated_at)}</TableCell>
                      <TableCell className='text-end'>
                        <div className='flex justify-end gap-1'>
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
            <DialogTitle>{editing ? '编辑功能' : '新建功能'}</DialogTitle>
          </DialogHeader>
          <div className='grid gap-4 py-2'>
            <div className='grid gap-2'>
              <Label>功能编码</Label>
              <Input
                value={form.code}
                onChange={(e) => setForm((s) => ({ ...s, code: e.target.value }))}
                placeholder='例如：advanced-export'
              />
            </div>
            <div className='grid gap-2'>
              <Label>功能名称</Label>
              <Input
                value={form.name}
                onChange={(e) => setForm((s) => ({ ...s, name: e.target.value }))}
                placeholder='请输入功能名称'
              />
            </div>
            <div className='grid gap-2'>
              <Label>功能描述</Label>
              <Textarea
                value={form.description}
                onChange={(e) => setForm((s) => ({ ...s, description: e.target.value }))}
                placeholder='请输入功能描述'
              />
            </div>
            <div className='flex items-center justify-between rounded-md border p-3'>
              <div>
                <div className='text-sm font-medium'>启用状态</div>
                <div className='text-xs text-muted-foreground'>停用后不可被新套餐选择</div>
              </div>
              <Switch
                checked={form.enabled}
                onCheckedChange={(enabled) => setForm((s) => ({ ...s, enabled }))}
              />
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

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title='删除功能'
        desc='确定删除该功能？已被套餐使用的功能不能删除。'
        confirmText={deleting ? '删除中...' : '删除'}
        destructive
        isLoading={deleting}
        handleConfirm={() => void confirmDelete()}
      />
    </>
  )
}
