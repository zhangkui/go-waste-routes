<template>
  <el-card>
    <template #header>数据导出</template>
    <el-form :inline="true">
      <el-form-item label="数据类型">
        <el-select v-model="kind" style="width: 220px">
          <el-option v-for="item in kinds" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="格式">
        <el-select v-model="format" style="width: 140px">
          <el-option label="JSON" value="json" />
          <el-option label="CSV" value="csv" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="preview">预览</el-button>
        <el-button @click="download">下载</el-button>
      </el-form-item>
    </el-form>
    <el-input v-model="content" type="textarea" :rows="18" readonly />
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import http from '@/api/http'

const kinds = ref<string[]>([])
const kind = ref('customers')
const format = ref<'json' | 'csv'>('json')
const content = ref('')

async function loadKinds() {
  const response = await http.get('/exports/kinds')
  kinds.value = response.data.data.kinds ?? []
  if (!kinds.value.includes(kind.value) && kinds.value.length > 0) {
    kind.value = kinds.value[0]
  }
}

async function preview() {
  const response = await http.get('/exports', { params: { kind: kind.value, format: format.value } })
  content.value = typeof response.data === 'string' ? response.data : JSON.stringify(response.data, null, 2)
}

async function download() {
  const response = await http.get('/exports', {
    params: { kind: kind.value, format: format.value },
    responseType: 'blob',
  })
  const url = window.URL.createObjectURL(new Blob([response.data]))
  const link = document.createElement('a')
  link.href = url
  link.download = `${kind.value}.${format.value}`
  link.click()
  window.URL.revokeObjectURL(url)
}

onMounted(async () => {
  await loadKinds()
  await preview()
})
</script>
