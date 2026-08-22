<template>
  <el-card>
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新增</el-button>
      <el-input v-model="query" placeholder="搜索" style="max-width: 260px" />
    </div>
    <el-table :data="items" border>
      <el-table-column prop="id" label="ID" width="90" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="status" label="状态" />
      <el-table-column label="操作" width="220">
        <template #default="{ row }">
          <el-button text @click="editRow(row)">编辑</el-button>
          <el-button v-if="row.status !== undefined" text @click="toggleStatus(row)">启停</el-button>
          <el-button text type="danger" @click="removeRow(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
  <el-dialog v-model="dialogVisible" :title="title">
    <el-input v-model="editor" type="textarea" :rows="14" />
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="save">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import http from '@/api/http'

const route = useRoute()
const apiPath = computed(() => String(route.meta.apiPath ?? '/users'))
const title = computed(() => String(route.meta.title ?? '实体'))
const items = ref<Array<Record<string, unknown>>>([])
const query = ref('')
const dialogVisible = ref(false)
const editor = ref('{}')
const editingId = ref<number | null>(null)

async function load() {
  const response = await http.get(apiPath.value, { params: { page: 1, page_size: 50, q: query.value } })
  items.value = response.data.data.items ?? []
}

function openCreate() {
  editingId.value = null
  editor.value = '{}'
  dialogVisible.value = true
}

function editRow(row: Record<string, unknown>) {
  editingId.value = Number(row.id)
  editor.value = JSON.stringify(row, null, 2)
  dialogVisible.value = true
}

async function save() {
  const payload = JSON.parse(editor.value)
  if (editingId.value) {
    await http.put(`${apiPath.value}/${editingId.value}`, payload)
  } else {
    await http.post(apiPath.value, payload)
  }
  dialogVisible.value = false
  await load()
}

async function removeRow(row: Record<string, unknown>) {
  await http.delete(`${apiPath.value}/${row.id}`)
  await load()
}

async function toggleStatus(row: Record<string, unknown>) {
  await http.put(`${apiPath.value}/${row.id}/status`, { status: row.status === 'enabled' ? 'disabled' : 'enabled' })
  await load()
}

watch(query, load)
onMounted(load)
</script>
